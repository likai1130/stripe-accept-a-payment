package payment

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"stripe-accept-a-payment/config"
)

type CheckoutParams struct {
	OrderId            string `json:"order_id"`            //业务订单号
	CustomerEmail      string `json:"customer_email"`      //客户邮箱
	ProductName        string `json:"product_name"`        //产品名称
	ProductDescription string `json:"product_description"` //产品描述
	UnitAmount         int64  `json:"unit_amount"`         //产品价格
}

func CreateCheckoutSession(cp CheckoutParams) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		ClientReferenceID: stripe.String(cp.OrderId),
		CustomerEmail:     stripe.String(cp.CustomerEmail),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Description: stripe.String(cp.ProductDescription),
						Name:        stripe.String(cp.ProductName),
						//	TaxCode:     stripe.String("txcd_10101000"),
					},
					UnitAmount: stripe.Int64(cp.UnitAmount),
				},
				Quantity: stripe.Int64(1),
				TaxRates: nil,
			},
		},
		//	Locale:     nil, // 可能就是地区
		Mode:         stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:   stripe.String(config.Cfg.Stripe.Payment.SuccessUrl),
		CancelURL:    stripe.String(config.Cfg.Stripe.Payment.CancelUrl),
		AutomaticTax: &stripe.CheckoutSessionAutomaticTaxParams{Enabled: stripe.Bool(true)},
	}
	s, err := session.New(params)
	if err != nil {
		return s, errors.New(fmt.Sprintf("error while creating session %v", err.Error()))
	}
	return s, nil
}

func HandleCheckoutSession(sessionID string) *stripe.CheckoutSession {
	s, _ := session.Get(sessionID, nil)
	return s
}