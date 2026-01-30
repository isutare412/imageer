package middleware

import (
	"context"
	"net/http"

	"github.com/felixge/httpsnoop"
)

type responseRecordContextKey struct{}

type ResponseRecord struct {
	Status       int
	ResponseSize int
}

func WithResponseRecord(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Initialize response metrics and store in context
		metrics := &ResponseRecord{}
		ctx = context.WithValue(ctx, responseRecordContextKey{}, metrics)

		// Wrap the ResponseWriter to capture metrics
		w = httpsnoop.Wrap(w, httpsnoop.Hooks{
			WriteHeader: func(whf httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
				return func(code int) {
					metrics.Status = code
					whf(code)
				}
			},
			Write: func(wf httpsnoop.WriteFunc) httpsnoop.WriteFunc {
				return func(b []byte) (int, error) {
					metrics.ResponseSize += len(b)
					return wf(b)
				}
			},
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetResponseRecord(ctx context.Context) (*ResponseRecord, bool) {
	metrics, ok := ctx.Value(responseRecordContextKey{}).(*ResponseRecord)
	return metrics, ok
}
