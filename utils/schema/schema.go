// Package schema implements OpenAPI 3 compatible JSON Schema which can be
// generated from structs.
package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ErrSchemaInvalid is sent when there is a problem building the schema.
var ErrSchemaInvalid = errors.New("schema is invalid")

// Mode defines whether the schema is being generated for read or
// write mode. Read-only fields are dropped when in write mode, for example.
type Mode int

const (
	// ModeAll is for general purpose use and includes all fields.
	ModeAll Mode = iota
	// ModeRead is for HTTP HEAD & GET and will hide write-only fields.
	ModeRead
	// ModeWrite is for HTTP POST, PUT, PATCH, DELETE and will hide
	// read-only fields.
	ModeWrite
)

// JSON Schema type constants
const (
	TypeBoolean = "boolean"
	TypeInteger = "integer"
	TypeNumber  = "number"
	TypeString  = "string"
	TypeArray   = "array"
	TypeObject  = "object"
	TypeMap     = "map"
)

var (
	timeType      = reflect.TypeOf(time.Time{})
	ipType        = reflect.TypeOf(net.IP{})
	uriType       = reflect.TypeOf(url.URL{})
	byteSliceType = reflect.TypeOf([]byte(nil))
)

// I returns a pointer to the given int. Useful helper function for pointer
// schema validators like MaxLength or MinItems.
func I(value uint64) *uint64 {
	return &value
}

// F returns a pointer to the given float64. Useful helper function for pointer
// schema validators like Maximum or Minimum.
func F(value float64) *float64 {
	return &value
}

// getTagValue returns a value of the schema's type for the given tag string.
// Uses JSON parsing if the schema is not a string.
func getTagValue(s *Schema, t reflect.Type, value string) (interface{}, error) {
	// Special case: strings don't need quotes.
	if s.Type == TypeString {
		return value, nil
	}

	// Special case: array of strings with comma-separated values and no quotes.
	if s.Type == TypeArray && s.Items != nil && s.Items.Type == TypeString && len(value) > 0 && value[0] != '[' {
		values := []string{}
		for _, s := range strings.Split(value, ",") {
			values = append(values, strings.TrimSpace(s))
		}
		return values, nil
	}

	var v interface{}
	if err := json.Unmarshal([]byte(value), &v); err != nil {
		return nil, err
	}

	vv := reflect.ValueOf(v)
	tv := reflect.TypeOf(v)
	if v != nil && tv != t {
		if tv.Kind() == reflect.Slice {
			// Slices can't be cast due to the different layouts. Instead, we make a
			// new instance of the destination slice, and convert each value in
			// the original to the new type.
			tmp := reflect.MakeSlice(t, 0, vv.Len())
			for i := 0; i < vv.Len(); i++ {
				if !vv.Index(i).Elem().Type().ConvertibleTo(t.Elem()) {
					return nil, fmt.Errorf("unable to convert %v to %v: %w", vv.Index(i).Interface(), t.Elem(), ErrSchemaInvalid)
				}

				tmp = reflect.Append(tmp, vv.Index(i).Elem().Convert(t.Elem()))
			}
			v = tmp.Interface()
		} else if !tv.ConvertibleTo(t) {
			return nil, fmt.Errorf("unable to convert %v to %v: %w", tv, t, ErrSchemaInvalid)
		}

		v = reflect.ValueOf(v).Convert(t).Interface()
	}

	return v, nil
}

// Schema represents a JSON Schema which can be generated from Go structs
type Schema struct {
	Name                 string              `json:"name,omitempty"`
	Type                 string              `json:"type,omitempty"`
	Description          string              `json:"description,omitempty"`
	Items                *Schema             `json:"items,omitempty"`
	Properties           []*Schema           `json:"properties,omitempty"`
	AdditionalProperties interface{}         `json:"additionalProperties,omitempty"`
	PatternProperties    map[string]*Schema  `json:"patternProperties,omitempty"`
	Required             bool                `json:"required,omitempty"`
	Format               string              `json:"format,omitempty"`
	Enum                 []interface{}       `json:"enum,omitempty"`
	Default              interface{}         `json:"default,omitempty"`
	Example              interface{}         `json:"example,omitempty"`
	Minimum              *float64            `json:"minimum,omitempty"`
	ExclusiveMinimum     *bool               `json:"exclusiveMinimum,omitempty"`
	Maximum              *float64            `json:"maximum,omitempty"`
	ExclusiveMaximum     *bool               `json:"exclusiveMaximum,omitempty"`
	MultipleOf           float64             `json:"multipleOf,omitempty"`
	MinLength            *uint64             `json:"minLength,omitempty"`
	MaxLength            *uint64             `json:"maxLength,omitempty"`
	Pattern              string              `json:"pattern,omitempty"`
	MinItems             *uint64             `json:"minItems,omitempty"`
	MaxItems             *uint64             `json:"maxItems,omitempty"`
	UniqueItems          bool                `json:"uniqueItems,omitempty"`
	MinProperties        *uint64             `json:"minProperties,omitempty"`
	MaxProperties        *uint64             `json:"maxProperties,omitempty"`
	AllOf                []*Schema           `json:"allOf,omitempty"`
	AnyOf                []*Schema           `json:"anyOf,omitempty"`
	OneOf                []*Schema           `json:"oneOf,omitempty"`
	Not                  *Schema             `json:"not,omitempty"`
	Nullable             bool                `json:"nullable,omitempty"`
	ReadOnly             bool                `json:"readOnly,omitempty"`
	WriteOnly            bool                `json:"writeOnly,omitempty"`
	Deprecated           bool                `json:"deprecated,omitempty"`
	ContentEncoding      string              `json:"contentEncoding,omitempty"`
	Ref                  string              `json:"$ref,omitempty"`
	Dependencies         map[string][]string `json:"dependencies,omitempty"`
}

// HasValidation returns true if at least one validator is set on the schema.
// This excludes the schema's type but includes most other fields and can be
// used to trigger additional slow validation steps when needed.
func (s *Schema) HasValidation() bool {
	if s.Items != nil || len(s.Properties) > 0 || s.AdditionalProperties != nil || len(s.PatternProperties) > 0 || len(s.Enum) > 0 || s.Minimum != nil || s.ExclusiveMinimum != nil || s.Maximum != nil || s.ExclusiveMaximum != nil || s.MultipleOf != 0 || s.MinLength != nil || s.MaxLength != nil || s.Pattern != "" || s.MinItems != nil || s.MaxItems != nil || s.UniqueItems || s.MinProperties != nil || s.MaxProperties != nil || len(s.AllOf) > 0 || len(s.AnyOf) > 0 || len(s.OneOf) > 0 || s.Not != nil || s.Ref != "" {
		return true
	}

	return false
}

// Generate creates a JSON schema for a Go type. Struct field tags
// can be used to provide additional metadata such as descriptions and
// validation.
func Generate(t reflect.Type, schema *Schema) (*Schema, error) {
	return GenerateWithMode(t, ModeAll, schema)
}

// getFields performs a breadth-first search for all fields including embedded
// ones. It may return multiple fields with the same name, the first of which
// represents the outer-most declaration.
func getFields(typ reflect.Type) []reflect.StructField {
	fields := make([]reflect.StructField, 0, typ.NumField())
	embedded := []reflect.StructField{}

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Anonymous {
			embedded = append(embedded, f)
			continue
		}

		fields = append(fields, f)
	}

	for _, f := range embedded {
		newTyp := f.Type
		if newTyp.Kind() == reflect.Ptr {
			newTyp = newTyp.Elem()
		}
		if newTyp.Kind() == reflect.Struct {
			fields = append(fields, getFields(newTyp)...)
		}
	}

	return fields
}

// generateFromField generates a schema for a single struct field. It returns
// the computed field name, whether it is optional, its schema, and any error
// which may have occurred.
func generateFromField(f reflect.StructField, mode Mode) (string, *Schema, error) {
	jsonTags := strings.Split(f.Tag.Get("json"), ",")
	name := strings.ToLower(f.Name)
	if len(jsonTags) > 0 && jsonTags[0] != "" {
		name = jsonTags[0]
	}

	if name == "-" {
		// Skip deliberately filtered out items
		return name, nil, nil
	}

	schema := &Schema{Name: name}

	if tag, ok := f.Tag.Lookup("required"); ok {
		if !(tag == "true" || tag == "false") {
			return name, nil, fmt.Errorf("%s required: boolean should be true or false: %w", f.Name, ErrSchemaInvalid)
		}
		schema.Required = tag == "true"
	}

	//找到dependencies并且生成
	if tag, ok := f.Tag.Lookup("dependencies"); ok {
		attrList := strings.Split(tag, " ")
		dependencies := make(map[string][]string)

		for _, attrs := range attrList {
			idx := strings.Index(attrs, ":")
			if idx == -1 {
				return name, nil, fmt.Errorf("Create Json Schema Fail. StructField %s: dependencies tag format err: %s. ", name, tag)
			}
			key := attrs[:idx]
			dps := strings.Split(attrs[idx+1:], ";")

			for _, dp := range dps {
				if dp == "" {
					return name, nil, fmt.Errorf("Create Json Schema Fail. StructField %s: dependencies tag format err: %s. ", name, tag)
				}
			}
			dependencies[key] = dps
		}
		schema.Dependencies = dependencies
	}

	//生成field 类型的对应schema
	s, err := GenerateWithMode(f.Type, mode, schema)
	if err != nil {
		return name, nil, err
	}

	if tag, ok := f.Tag.Lookup("description"); ok {
		s.Description = tag
	}

	if tag, ok := f.Tag.Lookup("doc"); ok {
		s.Description = tag
	}

	if tag, ok := f.Tag.Lookup("format"); ok {
		s.Format = tag
	}

	if tag, ok := f.Tag.Lookup("enum"); ok {
		s.Enum = []interface{}{}

		enumType := f.Type
		enumSchema := s
		if s.Type == TypeArray {
			// Enum values should be the type of the array elements, not the
			// array itself!
			enumType = f.Type.Elem()
			enumSchema = s.Items
		}

		for _, v := range strings.Split(tag, ",") {
			parsed, err := getTagValue(enumSchema, enumType, v)
			if err != nil {
				return name, nil, err
			}

			enumSchema.Enum = append(enumSchema.Enum, parsed)
		}
	}

	if tag, ok := f.Tag.Lookup("default"); ok {
		v, err := getTagValue(s, f.Type, tag)
		if err != nil {
			return name, nil, err
		}

		s.Default = v
	}

	if tag, ok := f.Tag.Lookup("example"); ok {
		v, err := getTagValue(s, f.Type, tag)
		if err != nil {
			return name, nil, err
		}

		s.Example = v
	}

	if tag, ok := f.Tag.Lookup("minimum"); ok {
		min, err := strconv.ParseFloat(tag, 64)
		if err != nil {
			return name, nil, err
		}
		s.Minimum = &min
	}

	if tag, ok := f.Tag.Lookup("exclusiveMinimum"); ok {
		min, err := strconv.ParseFloat(tag, 64)
		if err != nil {
			return name, nil, err
		}
		s.Minimum = &min
		t := true
		s.ExclusiveMinimum = &t
	}

	if tag, ok := f.Tag.Lookup("maximum"); ok {
		max, err := strconv.ParseFloat(tag, 64)
		if err != nil {
			return name, nil, err
		}
		s.Maximum = &max
	}

	if tag, ok := f.Tag.Lookup("exclusiveMaximum"); ok {
		max, err := strconv.ParseFloat(tag, 64)
		if err != nil {
			return name, nil, err
		}
		s.Maximum = &max
		t := true
		s.ExclusiveMaximum = &t
	}

	if tag, ok := f.Tag.Lookup("multipleOf"); ok {
		mof, err := strconv.ParseFloat(tag, 64)
		if err != nil {
			return name, nil, err
		}
		s.MultipleOf = mof
	}

	if tag, ok := f.Tag.Lookup("minLength"); ok {
		min, err := strconv.ParseUint(tag, 10, 64)
		if err != nil {
			return name, nil, err
		}
		s.MinLength = &min
	}

	if tag, ok := f.Tag.Lookup("maxLength"); ok {
		max, err := strconv.ParseUint(tag, 10, 64)
		if err != nil {
			return name, nil, err
		}
		s.MaxLength = &max
	}

	if tag, ok := f.Tag.Lookup("pattern"); ok {
		s.Pattern = tag

		if _, err := regexp.Compile(s.Pattern); err != nil {
			return name, nil, err
		}
	}

	if tag, ok := f.Tag.Lookup("minItems"); ok {
		min, err := strconv.ParseUint(tag, 10, 64)
		if err != nil {
			return name, nil, err
		}
		s.MinItems = &min
	}

	if tag, ok := f.Tag.Lookup("maxItems"); ok {
		max, err := strconv.ParseUint(tag, 10, 64)
		if err != nil {
			return name, nil, err
		}
		s.MaxItems = &max
	}

	if tag, ok := f.Tag.Lookup("uniqueItems"); ok {
		if !(tag == "true" || tag == "false") {
			return name, nil, fmt.Errorf("%s uniqueItems: boolean should be true or false: %w", f.Name, ErrSchemaInvalid)
		}
		s.UniqueItems = tag == "true"
	}

	if tag, ok := f.Tag.Lookup("minProperties"); ok {
		min, err := strconv.ParseUint(tag, 10, 64)
		if err != nil {
			return name, nil, err
		}
		s.MinProperties = &min
	}

	if tag, ok := f.Tag.Lookup("maxProperties"); ok {
		max, err := strconv.ParseUint(tag, 10, 64)
		if err != nil {
			return name, nil, err
		}
		s.MaxProperties = &max
	}

	if tag, ok := f.Tag.Lookup("nullable"); ok {
		if !(tag == "true" || tag == "false") {
			return name, nil, fmt.Errorf("%s nullable: boolean should be true or false but got %s: %w", f.Name, tag, ErrSchemaInvalid)
		}
		s.Nullable = tag == "true"
	}

	if tag, ok := f.Tag.Lookup("readOnly"); ok {
		if !(tag == "true" || tag == "false") {
			return name, nil, fmt.Errorf("%s readOnly: boolean should be true or false: %w", f.Name, ErrSchemaInvalid)
		}
		s.ReadOnly = tag == "true"
	}

	if tag, ok := f.Tag.Lookup("writeOnly"); ok {
		if !(tag == "true" || tag == "false") {
			return name, nil, fmt.Errorf("%s writeOnly: boolean should be true or false: %w", f.Name, ErrSchemaInvalid)
		}
		s.WriteOnly = tag == "true"
	}

	if tag, ok := f.Tag.Lookup("deprecated"); ok {
		if !(tag == "true" || tag == "false") {
			return name, nil, fmt.Errorf("%s deprecated: boolean should be true or false: %w", f.Name, ErrSchemaInvalid)
		}
		s.Deprecated = tag == "true"
	}

	return name, s, nil
}

// GenerateWithMode creates a JSON schema for a Go type. Struct field
// tags can be used to provide additional metadata such as descriptions and
// validation. The mode can be all, read, or write. In read or write mode
// any field that is marked as the opposite will be excluded, e.g. a
// write-only field would not be included in read mode. If a schema is given
// as input, add to it, otherwise creates a new schema.
func GenerateWithMode(t reflect.Type, mode Mode, schema *Schema) (*Schema, error) {
	if schema == nil {
		schema = &Schema{}
	}

	if t == ipType {
		// Special case: IP address.
		return &Schema{Type: TypeString, Format: "ipv4"}, nil
	}

	switch t.Kind() {
	case reflect.Struct:
		// Handle special cases.
		switch t {
		case timeType:
			return &Schema{Type: TypeString, Format: "date-time"}, nil
		case uriType:
			return &Schema{Type: TypeString, Format: "uri"}, nil
		}

		properties := make([]*Schema, 0)
		propertiesSet := make(map[string]struct{})
		schema.Type = TypeObject

		for _, f := range getFields(t) {
			name, s, err := generateFromField(f, mode)
			if err != nil {
				return nil, err
			}

			if s == nil {
				// Skip deliberately filtered out items
				continue
			}

			if _, ok := propertiesSet[name]; ok {
				// Item already exists, ignore it since we process embedded fields
				// after top-level ones.
				continue
			}

			if s.ReadOnly && mode == ModeWrite {
				continue
			}

			if s.WriteOnly && mode == ModeRead {
				continue
			}

			properties = append(properties, s)

		}

		if len(properties) > 0 {
			schema.Properties = properties
		}

		//判断scheme的dependencies存不存在，存在则校验里面的key及其依赖在properties里存在
		dependencies := schema.Dependencies
		if dependencies != nil {
			for key, dps := range dependencies {
				if _, has := propertiesSet[key]; !has {
					return nil, fmt.Errorf("Create Json Schema Fail: dependencies key:%s is not exist in properties. ", key)
				}
				for _, dp := range dps {
					if _, has := propertiesSet[dp]; !has {
						return nil, fmt.Errorf("Create Json Schema Fail: dependencies key:%s is not exist in properties. ", dp)
					}
				}
			}
		}

	case reflect.Map:
		schema.Type = TypeMap
		s, err := GenerateWithMode(t.Elem(), mode, nil)
		if err != nil {
			return nil, err
		}
		schema.Items = s
	case reflect.Slice, reflect.Array:
		if t.Elem().Kind() == reflect.Uint8 {
			// Special case: `[]byte` should be a Base-64 string.
			schema.Type = TypeString
		} else {
			schema.Type = TypeArray
			s, err := GenerateWithMode(t.Elem(), mode, nil)
			if err != nil {
				return nil, err
			}
			schema.Items = s
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		schema.Type = TypeInteger
		schema.Format = "int32"
	case reflect.Int64:
		schema.Type = TypeInteger
		schema.Format = "int64"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		// Unsigned integers can't be negative.
		schema.Type = TypeInteger
		schema.Format = "int32"
		schema.Minimum = F(0.0)
	case reflect.Uint64:
		schema.Type = TypeInteger
		schema.Format = "int64"
		schema.Minimum = F(0.0)
	case reflect.Float32:
		schema.Type = TypeNumber
		schema.Format = "float"
	case reflect.Float64:
		schema.Type = TypeNumber
		schema.Format = "double"
	case reflect.Bool:
		schema.Type = TypeBoolean
	case reflect.String:
		schema.Type = TypeString
	case reflect.Ptr:
		return GenerateWithMode(t.Elem(), mode, schema)
	case reflect.Interface:
		// Interfaces can be any type.
	case reflect.Uintptr, reflect.UnsafePointer, reflect.Func:
		// Ignored...
	default:
		return nil, fmt.Errorf("unsupported type %s from %s", t.Kind(), t)
	}

	return schema, nil
}
