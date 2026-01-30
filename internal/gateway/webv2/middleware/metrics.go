package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/isutare412/imageer/internal/gateway/metric"
)

func ObserveMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		before := time.Now()
		next.ServeHTTP(w, r)

		var status int
		if rec, ok := GetResponseRecord(ctx); ok {
			status = rec.Status
		}

		// NOTE: We use path template instead of r.URL.Path to reduce cardinality
		path, err := mux.CurrentRoute(r).GetPathTemplate()
		if err != nil {
			panic(fmt.Errorf("failed to get path template: %w", err))
		}

		metric.ObserveHTTPRequest(r.Method, path, status, time.Since(before))
	})
}
