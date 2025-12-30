package contextbag

import (
	"context"

	"github.com/isutare412/imageer/internal/gateway/domain"
)

// Bag is a placeholder for storing cross-cutting concerns in context.
type Bag struct {
	// Passport holds the authentication and authorization information.
	Passport domain.Passport
}

// WithBag returns a new context with an empty Bag.
func WithBag(ctx context.Context) context.Context {
	return context.WithValue(ctx, Bag{}, &Bag{})
}

// BagFromContext retrieves the Bag from the context. If no Bag is found in
// the context, bool return value will be false.
func BagFromContext(ctx context.Context) (*Bag, bool) {
	bag, ok := ctx.Value(Bag{}).(*Bag)
	return bag, ok
}
