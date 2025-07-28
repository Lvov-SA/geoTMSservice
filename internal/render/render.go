package render

import (
	"errors"
	"fmt"
	"geoserver/internal/db/models"
	"image"
	"math"
	"os"
	"os/exec"
)

func CliRender(layer models.Layer, z, x, y int) (image.Image, error) {
	coef := math.Pow(2, float64(z))
	maxSize := min(layer.Width, layer.Height)
	xFloat := float64(x)
	yFloat := float64(y)
	maxSizeFloat := float64(maxSize)
	readSize := maxSizeFloat / coef
	if xFloat*readSize >= float64(layer.Width) || yFloat*readSize >= float64(layer.Height) {
		return nil, errors.New("Выход за границы")
	}
	tmpFile, err := os.CreateTemp("../resource", "tile_*.png")
	if err != nil {
		return nil, err
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)
	defer os.Remove(tmpPath + ".aux.xml")
	cmd := exec.Command("gdal_translate", "-srcwin",
		fmt.Sprintf("%d", int(xFloat*readSize)),
		fmt.Sprintf("%d", int(yFloat*readSize)),
		fmt.Sprintf("%d", int(readSize)),
		fmt.Sprintf("%d", int(readSize)),
		"-outsize",
		fmt.Sprintf("%d", layer.TileSize),
		fmt.Sprintf("%d", layer.TileSize),
		"../resource/map/"+layer.SourcePath,
		tmpPath)
	cmd.Run()

	file, err := os.Open(tmpPath)
	if err != nil {
		return nil, err
	}
	imageRGBA, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return imageRGBA, nil
}
