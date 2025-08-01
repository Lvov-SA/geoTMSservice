package loader

import (
	"fmt"
	"geoserver/internal/db"
	"geoserver/internal/db/models"
	"math"
	"os"
	"path"
	"strings"

	"github.com/Lvov-SA/gdal"
	"gorm.io/gorm/clause"
)

var Layers map[string]models.Layer

func GeoTiff() error {
	err := Load()
	if err != nil {
		return err
	}
	var layers []models.Layer
	Layers = make(map[string]models.Layer)

	db, err := db.GetConnection()
	if err != nil {
		return err
	}

	db.Find(&layers)
	for _, layer := range layers {

		src := "../resource/map/" + layer.SourcePath
		dataset, err := gdal.Open(src, gdal.ReadOnly)
		if err != nil {
			fmt.Println("Ошибка загрузки файла " + layer.SourcePath)
			continue
		}
		dataset.Close()
		Layers[layer.Name] = layer
	}
	return nil
}

func Load() error {
	files, err := os.ReadDir("../resource/map")
	if err != nil {
		return err
	}
	db, _ := db.GetConnection()
	for _, file := range files {
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}
		var layer models.Layer
		src := "../resource/map/" + file.Name()
		dataset, err := gdal.Open(src, gdal.ReadOnly)
		if err != nil {
			return err
		}

		layer.Name = file.Name()[:len(file.Name())-len(path.Ext(file.Name()))]
		layer.Title = file.Name()[:len(file.Name())-len(path.Ext(file.Name()))]
		layer.SourcePath = file.Name()
		layer.Width = dataset.RasterXSize()
		layer.Height = dataset.RasterYSize()
		layer.IsActive = true
		layer.TileSize = 256
		layer.MinZoom = 0
		layer.MaxZoom = int(math.Log2(float64(min(layer.Width, layer.Height))))
		db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			UpdateAll: true,
		}).Create(&layer)
	}
	return nil
}
