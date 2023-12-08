package main

import (
	"context"
	"log"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/omeiirr/dummy-data-api/helpers"
	"github.com/DataDog/datadog-lambda-go"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	fibertrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gofiber/fiber.v2"
)

var fiberLambda *fiberadapter.FiberLambda

func main() {
	tracer.Start()
	defer tracer.Stop()

	app := fiber.New()

	app.Use(fibertrace.Middleware(fibertrace.WithServiceName("go")))

	app.Get("/", func(c *fiber.Ctx) error {
		// ตรวจสอบ Trace Context ใน context ของ Lambda function
		ctx := c.Context() // หรือใช้ c.Locals("ddtrace.context").(context.Context)
		span, _ := tracer.SpanFromContext(ctx)
	
		// Datadog Lambda Library ควรจะดูแลการสร้าง span ใน Lambda event
	
		// ใช้ span ในที่นี้ (ตัวอย่างเท่านี้)
		span.SetTag("example", "Lambda function span")
	
		return c.SendString("Server is running normally")
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		// ตรวจสอบ Trace Context ใน context ของ Lambda function
		ctx := c.Context() // หรือใช้ c.Locals("ddtrace.context").(context.Context)
		span, _ := tracer.SpanFromContext(ctx)
	
		// ใช้ span ในที่นี้ (ตัวอย่างเท่านี้)
		span.SetTag("example", "HTTP request span")
	
		users := helpers.GenerateUsers(5)
		return c.JSON(users)
	})

	if helpers.IsLambda() {
		fiberLambda = fiberadapter.New(app)
		// ใช้ ddlambda.Wrap แทน ddlambda.WrapFunction
		lambda.Start(ddlambda.WrapFunction(Handler, nil))
	} else {
		app.Listen(":3000")
	}
}

// ใน Handler function ของ Lambda function
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Datadog Lambda Library ควรจะดูแลการสร้าง span จาก Lambda event

	// สร้าง span ใน Lambda function
	span, ctx := tracer.StartSpanFromContext(ctx, "lambda-handler-span")
	defer span.Finish()

	// Log Trace Context ที่ได้รับ
	log.Printf("Lambda Trace Context: %v", span.Context())

	// สร้าง TextMapCarrier จาก Headers
	carrier := tracer.TextMapCarrier(request.Headers)

	// ตรวจสอบว่า carrier ไม่เป็น nil ก่อนใช้ Inject
	if carrier != nil {
		// Log carrier ที่ได้รับ
		log.Printf("Carrier Headers: %v", carrier)

		// ส่ง Trace Context ใน HTTP header ของ response
		tracer.Inject(span.Context(), carrier)
	}

	return fiberLambda.ProxyWithContext(ctx, request)
}

