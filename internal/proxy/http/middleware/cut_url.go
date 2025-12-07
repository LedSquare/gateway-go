package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type contextKey string

const NewRouteKey contextKey = "new_route"

func CutUrl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := chi.RouteContext(r.Context())

		if route == nil {
			http.Error(w, "контект роута не найден", http.StatusInternalServerError)
			return
		}

		path := r.URL.Path
		segments := strings.Split(strings.Trim(path, "/"), "/")
		if len(segments) < 2 {
			http.Error(w, "неправильный роут, не хватает названия сервиса", http.StatusInternalServerError)
		}

		version := segments[1]
		resourcePath := strings.Join(segments[2:], "/")
		newRoute := version + "/" + resourcePath

		ctx := context.WithValue(r.Context(), NewRouteKey, newRoute)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
