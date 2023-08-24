package payment

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
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
	if event.Type == "checkout.session.completed" {
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
	} else {
		log.Printf("未知的webhook事件[%s]", event.Type)
	}
	return &Response{
		HttpCode: http.StatusOK,
	}
}

func FulfillOrder(session stripe.CheckoutSession) {
	log.Printf("用户[%s]购买[%s]成功,金额为[$%d],订单id为[%s]: \n", session.CustomerEmail, "123", session.AmountTotal, session.ClientReferenceID)
	// TODO:修改订单状态
	log.Println("订单状态已修改")
	// TODO:发送邮件给用户
	log.Println("凭据已发送")
}

func jsonMarshal(session stripe.CheckoutSession) string {
	marshal, _ := json.Marshal(session)
	return string(marshal)
}
