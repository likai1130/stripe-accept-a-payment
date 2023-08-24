package api

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"stripe-accept-a-payment/payment"
)

// HandleCreateCheckout 创建结账
func HandleCreateCheckout(c *gin.Context) {
	/*params := payment.CheckoutParams{}
	if err := c.ShouldBindJSON(&params); err != nil {

		c.JSON(http.StatusBadRequest, errors.New(fmt.Sprintf("结算参数错误：%v", err)).Error())
		return
	}*/
	params := payment.CheckoutParams{
		OrderId:            "ws_1231312312",
		CustomerEmail:      "aa@bb.com",
		ProductName:        "基础版本套餐",
		ProductDescription: "这是测试stripe的套餐",
		UnitAmount:         9900,
	}
	//params.UnitAmount = params.UnitAmount * 100
	session, err := payment.CreateCheckoutSession(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, session.URL)
	return
}

// HandleCheckoutSession 处理结账结果
func HandleCheckoutSession(c *gin.Context) {
	sessionId := c.Query("sessionId")
	session := payment.HandleCheckoutSession(sessionId)
	c.JSON(http.StatusOK, session)
}

// HandleWebhook 调用webhook
func HandleWebhook(c *gin.Context) {
	signature := c.GetHeader("Stripe-Signature")
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("webhook read io error: %v \n", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	payment.HandleWebhook(signature, b)
}
