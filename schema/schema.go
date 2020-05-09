package schema

import (
	"gopkg.in/yaml.v2"
	"reflect"
	"strings"
)

const (
	schemaTypeString = "string"
	schemaTypeBool   = "boolean"
	schemaTypeInt    = "integer"
	schemaTypeNumber = "number"
	schemaTypeObject = "object"
	schemaTypeArray  = "array"

	formatInt32    = "int32"
	formatInt64    = "int64"
	formatFloat    = "float"
	formatDateTime = "date-time"
)

// Schema represents an OpenAPI Schema Object
//
// https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.3.md#schema-object
type Schema struct {
	// Type is the data type.
	Type string `yaml:"type"`
	// Format is the format of the type.
	Format string `yaml:"format,omitempty"`
	// Items describes the items in an array.
	Items *Schema `yaml:"items,omitempty"`
	// Properties describes the properties of an object.
	Properties map[string]Schema `yaml:"properties,omitempty"`
	// AdditionalProperties signifies an object with variable properties.
	AdditionalProperties *bool `yaml:"additionalProperties,omitempty"`
	// Example is an example of the given property.
	Example interface{} `yaml:"example,omitempty"`
}

func parseDeep(v reflect.Value, name string, out map[string]Schema) map[string]Schema {
	switch v.Kind() {
	case reflect.Ptr:
		if !v.IsNil() {
			return parseDeep(v.Elem(), name, out)
		}

		return parseDeep(reflect.New(v.Type().Elem()), name, out)
	case reflect.String:
		out[name] = Schema{Type: schemaTypeString}
	case reflect.Bool:
		out[name] = Schema{Type: schemaTypeBool}
	case reflect.Int, reflect.Int8, reflect.Int16:
		out[name] = Schema{Type: schemaTypeInt}
	case reflect.Int32:
		out[name] = Schema{Type: schemaTypeInt, Format: formatInt32}
	case reflect.Int64:
		out[name] = Schema{Type: schemaTypeInt, Format: formatInt64}
	case reflect.Float32, reflect.Float64:
		out[name] = Schema{Type: schemaTypeNumber, Format: formatFloat}
	case reflect.Struct:
		switch v.Type().String() {
		// time.Time marshals as RFC3339 string.
		case "time.Time":
			out[name] = Schema{Type: schemaTypeString, Format: formatDateTime}
		default:
			p := Schema{Type: schemaTypeObject, Properties: map[string]Schema{}}

			for i := 0; i < v.NumField(); i++ {
				jsonTag := strings.Split(v.Type().Field(i).Tag.Get("json"), ",")
				if jsonTag[0] != "" {
					p.Properties = parseDeep(v.Field(i), jsonTag[0], p.Properties)
				}
			}

			out[name] = p
		}
	case reflect.Slice, reflect.Array:
		p := Schema{Type: schemaTypeArray}
		// TODO: Get the value when the array is not nil.
		v2 := reflect.New(v.Type().Elem())
		dummy := parseDeep(v2, "dummy", map[string]Schema{})
		d := dummy["dummy"]
		p.Items = &d

		out[name] = p
	case reflect.Map:
		additionalProps := true
		p := Schema{
			Type:                 schemaTypeObject,
			Properties:           map[string]Schema{},
			AdditionalProperties: &additionalProps,
		}

		v3 := reflect.New(v.Type().Elem())
		p.Properties = parseDeep(v3, "example", p.Properties)
		out[name] = p
	}

	return out
}

// Generate parses the struct and returns a YAML string of the associated OpenAPI Schema.
func Generate(input interface{}) (string, error) {
	response := map[string]Schema{}

	v := reflect.ValueOf(input)
	response = parseDeep(v, "schema", response)

	d, err := yaml.Marshal(&response)
	if err != nil {
		return "", err
	}

	return string(d), nil
}
