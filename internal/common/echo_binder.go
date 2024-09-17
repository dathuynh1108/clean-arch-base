package common

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
)

type EchoBinder struct {
	echo.DefaultBinder
}

func (b *EchoBinder) Bind(i any, c echo.Context) (err error) {
	if err := b.BindPathParams(c, i); err != nil {
		return err
	}

	// Only bind query parameters for GET/DELETE/HEAD to avoid unexpected behavior with destination struct binding from body.
	// For example a request URL `&id=1&lang=en` with body `{"id":100,"lang":"de"}` would lead to precedence issues.
	// The HTTP method check restores pre-v4.1.11 behavior to avoid these problems (see issue #1670)
	method := c.Request().Method
	if method == http.MethodGet || method == http.MethodDelete || method == http.MethodHead {
		if err = b.BindQueryParams(c, i); err != nil {
			return err
		}
	}

	return b.BindBody(c, i)
}

func (b *EchoBinder) BindBody(c echo.Context, value any) (err error) {
	req := c.Request()
	if req.ContentLength == 0 {
		return
	}

	err = b.DefaultBinder.BindBody(c, value)
	if err != nil {
		return
	}

	ctype := req.Header.Get(echo.HeaderContentType)
	if strings.HasPrefix(ctype, echo.MIMEMultipartForm) {
		multipartForm, err := c.MultipartForm()
		if err != nil {
			return err
		}

		if err := bindFile(value, c, multipartForm.File); err != nil {
			return err
		}
	}

	return nil
}

func bindFile(i interface{}, c echo.Context, files map[string][]*multipart.FileHeader) error {
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
