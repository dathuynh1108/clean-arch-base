package middleware

import (
	"io"
	"net/http"
	"sync"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	DecompessorPoolMap = map[string]*sync.Pool{
		CompressSchemeZstd:   decompressPoolZstd(),
		CompressSchemeBrotli: decompressPoolBrotli(),
		CompressSchemeGzip:   decompressPoolGzip(),
	}
)

type DecompressConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper
}

type DecompessReader interface {
	io.ReadCloser
	Reset(io.Reader) error
}

func Decompress() echo.MiddlewareFunc {
	return DecompressWithConfig(DecompressConfig{})
}

func DecompressWithConfig(config DecompressConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			contentEncoding := c.Request().Header.Get(echo.HeaderContentEncoding)
			if contentEncoding == "" {
				return next(c)
			}

			selectedPool, ok := DecompessorPoolMap[contentEncoding]
			if !ok || selectedPool == nil {
				return next(c)
			}

			i := selectedPool.Get()
			dr, ok := i.(DecompessReader)
			if !ok || dr == nil {
				return echo.NewHTTPError(http.StatusInternalServerError, i.(error).Error())
			}
			defer selectedPool.Put(dr)

			b := c.Request().Body
			defer b.Close()

			if err := dr.Reset(b); err != nil {
				if err == io.EOF { //ignore if body is empty
					return next(c)
				}
				return err
			}

			defer dr.Close()

			c.Request().Body = dr

			return next(c)
		}
	}
}

func decompressPoolGzip() *sync.Pool {
	return &sync.Pool{New: func() interface{} { return new(gzip.Reader) }}
}

func decompressPoolZstd() *sync.Pool {
	return &sync.Pool{New: func() interface{} { return new(zstd.Decoder) }}
}

func decompressPoolBrotli() *sync.Pool {
	return &sync.Pool{New: func() interface{} { return new(brotli.Reader) }}
}
