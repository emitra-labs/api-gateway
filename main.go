package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/appleboy/graceful"
	"github.com/caitlinelfring/go-env-default"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/iancoleman/strcase"
)

type Service struct {
	Path    string
	Address string
}

var port = env.GetIntDefault("PORT", 3000)
var services = []Service{}

func init() {
	// Collect backend services from env variables
	for _, e := range os.Environ() {
		segments := strings.Split(e, "=")
		key := segments[0]
		value := segments[1]

		if strings.HasSuffix(key, "_HTTP_ADDRESS") {
			services = append(services, Service{
				Path:    "/" + strcase.ToKebab(strings.Replace(key, "_HTTP_ADDRESS", "", 1)),
				Address: value,
			})
		}
	}
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	for _, service := range services {
		app.Use(service.Path, proxy.Balancer(proxy.Config{
			Servers: []string{service.Address},
			ModifyRequest: func(c *fiber.Ctx) error {
				requestURI := string(c.Request().RequestURI())
				// Rewrite request uri by eliminating service path prefix,
				// then send it to the backend service.
				requestURI = strings.Replace(requestURI, service.Path, "", 1)
				c.Request().SetRequestURI(requestURI)
				return nil
			},
			Timeout: 5 * time.Second,
		}))
	}

	m := graceful.NewManager()

	m.AddRunningJob(func(ctx context.Context) error {
		return app.Listen(fmt.Sprintf(":%d", port))
	})

	m.AddShutdownJob(func() error {
		return app.Shutdown()
	})

	<-m.Done()
}
