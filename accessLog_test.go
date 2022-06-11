package accessLog

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.Use(LoggerWithWriter(buffer))
	router.GET("/example", func(c context.Context, ctx *app.RequestContext) {})
	router.POST("/example", func(c context.Context, ctx *app.RequestContext) {})
	router.PUT("/example", func(c context.Context, ctx *app.RequestContext) {})
	router.DELETE("/example", func(c context.Context, ctx *app.RequestContext) {})
	router.PATCH("/example", func(c context.Context, ctx *app.RequestContext) {})
	router.HEAD("/example", func(c context.Context, ctx *app.RequestContext) {})
	router.OPTIONS("/example", func(c context.Context, ctx *app.RequestContext) {})

	_ = ut.PerformRequest(router, "GET", "/example?a=100", nil)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")
	assert.Contains(t, buffer.String(), "a=100")

	buffer.Reset()
	_ = ut.PerformRequest(router, "POST", "/example", nil)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "POST")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	_ = ut.PerformRequest(router, "PUT", "/example", nil)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PUT")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	_ = ut.PerformRequest(router, "DELETE", "/example", nil)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "DELETE")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	_ = ut.PerformRequest(router, "PATCH", "/example", nil)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PATCH")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	_ = ut.PerformRequest(router, "HEAD", "/example", nil)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "HEAD")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	_ = ut.PerformRequest(router, "OPTIONS", "/example", nil)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "OPTIONS")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	_ = ut.PerformRequest(router, "GET", "/notfound", nil)
	assert.Contains(t, buffer.String(), "404")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/notfound")
}

func TestLoggerWithConfig(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.Use(LoggerWithConfig(LoggerConfig{Output: buffer}))
	router.GET("/example", func(c context.Context, ctx *app.RequestContext) {})
	router.POST("/example", func(c context.Context, ctx *app.RequestContext) {})
	router.PUT("/example", func(c context.Context, ctx *app.RequestContext) {})
	router.DELETE("/example", func(c context.Context, ctx *app.RequestContext) {})
	router.PATCH("/example", func(c context.Context, ctx *app.RequestContext) {})
	router.HEAD("/example", func(c context.Context, ctx *app.RequestContext) {})
	router.OPTIONS("/example", func(c context.Context, ctx *app.RequestContext) {})

	_ = ut.PerformRequest(router, "GET", "/example?a=100", nil)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")
	assert.Contains(t, buffer.String(), "a=100")

	buffer.Reset()
	_ = ut.PerformRequest(router, "POST", "/example", nil)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "POST")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	_ = ut.PerformRequest(router, "PUT", "/example", nil)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PUT")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	_ = ut.PerformRequest(router, "DELETE", "/example", nil)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "DELETE")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	_ = ut.PerformRequest(router, "PATCH", "/example", nil)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "PATCH")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	_ = ut.PerformRequest(router, "HEAD", "/example", nil)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "HEAD")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	_ = ut.PerformRequest(router, "OPTIONS", "/example", nil)
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "OPTIONS")
	assert.Contains(t, buffer.String(), "/example")

	buffer.Reset()
	_ = ut.PerformRequest(router, "GET", "/notfound", nil)
	assert.Contains(t, buffer.String(), "404")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/notfound")
}

func TestLoggerWithConfigFormatting(t *testing.T) {
	var gotParam LogFormatterParams
	var gotKeys map[string]any
	buffer := new(bytes.Buffer)

	router := route.NewEngine(config.NewOptions([]config.Option{}))

	router.Use(LoggerWithConfig(LoggerConfig{
		Output: buffer,
		Formatter: func(param LogFormatterParams) string {
			// for assert test
			gotParam = param

			return fmt.Sprintf("[FORMATTER TEST] %v | %3d | %13v | %15s | %-7s %s\n%s",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				param.StatusCode,
				param.Latency,
				param.ClientIP,
				param.Method,
				param.Path,
				param.ErrorMessage,
			)
		},
	}))
	router.GET("/example", func(c context.Context, ctx *app.RequestContext) {
		// set dummy ClientIP
		ctx.Request.Header.Set("X-Forwarded-For", "20.20.20.20")
		gotKeys = ctx.Keys
		time.Sleep(time.Millisecond)
	})
	_ = ut.PerformRequest(router, "GET", "/example?a=100", nil)

	// output test
	assert.Contains(t, buffer.String(), "[FORMATTER TEST]")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")
	assert.Contains(t, buffer.String(), "a=100")

	// LogFormatterParams test
	assert.NotNil(t, gotParam.Request)
	assert.NotEmpty(t, gotParam.TimeStamp)
	assert.Equal(t, 200, gotParam.StatusCode)
	assert.NotEmpty(t, gotParam.Latency)
	assert.Equal(t, "20.20.20.20", gotParam.ClientIP)
	assert.Equal(t, "GET", gotParam.Method)
	assert.Equal(t, "/example?a=100", gotParam.Path)
	assert.Empty(t, gotParam.ErrorMessage)
	assert.Equal(t, gotKeys, gotParam.Keys)
}

func TestDefaultLogFormatter(t *testing.T) {
	timeStamp := time.Unix(1544173902, 0).UTC()

	termFalseParam := LogFormatterParams{
		TimeStamp:    timeStamp,
		StatusCode:   200,
		Latency:      time.Second * 5,
		ClientIP:     "20.20.20.20",
		Method:       "GET",
		Path:         "/",
		ErrorMessage: "",
		isTerm:       false,
	}

	termTrueParam := LogFormatterParams{
		TimeStamp:    timeStamp,
		StatusCode:   200,
		Latency:      time.Second * 5,
		ClientIP:     "20.20.20.20",
		Method:       "GET",
		Path:         "/",
		ErrorMessage: "",
		isTerm:       true,
	}
	termTrueLongDurationParam := LogFormatterParams{
		TimeStamp:    timeStamp,
		StatusCode:   200,
		Latency:      time.Millisecond * 9876543210,
		ClientIP:     "20.20.20.20",
		Method:       "GET",
		Path:         "/",
		ErrorMessage: "",
		isTerm:       true,
	}

	termFalseLongDurationParam := LogFormatterParams{
		TimeStamp:    timeStamp,
		StatusCode:   200,
		Latency:      time.Millisecond * 9876543210,
		ClientIP:     "20.20.20.20",
		Method:       "GET",
		Path:         "/",
		ErrorMessage: "",
		isTerm:       false,
	}

	assert.Equal(t, "[Hertz] 2018/12/07 - 09:11:42 | 200 |            5s |     20.20.20.20 | GET      \"/\"\n", defaultLogFormatter(termFalseParam))
	assert.Equal(t, "[Hertz] 2018/12/07 - 09:11:42 | 200 |    2743h29m3s |     20.20.20.20 | GET      \"/\"\n", defaultLogFormatter(termFalseLongDurationParam))

	assert.Equal(t, "[Hertz] 2018/12/07 - 09:11:42 |\x1b[97;42m 200 \x1b[0m|            5s |     20.20.20.20 |\x1b[97;44m GET     \x1b[0m \"/\"\n", defaultLogFormatter(termTrueParam))
	assert.Equal(t, "[Hertz] 2018/12/07 - 09:11:42 |\x1b[97;42m 200 \x1b[0m|    2743h29m3s |     20.20.20.20 |\x1b[97;44m GET     \x1b[0m \"/\"\n", defaultLogFormatter(termTrueLongDurationParam))
}

func TestColorForMethod(t *testing.T) {
	colorForMethod := func(method string) string {
		p := LogFormatterParams{
			Method: method,
		}
		return p.MethodColor()
	}

	assert.Equal(t, blue, colorForMethod("GET"), "get should be blue")
	assert.Equal(t, cyan, colorForMethod("POST"), "post should be cyan")
	assert.Equal(t, yellow, colorForMethod("PUT"), "put should be yellow")
	assert.Equal(t, red, colorForMethod("DELETE"), "delete should be red")
	assert.Equal(t, green, colorForMethod("PATCH"), "patch should be green")
	assert.Equal(t, magenta, colorForMethod("HEAD"), "head should be magenta")
	assert.Equal(t, white, colorForMethod("OPTIONS"), "options should be white")
	assert.Equal(t, reset, colorForMethod("TRACE"), "trace is not defined and should be the reset color")
}

func TestColorForStatus(t *testing.T) {
	colorForStatus := func(code int) string {
		p := LogFormatterParams{
			StatusCode: code,
		}
		return p.StatusCodeColor()
	}

	assert.Equal(t, green, colorForStatus(http.StatusOK), "2xx should be green")
	assert.Equal(t, white, colorForStatus(http.StatusMovedPermanently), "3xx should be white")
	assert.Equal(t, yellow, colorForStatus(http.StatusNotFound), "4xx should be yellow")
	assert.Equal(t, red, colorForStatus(2), "other things should be red")
}

func TestResetColor(t *testing.T) {
	p := LogFormatterParams{}
	assert.Equal(t, string([]byte{27, 91, 48, 109}), p.ResetColor())
}

func TestIsOutputColor(t *testing.T) {
	// test with isTerm flag true.
	p := LogFormatterParams{
		isTerm: true,
	}

	consoleColorMode = autoColor
	assert.Equal(t, true, p.IsOutputColor())

	ForceConsoleColor()
	assert.Equal(t, true, p.IsOutputColor())

	DisableConsoleColor()
	assert.Equal(t, false, p.IsOutputColor())

	// test with isTerm flag false.
	p = LogFormatterParams{
		isTerm: false,
	}

	consoleColorMode = autoColor
	assert.Equal(t, false, p.IsOutputColor())

	ForceConsoleColor()
	assert.Equal(t, true, p.IsOutputColor())

	DisableConsoleColor()
	assert.Equal(t, false, p.IsOutputColor())

	// reset console color mode.
	consoleColorMode = autoColor
}

func TestLoggerWithFormatter(t *testing.T) {
	buffer := new(bytes.Buffer)

	d := DefaultWriter
	DefaultWriter = buffer
	defer func() {
		DefaultWriter = d
	}()

	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.Use(LoggerWithFormatter(func(param LogFormatterParams) string {
		return fmt.Sprintf("[FORMATTER TEST] %v | %3d | %13v | %15s | %-7s %#v\n%s",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
			param.ErrorMessage,
		)
	}))
	router.GET("/example", func(c context.Context, ctx *app.RequestContext) {})

	_ = ut.PerformRequest(router, "GET", "/example?a=100", nil)

	// output test
	assert.Contains(t, buffer.String(), "[FORMATTER TEST]")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")
	assert.Contains(t, buffer.String(), "a=100")
}

func TestLoggerWithWriterSkippingPaths(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.Use(LoggerWithWriter(buffer, "/skipped"))
	router.GET("/logged", func(c context.Context, ctx *app.RequestContext) {})
	router.GET("/skipped", func(c context.Context, ctx *app.RequestContext) {})

	_ = ut.PerformRequest(router, "GET", "/logged", nil)
	assert.Contains(t, buffer.String(), "200")

	buffer.Reset()
	_ = ut.PerformRequest(router, "GET", "/skipped", nil)
	assert.Contains(t, buffer.String(), "")
}

func TestLoggerWithConfigSkippingPaths(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.Use(LoggerWithConfig(LoggerConfig{
		Output:    buffer,
		SkipPaths: []string{"/skipped"},
	}))
	router.GET("/logged", func(c context.Context, ctx *app.RequestContext) {})
	router.GET("/skipped", func(c context.Context, ctx *app.RequestContext) {})

	_ = ut.PerformRequest(router, "GET", "/logged", nil)
	assert.Contains(t, buffer.String(), "200")

	buffer.Reset()
	_ = ut.PerformRequest(router, "GET", "/skipped", nil)
	assert.Contains(t, buffer.String(), "")
}

func TestDisableConsoleColor(t *testing.T) {
	server.New()
	assert.Equal(t, autoColor, consoleColorMode)
	DisableConsoleColor()
	assert.Equal(t, disableColor, consoleColorMode)

	// reset console color mode.
	consoleColorMode = autoColor
}

func TestForceConsoleColor(t *testing.T) {
	server.New()
	assert.Equal(t, autoColor, consoleColorMode)
	ForceConsoleColor()
	assert.Equal(t, forceColor, consoleColorMode)

	// reset console color mode.
	consoleColorMode = autoColor
}
