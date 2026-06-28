package config

import "time"

// App represents application configuration
type App struct {
	ServiceName    string        `envconfig:"APP_SERVICE_NAME" default:"boilerplate-service"`
	Name           string        `envconfig:"APP_NAME" default:"Fiber Boilerplate"`
	Env            string        `envconfig:"APP_ENV" default:"development"`
	Port           string        `envconfig:"APP_PORT" default:"3000"`
	URL            string        `envconfig:"APP_URL" default:"http://localhost:3000"`
	ContextTimeout time.Duration `envconfig:"CONTEXT_TIMEOUT" default:"30s"`
}
