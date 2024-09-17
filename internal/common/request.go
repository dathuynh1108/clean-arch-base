package common

import (
	"bytes"
	"io"

	"github.com/dathuynh1108/clean-arch-base/pkg/comerr"
	"github.com/labstack/echo/v4"
)

func ReadContextRequestBody(c echo.Context) ([]byte, error) {
	request := c.Request()
	if request.Body == nil {
		return nil, nil
	}
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, comerr.WrapMessage(nil, "echo: read request body")
	}
	if err := request.Body.Close(); err != nil {
		return nil, comerr.WrapMessage(nil, "echo: close request body")
	}
	request.Body = io.NopCloser(bytes.NewBuffer(body))
	return body, nil
}
