package config

// Email represents email configuration
type Email struct {
	From     string `envconfig:"EMAIL_FROM" default:"noreply@fiberboilerplate.com"`
	Host     string `envconfig:"EMAIL_HOST" default:"smtp.ethereal.email"`
	Port     string `envconfig:"EMAIL_PORT" default:"587"`
	Username string `envconfig:"EMAIL_USERNAME" default:""`
	Password string `envconfig:"EMAIL_PASSWORD" default:""`
}
