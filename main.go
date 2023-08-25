package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v75"
	"stripe-accept-a-payment/api"
	"stripe-accept-a-payment/config"
)

func init() {
	config.Setup("./config/config.yaml")
	stripe.Key = config.Cfg.Stripe.SecretKey
}

func main() {
	r := gin.Default()
	router(r)
	err := r.Run(fmt.Sprintf(":%d", config.Cfg.Server.Port))
	if err != nil {
		panic(err)
	}
}

func router(r *gin.Engine) {
	r.Static("/static", "./web")
	r.GET("/checkout-session", api.HandleCheckoutSession)
	r.POST("/create-checkout-session", api.HandleCreateCheckout)
	r.POST("/webhook", api.HandleWebhook)

}
