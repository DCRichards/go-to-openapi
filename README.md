# go-to-openapi

> Convert Go structs to [OpenAPI Schemas (Data Models)](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.3.md#schemaObject).

## Motivation

Documenting a large API can be an extremely cumbersome task, especially if it isn't your API or you've simply lost familiarity with it. Go's structs are the ideal tool for figuring out requests and responses without having to query each endpoint, but writing an OpenAPI with them is equally as cumbersome. Enter `go-to-openapi` which reflectively parses any struct and returns the associated schema.

## Usage

Here's a simple example.

```go
import (
        "fmt"
        "github.com/dcrichards/go-to-openapi/schema"
)

type User struct {
       	ID         string            `json:"id"`
       	Email      string            `json:"email"`
       	Tags       []Tag             `json:"tags"`
       	Properties map[string]string `json:"props"`
}

type Tag struct {
        Name   string `json:"name"`
       	Active bool   `json:"active"`
}

func main() {
        yml, err := schema.Generate(User{})
        if err != nil {
                fmt.Printf("Error: %s", err.Error())
        }

        log.Println(yml)
}
```
This will generate the following schema:

```yaml
schema:
  type: object
  properties:
    email:
      type: string
    id:
      type: string
    props:
      type: object
      properties:
        example:
          type: string
      additionalProperties: true
    tags:
      type: array
      items:
        type: object
        properties:
          active:
            type: boolean
          name:
            type: string
```