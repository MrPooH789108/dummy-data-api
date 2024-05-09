package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/omeiirr/dummy-data-api/helpers"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"github.com/DataDog/datadog-lambda-go"
)

var fiberLambda *fiberadapter.FiberLambda

func main() {
    tracer.Start(tracer.WithService("go-fiber"))
    defer tracer.Stop()

	app := fiber.New()
	setupRoutes(app)

	if helpers.IsLambda() {
		fiberLambda = fiberadapter.New(app)
		lambda.Start(ddlambda.WrapFunction(Handler, nil))
	} else {
		app.Listen(":3000")
	}
}

func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		span := tracer.StartSpan("web.request", tracer.ResourceName("/"))
		defer span.Finish()
		span.SetTag("url", c.BaseURL())

		return c.SendString("Server is running normally")
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		span := tracer.StartSpan("web.request", tracer.ResourceName("/users"))
		defer span.Finish()
		span.SetTag("url", c.BaseURL())

		users := helpers.GenerateUsers(5)
		return c.JSON(users)
	})
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	app := fiber.New()
	setupRoutes(app)
	return fiberLambda.ProxyWithContext(ctx, request)
}
