package main

import (
	"context"
	"log"

	"github.com/corezoid/gitcall-go-runner/runner"
	"github.com/xeipuuv/gojsonschema"
	"fmt"
)

var schema *gojsonschema.Schema

func usercode(_ context.Context, data map[string]interface{}) error {
	result, err := schema.Validate(gojsonschema.NewGoLoader(data))
	if err != nil {
		return err
	}

	if !result.Valid() {
		endErr := fmt.Errorf("object invalid")
		for i, err := range result.Errors() {
			endErr = fmt.Errorf("%s: %d# %s", endErr, i, err)
		}

		return endErr
	}

	return nil
}

func main() {
	runner.Run(usercode)
}

func init() {
	var err error

	rootSchema := gojsonschema.NewStringLoader(`
{
    "type": "object",
	"required": ["foo", "bar"],
    "properties": {
		"foo": { "type": "string", "maxLength": 10 },
		"bar": { "type": "number", "min": 10, "max": 100 }
	}
}
`)
	sl := gojsonschema.NewSchemaLoader()
	sl.Validate = true

	schema, err = sl.Compile(rootSchema)
	if err != nil {
		log.Fatal(err)
	}
}
