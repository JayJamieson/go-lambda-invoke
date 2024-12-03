package lambda_invoke_test

import (
	"context"
	"testing"

	invoke "github.com/JayJamieson/go-lambda-invoke"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/matryer/is"
)

type client struct {
	FunctionError *string
}

func (c *client) Invoke(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error) {
	return &lambda.InvokeOutput{
		FunctionError: c.FunctionError,
		Payload:       params.Payload,
	}, nil
}

type input struct {
	Value string `json:"name"`
}

type output struct {
	Value string `json:"name"`
}

func TestInvokeSync(t *testing.T) {
	assert := is.New(t)
	c := new(client)
	var out output

	err := invoke.InvokeSync(context.TODO(), c, &invoke.InvokeInput{
		Name:      "test",
		Qualifier: invoke.DefaultAlias,
		Payload:   input{"hello"},
	}, &out)

	assert.NoErr(err)
	assert.Equal("hello", out.Value)
}

func TestInvokeSync_noInput(t *testing.T) {
	assert := is.New(t)
	c := new(client)
	var out output

	err := invoke.InvokeSync(context.TODO(), c, &invoke.InvokeInput{
		Name:      "test",
		Qualifier: invoke.DefaultAlias,
		Payload:   nil,
	}, &out)

	assert.NoErr(err)
	assert.Equal("", out.Value)
}

func TestInvokeSync_noOutput(t *testing.T) {
	assert := is.New(t)
	c := new(client)

	err := invoke.InvokeSync(context.TODO(), c, &invoke.InvokeInput{
		Name:      "test",
		Qualifier: invoke.DefaultAlias,
		Payload:   input{"hello"},
	}, nil)

	assert.NoErr(err)
}

func TestInvokeSync_noInput_noOutput(t *testing.T) {
	assert := is.New(t)
	c := new(client)

	err := invoke.InvokeSync(context.TODO(), c, &invoke.InvokeInput{
		Name:      "test",
		Qualifier: invoke.DefaultAlias,
		Payload:   nil,
	}, nil)

	assert.NoErr(err)
}

func TestInvokeSync_error(t *testing.T) {
	assert := is.New(t)
	c := new(client)
	var out output
	c.FunctionError = aws.String("Unhandled")

	err := invoke.InvokeSync(context.TODO(), c, &invoke.InvokeInput{
		Name:      "test",
		Qualifier: invoke.DefaultAlias,
		Payload:   &invoke.InvokeError{Message: "Task timed out after 5.00 seconds"},
	}, &out)
	assert.Equal("unhandled: Task timed out after 5.00 seconds", err.Error())

	e := err.(*invoke.InvokeError)
	assert.True(!e.Handled)
	assert.Equal("Task timed out after 5.00 seconds", e.Message)
}

func TestInvokeAsync(t *testing.T) {
	assert := is.New(t)
	c := new(client)
	err := invoke.InvokeAsync(context.TODO(), c, &invoke.InvokeInput{
		Name:      "test",
		Qualifier: invoke.DefaultAlias,
		Payload:   input{"hello"},
	})
	assert.NoErr(err)
}
