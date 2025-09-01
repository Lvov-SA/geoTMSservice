package loader

import (
	"fmt"
	"geoserver/internal/db"
	"geoserver/internal/db/models"

	"github.com/Lvov-SA/gdal"
)

var Layers map[string]LayerGD

type LayerGD struct {
	Gd gdal.Dataset
	models.Layer
}

func GeoTiff() error {
	gdal.Init()
	var layers []models.Layer
	Layers = make(map[string]LayerGD)

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
		Layers[layer.Name] = LayerGD{Gd: dataset, Layer: layer}
	}
	return nil
}
