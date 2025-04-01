package services

import (
	"github.com/stripe/stripe-go/v81"
	stripePrice "github.com/stripe/stripe-go/v81/price"
	"github.com/stripe/stripe-go/v81/product"
)

type StripeService struct{}

func NewStripeService() *StripeService {
	return &StripeService{}
}

func (s *StripeService) CreateNewProduct(name string, price float32) (*stripe.Product, error) {
	params := &stripe.ProductParams{
		Name: stripe.String(name),
		DefaultPriceData: &stripe.ProductDefaultPriceDataParams{
			Currency:   stripe.String(string(stripe.CurrencySGD)),
			UnitAmount: stripe.Int64(int64(price * 100)),
		},
	}

	result, err := product.New(params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *StripeService) UpdateProductPrice(stripeProductId, stripePriceId string, price float32) (*stripe.Product, error) {
	priceParams := &stripe.PriceParams{
		Product:    stripe.String(stripeProductId),
		Currency:   stripe.String(string(stripe.CurrencySGD)),
		UnitAmount: stripe.Int64(int64(price * 100)),
	}

	priceData, err := stripePrice.New(priceParams)
	if err != nil {
		return nil, err
	}

	productParams := &stripe.ProductParams{
		DefaultPrice: stripe.String(priceData.ID),
	}
	result, err := product.Update(stripeProductId, productParams)
	if err != nil {
		return nil, err
	}

	return result, nil
}
