package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/omeiirr/dummy-data-api/handlers"
	"github.com/omeiirr/dummy-data-api/helpers"
)

var fiberLambda *fiberadapter.FiberLambda

func main() {
	app := fiber.New()

	app.Get("/", handlers.HealthCheck)
	app.Get("/users", handlers.ReturnUsers)

	if helpers.IsLambda() {
		fiberLambda = fiberadapter.New(app)
		lambda.Start(Handler)
	} else {
		app.Listen(":3000")
	}

}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return fiberLambda.ProxyWithContext(ctx, request)
}
