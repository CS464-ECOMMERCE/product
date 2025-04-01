package services

import (
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/product"
)

type StripeService struct{}

func NewStripeService() *StripeService {
	return &StripeService{}
}

func (s *StripeService) CreateNewProduct(name string, price float32) (*stripe.Product, error) {
	stripePrice := int64(price * 100) // Convert to cents
	params := &stripe.ProductParams{
		Name: stripe.String(name),
		DefaultPriceData: &stripe.ProductDefaultPriceDataParams{
			Currency:   stripe.String(string(stripe.CurrencySGD)),
			UnitAmount: stripe.Int64(stripePrice),
		},
	}

	result, err := product.New(params)
	if err != nil {
		return nil, err
	}

	return result, nil
}
