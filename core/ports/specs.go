package ports

import (
	"context"
	"net/http"
)

type SpecParser interface {
	Handler(ctx context.Context) (http.Handler, error)
}
