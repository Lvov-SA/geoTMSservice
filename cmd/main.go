package main

import (
	"fmt"
	"geoserver/internal/config"
	"geoserver/internal/db"
	"geoserver/internal/loader"
	"geoserver/internal/render"
	"geoserver/internal/route"
	"strconv"

	"log"
	"net/http"
)

func main() {

	err := config.Init()
	if err != nil {
		fmt.Printf("Ошибка инициализации конфигурации приложения: %v", err)
		return
	}
	err = db.Init()
	if err != nil {
		fmt.Printf("Ошибка инициализации базы данных: %v", err)
		return
	}
	err = loader.GeoTiff()
	if err != nil {
		fmt.Printf("Ошибка загрузки данных слоев: %v", err)
		return
	}

	for _, layer := range loader.Layers {
		defer layer.Gd.Close()
	}

	render.InitWorkers()

	handler := route.Init()

	appUrl := config.Configs.HOST + ":" + strconv.Itoa(config.Configs.APP_PORT)
	log.Println("Server started at " + appUrl)
	log.Println("Access example: http://" + appUrl + "/tile/0/0/0.png")
	log.Println("Look at map: http://" + appUrl)

	err = http.ListenAndServe("0.0.0.0:"+strconv.Itoa(config.Configs.APP_PORT), handler)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v", err)
		return
	}
}
