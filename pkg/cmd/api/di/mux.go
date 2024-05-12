package di

import "net/http"

func newServeMux() (*http.ServeMux, http.Handler) {
	mux := http.NewServeMux()
	return mux, mux
}
