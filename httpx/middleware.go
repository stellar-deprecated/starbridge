package httpx

import (
	"context"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stellar/go/support/log"
)

var routeRegexp = regexp.MustCompile("{([^:}]*):[^}]*}")

func newWrapResponseWriter(w http.ResponseWriter, r *http.Request) middleware.WrapResponseWriter {
	mw, ok := w.(middleware.WrapResponseWriter)
	if !ok {
		mw = middleware.NewWrapResponseWriter(w, r.ProtoMajor)
	}

	return mw
}

// loggerMiddleware logs http requests and resposnes to the logging subsytem of horizon.
func loggerMiddleware(serverMetrics *ServerMetrics) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			mw := newWrapResponseWriter(w, r)

			logger := log.WithField("req", middleware.GetReqID(ctx))
			ctx = log.Set(ctx, logger)

			then := time.Now()
			next.ServeHTTP(mw, r.WithContext(ctx))
			duration := time.Since(then)
			logEndOfRequest(ctx, r, serverMetrics.RequestDurationSummary, duration, mw)
		})
	}
}

// getClientData gets client data (name or version) from header or GET parameter
// (useful when not possible to set headers, like in EventStream).
func getClientData(r *http.Request, headerName string) string {
	value := r.Header.Get(headerName)
	if value != "" {
		return value
	}

	value = r.URL.Query().Get(headerName)
	if value == "" {
		value = "undefined"
	}

	return value
}

func remoteAddrIP(r *http.Request) string {
	// To support IPv6
	lastSemicolon := strings.LastIndex(r.RemoteAddr, ":")
	if lastSemicolon == -1 {
		return r.RemoteAddr
	} else {
		return r.RemoteAddr[0:lastSemicolon]
	}
}

// https://prometheus.io/docs/instrumenting/exposition_formats/
// label_value can be any sequence of UTF-8 characters, but the backslash (\),
// double-quote ("), and line feed (\n) characters have to be escaped as \\,
// \", and \n, respectively.
func sanitizeMetricRoute(routePattern string) string {
	route := routeRegexp.ReplaceAllString(routePattern, "{$1}")
	route = strings.ReplaceAll(route, "\\", "\\\\")
	route = strings.ReplaceAll(route, "\"", "\\\"")
	route = strings.ReplaceAll(route, "\n", "\\n")
	if route == "" {
		// Can be empty when request did not reach the final route (ex. blocked by
		// a middleware). More info: https://github.com/go-chi/chi/issues/270
		return "undefined"
	}
	return route
}

// Author: https://github.com/rliebz
// From: https://github.com/go-chi/chi/issues/270#issuecomment-479184559
// https://github.com/go-chi/chi/blob/master/LICENSE
func getRoutePattern(r *http.Request) string {
	rctx := chi.RouteContext(r.Context())
	if pattern := rctx.RoutePattern(); pattern != "" {
		// Pattern is already available
		return pattern
	}

	routePath := r.URL.Path
	if r.URL.RawPath != "" {
		routePath = r.URL.RawPath
	}

	tctx := chi.NewRouteContext()
	if !rctx.Routes.Match(tctx, r.Method, routePath) {
		return ""
	}

	// tctx has the updated pattern, since Match mutates it
	return tctx.RoutePattern()
}

func logEndOfRequest(ctx context.Context, r *http.Request, requestDurationSummary *prometheus.SummaryVec, duration time.Duration, mw middleware.WrapResponseWriter) {
	route := sanitizeMetricRoute(getRoutePattern(r))

	referer := r.Referer()
	if referer == "" {
		referer = r.Header.Get("Origin")
	}
	if referer == "" {
		referer = "undefined"
	}

	log.Ctx(ctx).WithFields(log.F{
		"bytes":           mw.BytesWritten(),
		"duration":        duration.Seconds(),
		"x_forwarder_for": r.Header.Get("X-Forwarded-For"),
		"host":            r.Host,
		"ip":              remoteAddrIP(r),
		"ip_port":         r.RemoteAddr,
		"method":          r.Method,
		"path":            r.URL.String(),
		"route":           route,
		"status":          mw.Status(),
		"referer":         referer,
	}).Info("Finished request")

	requestDurationSummary.With(prometheus.Labels{
		"status": strconv.FormatInt(int64(mw.Status()), 10),
		"route":  route,
		"method": r.Method,
	}).Observe(float64(duration.Seconds()))
}
