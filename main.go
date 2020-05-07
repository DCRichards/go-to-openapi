package main

import (
	"gopkg.in/yaml.v2"
	"reflect"
	"strings"
)

type Prop struct {
	// Type is the data type.
	Type string `yaml:"type"`
	// Format is the format of the type.
	Format string `yaml:"format,omitempty"`
	// Example is example to make OpenAPI specification of your web service clearer.
	Example interface{} `yaml:"example,omitempty"`
	// Items describes the items in an array.
	Items *Prop `yaml:"items,omitempty"`
	// Properties describes the properties of an object.
	Properties map[string]Prop `yaml:"properties,omitempty"`
	// AdditionalProperties is used when the object properties are variable
	AdditionalProperties *bool `yaml:"additionalProperties,omitempty"`
}

func parseDeep(v reflect.Value, name string, out map[string]Prop) map[string]Prop {
	switch v.Kind() {
	case reflect.Ptr:
		if !v.IsNil() {
			return parseDeep(v.Elem(), name, out)
		}
		// TODO: if it is nil, generate a new one of the underlying type so it is included.
	case reflect.String:
		out[name] = Prop{Type: "string"}
	case reflect.Bool:
		out[name] = Prop{Type: "boolean"}
	case reflect.Int, reflect.Int8, reflect.Int16:
		out[name] = Prop{Type: "integer"}
	case reflect.Int32:
		out[name] = Prop{Type: "integer", Format: "int32"}
	case reflect.Int64:
		out[name] = Prop{Type: "integer", Format: "int64"}
	case reflect.Float32, reflect.Float64:
		out[name] = Prop{Type: "number", Format: "float"}
	case reflect.Struct:
		switch v.Type().String() {
		// time.Time marshals as RFC3339 string.
		case "time.Time":
			out[name] = Prop{Type: "string", Format: "date-time"}
		default:
			p := Prop{Type: "object", Properties: map[string]Prop{}}

			for i := 0; i < v.NumField(); i++ {
				jsonTag := strings.Split(v.Type().Field(i).Tag.Get("json"), ",")
				if jsonTag[0] != "" {
					p.Properties = parseDeep(v.Field(i), jsonTag[0], p.Properties)
				}
			}

			out[name] = p
		}
	case reflect.Slice, reflect.Array:
		p := Prop{Type: "array"}
		// TODO: Get the value when the array is not nil.
		v2 := reflect.New(v.Type().Elem())
		dummy := parseDeep(v2, "dummy", map[string]Prop{})
		d := dummy["dummy"]
		p.Items = &d

		out[name] = p
	case reflect.Map:
		additionalProps := true
		p := Prop{
			Type: "object",
			Properties: map[string]Prop{},
			AdditionalProperties: &additionalProps,
		}

		v3 := reflect.New(v.Type().Elem())
		p.Properties = parseDeep(v3, "example", p.Properties)
		out[name] = p
	}

	return out
}

func Parse(input interface{}) (string, error) {
	response := map[string]Prop{}

	v := reflect.ValueOf(input)
	response = parseDeep(v, "schema", response)

	d, err := yaml.Marshal(&response)
	if err != nil {
		return "", err
	}

	return string(d), nil
}
