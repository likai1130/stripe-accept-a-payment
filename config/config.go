package config

var Cfg Config

type Config struct {
	Server Server
	Stripe Stripe
}
type Server struct {
	Port int `json:"port"`
}

type Stripe struct {
	PublishableKey   string  `json:"publishable_key"`
	SecretKey        string  `json:"secret_key"`
	WebhookSecretKey string  `json:"webhook_secret_key"`
	Payment          Payment `json:"payment"`
}

type Payment struct {
	SuccessUrl  string `json:"success_url"`
	CancelUrl   string `json:"cancel_url"`
	MethodTypes string `json:"method_types"`
}
