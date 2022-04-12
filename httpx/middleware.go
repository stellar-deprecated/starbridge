package httpx

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

func newWrapResponseWriter(w http.ResponseWriter, r *http.Request) middleware.WrapResponseWriter {
	mw, ok := w.(middleware.WrapResponseWriter)
	if !ok {
		mw = middleware.NewWrapResponseWriter(w, r.ProtoMajor)
	}

	return mw
}

// prometheusMiddleware gathers http requests data.
func prometheusMiddleware(serverMetrics *ServerMetrics) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mw := newWrapResponseWriter(w, r)

			then := time.Now()
			next.ServeHTTP(mw, r)
			duration := time.Since(then)

			serverMetrics.RequestDurationSummary.With(prometheus.Labels{
				"status": strconv.FormatInt(int64(mw.Status()), 10),
				"method": r.Method,
			}).Observe(float64(duration.Seconds()))
		})
	}
}
