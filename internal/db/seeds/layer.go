package seeds

import (
	"geoserver/internal/db/models"
	"geoserver/internal/parser"
	"math"
	"os"
	"path"
	"strings"

	"github.com/Lvov-SA/gdal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Layers(db *gorm.DB) error {
	files, err := os.ReadDir("../resource/map")
	if err != nil {
		return err
	}
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
		projection, err := parser.ExtractProjectionFromWKT(dataset.Projection())
		layer.Projection = projection

		layer.UpperLeftX, layer.UpperLeftY, layer.LowerRightX, layer.LowerRightY, err = parser.ExtractBounds(src)
		if err != nil {
			return err
		}
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
