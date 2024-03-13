package binder

import (
	"encoding"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type (
	// Binder is the interface that wraps the Bind method.
	Binder interface {
		Bind(i interface{}, r *fiber.Ctx) error
	}

	// DefaultBinder is the default implementation of the Binder interface.
	DefaultBinder struct{}

	// BindUnmarshaler is the interface used to wrap the UnmarshalParam method.
	// Types that don't implement this, but do implement encoding.TextUnmarshaler
	// will use that interface instead.
	BindUnmarshaler interface {
		// UnmarshalParam decodes and assigns a value from an form or query param.
		UnmarshalParam(param string) error
	}
)

// BindPathParams binds path params to bindable object
func (b *DefaultBinder) BindPathParams(ctx *fiber.Ctx, i interface{}) error {
	names := ctx.Route().Params
	params := map[string][]string{}
	for _, name := range names {
		params[name] = []string{ctx.Params(name)}
	}
	if err := b.bindData(i, params, "param"); err != nil {
		return err
	}
	return nil
}

// BindQueryParams binds query params to bindable object
func (b *DefaultBinder) BindQueryParams(ctx *fiber.Ctx, i interface{}) error {
	return ctx.QueryParser(i)
}

// BindBody binds request body contents to bindable object
// NB: then binding forms take note that this implementation uses standard library form parsing
// which parses form data from BOTH URL and BODY if content type is not MIMEMultipartForm
// See non-MIMEMultipartForm: https://golang.org/pkg/net/http/#Request.ParseForm
// See MIMEMultipartForm: https://golang.org/pkg/net/http/#Request.ParseMultipartForm
func (b *DefaultBinder) BindBody(ctx *fiber.Ctx, i interface{}) (err error) {
	return ctx.BodyParser(i)
}

// BindHeaders binds HTTP headers to a bindable object
func (b *DefaultBinder) BindHeaders(ctx *fiber.Ctx, i interface{}) error {
	if err := b.bindData(i, ctx.GetReqHeaders(), "header"); err != nil {
		return err
	}
	return nil
}

// Bind implements the `Binder#Bind` function.
// Binding is done in following order: 1) path params; 2) query params; 3) request body. Each step COULD override previous
// step binded values. For single source binding use their own methods BindBody, BindQueryParams, BindPathParams.
func (b *DefaultBinder) Bind(i interface{}, ctx *fiber.Ctx) (err error) {
	if err := b.BindPathParams(ctx, i); err != nil {
		return err
	}
	// Only bind query parameters for GET/DELETE/HEAD to avoid unexpected behavior with destination struct binding from body.
	// For example a request URL `&id=1&lang=en` with body `{"id":100,"lang":"de"}` would lead to precedence issues.
	// The HTTP method check restores pre-v4.1.11 behavior to avoid these problems (see issue #1670)
	method := ctx.Method()
	if method == http.MethodGet || method == http.MethodDelete || method == http.MethodHead {
		if err = b.BindQueryParams(ctx, i); err != nil {
			return err
		}
	}
	return b.BindBody(ctx, i)
}

// bindData will bind data ONLY fields in destination struct that have EXPLICIT tag
func (b *DefaultBinder) bindData(destination interface{}, data map[string][]string, tag string) error {
	if destination == nil || len(data) == 0 {
		return nil
	}
	typ := reflect.TypeOf(destination).Elem()
	val := reflect.ValueOf(destination).Elem()

	// Map
	if typ.Kind() == reflect.Map {
		for k, v := range data {
			val.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v[0]))
		}
		return nil
	}

	// !struct
	if typ.Kind() != reflect.Struct {
		if tag == "param" || tag == "query" || tag == "header" {
			// incompatible type, data is probably to be found in the body
			return nil
		}
		return errors.New("binding element must be a struct")
	}

	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		if typeField.Anonymous {
			if structField.Kind() == reflect.Ptr {
				structField = structField.Elem()
			}
		}
		if !structField.CanSet() {
			continue
		}
		structFieldKind := structField.Kind()
		inputFieldName := typeField.Tag.Get(tag)
		if typeField.Anonymous && structField.Kind() == reflect.Struct && inputFieldName != "" {
			// if anonymous struct with query/param/form tags, report an error
			return errors.New("query/param/form tags are not allowed with anonymous struct field")
		}

		if inputFieldName == "" {
			// If tag is nil, we inspect if the field is a not BindUnmarshaler struct and try to bind data into it (might contains fields with tags).
			// structs that implement BindUnmarshaler are binded only when they have explicit tag
			if _, ok := structField.Addr().Interface().(BindUnmarshaler); !ok && structFieldKind == reflect.Struct {
				if err := b.bindData(structField.Addr().Interface(), data, tag); err != nil {
					return err
				}
			}
			// does not have explicit tag and is not an ordinary struct - so move to next field
			continue
		}

		inputValue, exists := data[inputFieldName]
		if !exists {
			// Go json.Unmarshal supports case insensitive binding.  However the
			// url params are bound case sensitive which is inconsistent.  To
			// fix this we must check all of the map values in a
			// case-insensitive search.
			for k, v := range data {
				if strings.EqualFold(k, inputFieldName) {
					inputValue = v
					exists = true
					break
				}
			}
		}

		if !exists {
			continue
		}

		// Call this first, in case we're dealing with an alias to an array type
		if ok, err := unmarshalField(typeField.Type.Kind(), inputValue[0], structField); ok {
			if err != nil {
				return err
			}
			continue
		}

		numElems := len(inputValue)
		if structFieldKind == reflect.Slice && numElems > 0 {
			sliceOf := structField.Type().Elem().Kind()
			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
			for j := 0; j < numElems; j++ {
				if err := setWithProperType(sliceOf, inputValue[j], slice.Index(j)); err != nil {
					return err
				}
			}
			val.Field(i).Set(slice)
		} else if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil {
			return err

		}
	}
	return nil
}

func bindFile(i interface{}, r *http.Request, files map[string][]*multipart.FileHeader) error {
	iValue := reflect.Indirect(reflect.ValueOf(i))
	// check bind type is struct pointer
	if iValue.Kind() != reflect.Struct {
		return fmt.Errorf("bindFile input not is struct pointer, indirect type is %s", iValue.Type().String())
	}

	iType := iValue.Type()
	for i := 0; i < iType.NumField(); i++ {
		fType := iType.Field(i)
		// check canset field
		fValue := iValue.Field(i)
		if !fValue.CanSet() {
			continue
		}
		// revc type must *multipart.FileHeader or []*multipart.FileHeader
		switch fType.Type {
		case reflect.TypeOf((*multipart.FileHeader)(nil)):
			file := getFiles(files, fType.Name, fType.Tag.Get("form"))
			if len(file) > 0 {
				fValue.Set(reflect.ValueOf(file[0]))
			}
		case reflect.TypeOf(([]*multipart.FileHeader)(nil)):
			file := getFiles(files, fType.Name, fType.Tag.Get("form"))
			if len(file) > 0 {
				fValue.Set(reflect.ValueOf(file))
			}
		}
	}
	return nil
}

func getFiles(files map[string][]*multipart.FileHeader, names ...string) []*multipart.FileHeader {
	for _, name := range names {
		file, ok := files[name]
		if ok {
			return file
		}
	}
	return nil
}

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	// But also call it here, in case we're dealing with an array of BindUnmarshalers
	if ok, err := unmarshalField(valueKind, val, structField); ok {
		return err
	}

	switch valueKind {
	case reflect.Ptr:
		return setWithProperType(structField.Elem().Kind(), val, structField.Elem())
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32:
		return setFloatField(val, 32, structField)
	case reflect.Float64:
		return setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	default:
		return errors.New("unknown type")
	}
	return nil
}

func unmarshalField(valueKind reflect.Kind, val string, field reflect.Value) (bool, error) {
	switch valueKind {
	case reflect.Ptr:
		return unmarshalFieldPtr(val, field)
	default:
		return unmarshalFieldNonPtr(val, field)
	}
}

func unmarshalFieldNonPtr(value string, field reflect.Value) (bool, error) {
	fieldIValue := field.Addr().Interface()
	if unmarshaler, ok := fieldIValue.(BindUnmarshaler); ok {
		return true, unmarshaler.UnmarshalParam(value)
	}
	if unmarshaler, ok := fieldIValue.(encoding.TextUnmarshaler); ok {
		return true, unmarshaler.UnmarshalText([]byte(value))
	}

	return false, nil
}

func unmarshalFieldPtr(value string, field reflect.Value) (bool, error) {
	if field.IsNil() {
		// Initialize the pointer to a nil value
		field.Set(reflect.New(field.Type().Elem()))
	}
	return unmarshalFieldNonPtr(value, field.Elem())
}

func setIntField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	intVal, err := strconv.ParseInt(value, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	uintVal, err := strconv.ParseUint(value, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(value string, field reflect.Value) error {
	if value == "" {
		value = "false"
	}
	boolVal, err := strconv.ParseBool(value)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0.0"
	}
	floatVal, err := strconv.ParseFloat(value, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}
