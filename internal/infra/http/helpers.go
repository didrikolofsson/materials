package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetPathParams(r *http.Request) map[string]string {
	rc := chi.RouteContext(r.Context())
	if rc == nil {
		return map[string]string{}
	}
	params := make(map[string]string, len(rc.URLParams.Keys))
	for i, key := range rc.URLParams.Keys {
		params[key] = rc.URLParams.Values[i]
	}
	return params
}
