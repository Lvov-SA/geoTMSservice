package main

import (
	"flag"
	"fmt"
	"geoserver/internal/config"
	"geoserver/internal/db"
	"geoserver/internal/loader"
	"geoserver/internal/render"
	"math"
	"sync"
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

	render.InitWorkers()

	layerName := flag.String("layer", "tile", "Название слоя")
	flag.Parse()

	layer := loader.Layers[*layerName]
	totalTiles := 0
	for i := 0; i < layer.MaxZoom; i++ {
		wight, height := CalculateTiles(i)
		totalTiles = totalTiles + wight*height
	}
	fmt.Printf("Всего тайлов: %v", totalTiles)
	fmt.Println()
	i := 0
	for z := 0; z < layer.MaxZoom; z++ {
		wight, height := CalculateTiles(z)
		var wg sync.WaitGroup
		for x := 0; x < wight; x++ {
			for y := 0; y < height; y++ {
				totalTiles--
				i++
				wg.Add(1)
				render.MakeTask(layer, z, x, y, nil, &wg)
			}
		}
		fmt.Printf("Set in queqe all tiles for zoom %v", z)
		fmt.Println()
		wg.Wait()
		fmt.Printf("Total tales ready: %v", i)
		fmt.Printf(" Total tales lest: %v", totalTiles)
		fmt.Printf(" Zoom ready: %v", z)
		fmt.Println()
	}
}

func CalculateTiles(zoom int) (width, height int) {
	tiles := math.Pow(2, float64(zoom))
	return int(tiles), int(tiles)
}
