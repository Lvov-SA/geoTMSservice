package route

import (
	"geoserver/internal/handlers"
	"geoserver/internal/middlware"
	"net/http"
)

func Init() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.IndexHandler)
	mux.HandleFunc("GET /{tile}/{z}/{x}/{y}", handlers.TileHandler)
	handler := middlware.Init(mux)

	handlers.CasheHandler = http.FileServer(http.Dir("../resource/cache"))

	return handler
}
