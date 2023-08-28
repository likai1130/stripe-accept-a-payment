package payment

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/webhook"
	"log"
	"net/http"
	"stripe-accept-a-payment/config"
)

type Response struct {
	HttpCode int
	Data     []byte
	Err      error
}

func HandleWebhook(signature string, data []byte) *Response {
	event, err := webhook.ConstructEvent(data, signature, config.Cfg.Stripe.WebhookSecretKey)
	if err != nil {
		log.Printf("webhook.ConstructEvent: %v \n", err)
		return &Response{
			HttpCode: http.StatusBadRequest,
			Err:      err,
		}
	}

	switch event.Type {
	case "checkout.session.completed":
		log.Println("Checkout Session completed!")
		var session stripe.CheckoutSession
		err = json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			return &Response{
				HttpCode: http.StatusBadRequest,
				Err:      errors.WithMessage(err, "Error parsing webhook JSON:"),
			}
		}
		//检查订单是否已付款（例如，通过卡付款）
		orderPaid := session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid
		if orderPaid {
			fmt.Printf("[%s]付款成功,订货状态[%s],业务订单号为[%s],StripeCheckout单号[%s] \n", session.CustomerEmail, orderPaid, session.ClientReferenceID, session.ID)
		} else {
			fmt.Printf("[%s]付款成功,订货状态[%s],业务订单号为[%s],StripeCheckout单号[%s] \n", session.CustomerEmail, orderPaid, session.ClientReferenceID, session.ID)
		}
	case "invoice.paid": //生成票据
		log.Println("Invoice paid completed!")
		var invoice stripe.Invoice
		err = json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			return &Response{
				HttpCode: http.StatusBadRequest,
				Err:      errors.WithMessage(err, "Error parsing webhook JSON:"),
			}
		}

		fmt.Println("HostedInvoiceURL: ", invoice.HostedInvoiceURL)
		fmt.Println(fmt.Sprintf("用户[%s]-金额[$%d]-业务订单id[%s]", invoice.CustomerEmail, invoice.AmountPaid/100, invoice.Metadata["order_id"]))
		FulfillOrder(invoice)
	default:
		log.Printf("未知的webhook事件[%s]", event.Type)
	}
	return &Response{
		HttpCode: http.StatusOK,
	}
}

// FulfillOrder 处理业务订单状态
func FulfillOrder(invoice stripe.Invoice) {
	fmt.Printf("凭据URL: %s\n", invoice.HostedInvoiceURL)
	fmt.Printf("用户[%s]-金额[$%d]-业务订单id[%s]-invoiceId[%s] \n", invoice.CustomerEmail, invoice.AmountPaid/100, invoice.CustomFields[0].Value, invoice.ID)
	//todo 更新数据库订单状态
	fmt.Printf("业务订单id[%s] 更新成功 \n", invoice.CustomFields[0])
}
