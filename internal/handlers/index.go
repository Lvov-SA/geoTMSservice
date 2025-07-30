package handlers

import (
	"geoserver/internal/config"
	"html/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("../public/index.html")
	data := struct {
		Host string
		Port int
	}{Host: config.Configs.HOST, Port: config.Configs.APP_PORT}
	t.Execute(w, data)
	if err != nil {
		http.Error(w, "Ошибка загрузки main_page: "+err.Error(), http.StatusBadRequest)
		return
	}
}
