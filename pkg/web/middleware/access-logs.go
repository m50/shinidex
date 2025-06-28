package middleware

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gookit/slog"
	"github.com/labstack/echo/v4"
	"github.com/m50/shinidex/pkg/context"
	"github.com/m50/shinidex/pkg/web/errors"
	"github.com/spf13/viper"
)

const (
	keyTimeTaken = "TimeTaken"
)

var (
	staticPath = regexp.MustCompile("^/(assets|static|icons|.well-known)")
)

func LoggingHandler() func(echo.HandlerFunc) echo.HandlerFunc {
	accessLogs := viper.GetBool("logs.access")
	staticAccessLogs := viper.GetBool("logs.static-access")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			e := next(c)
			if !accessLogs {
				return e
			}
			if !staticAccessLogs && staticPath.Match([]byte(c.Request().RequestURI)) {
				return e
			}
			status := c.Response().Status
			if apiErr, ok := e.(*errors.APIError); ok {
				status = apiErr.StatusCode
			}
			log := slog.
				WithContext(context.FromEcho(c)).
				AddData(slog.M{
					keyTimeTaken: fmt.Sprintf("%dms", time.Since(start).Milliseconds()),
				})
			if e != nil {
				log.Errorf("%s %s: %v %v", c.Request().Method, c.Request().RequestURI, status, e)
			} else {
				log.Infof("%s %s: %v", c.Request().Method, c.Request().RequestURI, status)
			}
			return e
		}
	}
}
