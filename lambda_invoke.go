package lambda_invoke

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

// DefaultAlias is the alias for function invocations.
const DefaultAlias = "$LATEST"

type LambdaClient interface {
	Invoke(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error)
}

type InvokeInput struct {
	Name      string
	Qualifier string
	Payload   any
}

type InvokeError struct {
	// Message is the error message returned from Lambda.
	Message string `json:"errorMessage"`

	// Handled specifies if the error was controlled or not.
	// For example a timeout is unhandled, while an error returned from
	// the function is handled.
	Handled bool
}

// Error message.
func (e *InvokeError) Error() string {
	if e.Handled {
		return fmt.Sprintf("handled: %s", e.Message)
	} else {
		return fmt.Sprintf("unhandled: %s", e.Message)
	}
}

func NewDefaultClient(ctx context.Context) (*lambda.Client, error) {
	config, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		return nil, err
	}

	client := lambda.NewFromConfig(config)
	return client, nil
}

func InvokeSync(ctx context.Context, client LambdaClient, in *InvokeInput, out any) error {
	b, err := json.Marshal(in.Payload)

	if err != nil {
		return fmt.Errorf("marshalling input: %w", err)
	}

	res, err := client.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   &in.Name,
		InvocationType: types.InvocationTypeRequestResponse,
		Qualifier:      &in.Qualifier,
		Payload:        b,
	})

	if err != nil {
		return fmt.Errorf("invoking function: %w", err)
	}

	if res.FunctionError != nil {
		err := &InvokeError{
			Handled: *res.FunctionError == "Handled",
		}

		if e := json.Unmarshal(res.Payload, &err); e != nil {
			return fmt.Errorf("unmarshalling error response: %w", e)
		}

		return err
	}

	if err := json.Unmarshal(res.Payload, &out); err != nil {
		return fmt.Errorf("unmarshalling response: %w", err)
	}

	return nil
}

func InvokeAsync(ctx context.Context, client LambdaClient, in *InvokeInput) error {
	b, err := json.Marshal(in.Payload)

	if err != nil {
		return fmt.Errorf("marshalling input: %w", err)
	}

	_, err = client.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   &in.Name,
		InvocationType: types.InvocationTypeEvent,
		Qualifier:      &in.Qualifier,
		Payload:        b,
	})

	if err != nil {
		return fmt.Errorf("invoking function: %w", err)
	}

	return nil
}
