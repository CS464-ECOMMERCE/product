package services

import (
	"fmt"
	"product/configs"
	"product/models"
	pb "product/proto"
	"product/storage"
	"time"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"gorm.io/gorm"
)

// PaymentItem represents an item with corresponding quantity in a payment
type PaymentItem struct {
	StripePriceId string `json:"stripe_price_id"`
	Quantity      uint64 `json:"quantity"`
}

type OrderService struct {
	cartClient *CartService
}

func NewOrderService(cartClient *CartService) *OrderService {
	return &OrderService{
		cartClient: cartClient,
	}
}

func (o *OrderService) PlaceOrder(req *pb.PlaceOrderRequest) (*pb.PlaceOrderResponse, error) {
	resp := &pb.PlaceOrderResponse{}

	// Get the user's cart
	cart, err := o.cartClient.GetCart(req.SessionId)
	if err != nil {
		return resp, fmt.Errorf("failed to get cart: %w", err)
	}

	// Validate cart is not empty
	if len(cart.Items) == 0 {
		return resp, fmt.Errorf("cart is empty")
	}

	// Create the order
	order := &models.Order{
		UserId: req.UserId,
		Status: models.OrderStatusProcessing,
	}

	// Start transaction
	tx := storage.StorageInstance.BeginTransaction()

	// Validate cart and get payment items
	paymentItem, err := o.validateCartAndGetPaymentItems(req.UserId, order, cart.Items, tx)
	if err != nil {
		tx.Rollback()
		return resp, fmt.Errorf("failed to validate cart: %w", err)
	}

	// Create order to attach ID into Stripe Payment
	order, err = storage.StorageInstance.Order.CreateOrder(order, tx)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	// Create payment session
	sess, err := o.createNewPayment(order.Id, req.UserEmail, paymentItem)
	if err != nil {
		tx.Rollback()
		return resp, fmt.Errorf("failed to create payment: %w", err)
	}

	// attach checkout session ID
	order.CheckoutSessionId = sess.ID
	if err := storage.StorageInstance.Order.UpdateOrder(order, tx); err != nil {
		tx.Rollback()
		return resp, fmt.Errorf("unable to update order with checkout session id: %w", err)
	}

	// Delete cart last
	if err := o.cartClient.DeleteCart(req.SessionId); err != nil {
		tx.Rollback()
		return resp, fmt.Errorf("failed to delete cart: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	resp.CheckoutUrl = sess.URL
	return resp, nil
}

func (o *OrderService) validateCartAndGetPaymentItems(buyerId uint64, order *models.Order, cartItems []*pb.CartItem, tx *gorm.DB) ([]*PaymentItem, error) {
	paymentItems := make([]*PaymentItem, 0, len(cartItems))
	orderItems := make([]models.OrderItem, 0, len(cartItems))

	for _, item := range cartItems {
		product, err := storage.StorageInstance.Product.Get(item.Id, tx)
		if err != nil {
			return nil, fmt.Errorf("failed to get product: %w", err)
		}

		if product.MerchantId == buyerId {
			return nil, fmt.Errorf("cannot buy your own product: %s", product.Name)
		}

		if product.Inventory < item.Quantity {
			return nil, fmt.Errorf("not enough stock for product: %s", product.Name)
		}

		// Update inventory
		updatedInventory := product.Inventory - item.Quantity
		err = storage.StorageInstance.Product.UpdateInventory(product.Id, updatedInventory, tx)
		if err != nil {
			return nil, fmt.Errorf("failed to update inventory: %w", err)
		}

		paymentItems = append(paymentItems, &PaymentItem{
			StripePriceId: product.StripePriceId,
			Quantity:      item.Quantity,
		})

		// Add Order Items
		orderItems = append(orderItems, models.OrderItem{
			ProductId: product.Id,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})

		// Update Total
		order.Total += product.Price * float32(item.Quantity)
	}

	// Save order items to database
	order.OrderItems = orderItems

	return paymentItems, nil
}

// Create a new payment details for order
func (o *OrderService) createNewPayment(orderId uint64, email string, paymentItemList []*PaymentItem) (*stripe.CheckoutSession, error) {
	var lineItems []*stripe.CheckoutSessionLineItemParams
	for _, item := range paymentItemList {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(item.StripePriceId),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	expiresAt := time.Now().Add(30 * time.Minute).Unix() // default 30 mins to checkout

	// Custom metadata for Stripe Event
	metadata := map[string]string{
		"orderId": fmt.Sprint(orderId),
	}

	// Create a new checkout session
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
			"paynow",
		}),
		LineItems:     lineItems,
		Mode:          stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:    stripe.String(fmt.Sprintf("%s/stripe/success", configs.FRONTEND_URL)),
		CancelURL:     stripe.String(fmt.Sprintf("%s/stripe/cancel?session_id={CHECKOUT_SESSION_ID}", configs.FRONTEND_URL)),
		ExpiresAt:     stripe.Int64(expiresAt),
		CustomerEmail: stripe.String(email),
		Metadata:      metadata,
	}

	sess, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout session: %v", err.Error())
	}

	return sess, nil
}
