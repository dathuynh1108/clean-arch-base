package middleware

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	CompressSchemeGzip   = "gzip"
	CompressSchemeBrotli = "br"
	CompressSchemeZstd   = "zstd"

	CompressLevelDefault = 1
	CompressLevelFast    = 2
	CompressLevelBest    = 3
)

type (
	CompressConfig struct {
		Skipper     middleware.Skipper
		HandleError bool
		Level       int
	}
	CompressWriter interface {
		io.WriteCloser
		Flush() error
		// Reset reset state of writer and set writer write to w (write data after compress)
		Reset(io.Writer)
	}
)

func Compress() echo.MiddlewareFunc {
	return CompressWithConfig(CompressConfig{})
}

func CompressWithConfig(conf CompressConfig) echo.MiddlewareFunc {
	if conf.Skipper == nil {
		conf.Skipper = middleware.DefaultSkipper
	}
	if conf.Level == 0 {
		conf.Level = CompressLevelDefault
	}
	var (
		compressorPoolMap = map[string]*sync.Pool{
			CompressSchemeZstd:   compressPoolZstd(conf),
			CompressSchemeBrotli: compressPoolBrotli(conf),
			CompressSchemeGzip:   compressPoolGzip(conf),
		}
	)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if conf.Skipper(c) {
				return next(c)
			}
			var (
				res                = c.Response()
				respHeader         = res.Header()
				acceptEncodingsRaw = c.Request().Header.Values(echo.HeaderAcceptEncoding)
				acceptEncodings    = make([]string, 0)
			)

			if len(acceptEncodingsRaw) == 0 {
				return next(c)
			}

			for _, acceptEncodingRaw := range acceptEncodingsRaw {
				// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding
				// It can be a string of gzip, deflate, br, zstd --> need split by comma
				for _, acceptEncoding := range strings.Split(acceptEncodingRaw, ",") {
					acceptEncoding = strings.TrimSpace(acceptEncoding)
					if acceptEncoding != "" {
						acceptEncodings = append(acceptEncodings, acceptEncoding)
					}
				}
			}

			var selectedPool *sync.Pool
			for _, acceptEncoding := range acceptEncodings {
				if pool, ok := compressorPoolMap[acceptEncoding]; ok {
					respHeader.Set(echo.HeaderContentEncoding, acceptEncoding)
					selectedPool = pool
					break
				}
			}
			if selectedPool == nil {
				return next(c)
			}

			respHeader.Add(echo.HeaderVary, echo.HeaderAcceptEncoding)
			cw, ok := selectedPool.Get().(CompressWriter)
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, cw.(error).Error())
			}

			rw := res.Writer

			// Reset reset state of writer and set writer write to rw (res.Writer)
			cw.Reset(rw)

			defer func() {
				if res.Size == 0 {
					contentEncoding := respHeader.Get(echo.HeaderContentEncoding)
					if _, ok := compressorPoolMap[contentEncoding]; ok {
						respHeader.Del(echo.HeaderContentEncoding)
					}
					cw.Reset(io.Discard)
				}
				_ = cw.Close()
				// We have to reset response to it's pristine state when
				// nothing is written to body or error is returned.
				res.Writer = rw
				selectedPool.Put(cw)
			}()
			res.Writer = &compressResponseWriter{
				ResponseWriter: rw,
				cw:             cw,
			}
			err = next(c)
			if err != nil && conf.HandleError {
				c.Echo().HTTPErrorHandler(err, c)
				return nil
			}
			return
		}
	}
}

func compressPoolZstd(config CompressConfig) *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			var level zstd.EncoderLevel
			switch config.Level {
			case CompressLevelDefault:
				level = zstd.SpeedDefault
			case CompressLevelBest:
				level = zstd.SpeedBestCompression
			default:
				level = zstd.SpeedFastest
			}
			writer, err := zstd.NewWriter(io.Discard, zstd.WithEncoderLevel(level))
			if err != nil {
				panic(err)
			}
			return writer
		},
	}
}

func compressPoolGzip(config CompressConfig) *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			var gzipLevel int
			switch config.Level {
			case CompressLevelDefault:
				gzipLevel = gzip.DefaultCompression
			case CompressLevelBest:
				gzipLevel = gzip.BestCompression
			default:
				gzipLevel = gzip.BestSpeed
			}
			writer, err := gzip.NewWriterLevel(io.Discard, gzipLevel)
			if err != nil {
				panic(err)
			}
			return writer
		},
	}
}

func compressPoolBrotli(config CompressConfig) *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			var brotliLevel int
			switch config.Level {
			case CompressLevelDefault:
				brotliLevel = brotli.DefaultCompression
			case CompressLevelBest:
				brotliLevel = brotli.BestCompression
			default:
				brotliLevel = brotli.BestSpeed
			}
			return brotli.NewWriterLevel(io.Discard, brotliLevel)
		},
	}
}

type compressResponseWriter struct {
	http.ResponseWriter
	cw CompressWriter
}

func (w *compressResponseWriter) WriteHeader(code int) {
	if code == http.StatusNoContent {
		w.Header().Del(echo.HeaderContentEncoding)
	}
	w.Header().Del(echo.HeaderContentLength)
	w.ResponseWriter.WriteHeader(code)
}

func (w *compressResponseWriter) Write(b []byte) (int, error) {
	if w.Header().Get(echo.HeaderContentType) == "" {
		w.Header().Set(echo.HeaderContentType, http.DetectContentType(b))
	}

	return w.cw.Write(b)
}

func (w *compressResponseWriter) Flush() {
	_ = w.cw.Flush()
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *compressResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w *compressResponseWriter) Push(target string, opts *http.PushOptions) error {
	if p, ok := w.ResponseWriter.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return http.ErrNotSupported
}
