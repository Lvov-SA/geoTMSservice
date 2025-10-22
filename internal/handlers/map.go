package handlers

import (
	"geoserver/internal/loader"
	"geoserver/internal/render"
	"net/http"
	"strconv"
	"strings"
)

var CasheHandler http.Handler

func TileHandler(w http.ResponseWriter, r *http.Request) {
	tileModel, exist := loader.Layers[r.PathValue("tile")]
	if !exist {
		http.Error(w, "Invalid Tile parameter", http.StatusBadRequest)
		return
	}

	z, err := strconv.Atoi(r.PathValue("z"))
	if err != nil || z < tileModel.MinZoom || z > tileModel.MaxZoom {
		http.Error(w, "Invalid z parameter", http.StatusBadRequest)
		return
	}

	x, err := strconv.Atoi(r.PathValue("x"))
	if err != nil || x < 0 {
		http.Error(w, "Invalid x parameter", http.StatusBadRequest)
		return
	}
	parts := strings.Split(r.PathValue("y"), ".")
	y, err := strconv.Atoi(parts[0])
	if err != nil || y < 0 {
		http.Error(w, "Invalid y parameter", http.StatusBadRequest)
		return
	}
	ext := parts[1]
	if ext != tileModel.TileExt {
		http.Error(w, "Invalid request extention", http.StatusBadRequest)
		return
	}
	err = render.Tiler(tileModel, z, x, y)
	if err != nil {
		http.Error(w, "Ошибка генерации тайла: "+err.Error(), http.StatusBadRequest)
		return
	}
	CasheHandler.ServeHTTP(w, r)
}
