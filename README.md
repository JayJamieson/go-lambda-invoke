# go-lambda-invoke

Wrapper library around `aws-sdk-go-v2/service/lambda` to simplify invoking lambda functions with less ceremony.

Installing

To start using go-invoke-lambda, install Go and run `go get`:

```sh
go get github.com/JayJamieson/go-lambda-invoke
```

## Usage

```go
package main

import (
 "log"

 invoke "github.com/JayJamieson/go-lambda-invoke"
)

type input struct {
 Value string `json:"name"`
}

type output struct {
 Value string `json:"name"`
}

func main() {
  var client // create your lambda client
  // or use our thing wrapper

  // client, err := invoke.NewDefaultClient(context.TODO())

  var out output

  // Synchronous invoke
  err := invoke.InvokeSync(context.TODO(), client, &invoke.InvokeInput{
    Name:      "test",
    Qualifier: invoke.DefaultAlias,
    Payload:   input{"hello"},
  }, &out)

 // Asynchronous invoke
  err := invoke.InvokeAsync(context.TODO(), client, &invoke.InvokeInpu{
    Name:      "test",
    Qualifier: invoke.DefaultAlias,
    Payload:   input{"hello"},
  })

}
```
