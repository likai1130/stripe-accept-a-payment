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
			// 完成购买
			FulfillOrder(session)
		} else {
			log.Println("付款失败")
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
	default:
		log.Printf("未知的webhook事件[%s]", event.Type)
	}
	return &Response{
		HttpCode: http.StatusOK,
	}
}

func FulfillOrder(s stripe.CheckoutSession) {
	log.Printf("用户[%s]购买[%s]成功,金额为[$%d],订单id为[%s]: \n", s.CustomerEmail, "123", s.AmountTotal, s.ClientReferenceID)
	// TODO:修改订单状态
	log.Println("订单状态已修改")
	// TODO:发送邮件给用户
	log.Println("凭据已发送")

	log.Println("凭据：", s.Invoice.HostedInvoiceURL)
	//sendInvoice(s.CustomerEmail, s.Customer.ID)

}
