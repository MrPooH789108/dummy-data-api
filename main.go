package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/omeiirr/dummy-data-api/handlers"
	"github.com/omeiirr/dummy-data-api/helpers"
	"github.com/DataDog/datadog-lambda-go"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	fibertrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gofiber/fiber.v2"
)

var fiberLambda *fiberadapter.FiberLambda

func main() {
	tracer.Start()
	defer tracer.Stop()

    //span,_ := tracer.StartSpanFromContext(context.Background(),"")
    //defer span.Finish()

	app := fiber.New()

	app.Use(fibertrace.Middleware(fibertrace.WithServiceName("go-fiber")))

	app.Get("/", handlers.HealthCheck)
	app.Get("/users", handlers.ReturnUsers)

	if helpers.IsLambda() {
		fiberLambda = fiberadapter.New(app)
		lambda.Start(ddlambda.WrapFunction(Handler, nil))
	} else {
		app.Listen(":3000")
	}

}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return fiberLambda.ProxyWithContext(ctx, request)
}
