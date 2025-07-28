package render

import (
	"errors"
	"fmt"
	"geoserver/internal/db/models"
	"image"
	"math"
	"os"
	"os/exec"
	"strconv"
)

func CliRender(layer models.Layer, z, x, y int) (image.Image, error) {
	filePath := "../resource/cache/" + layer.Name + "/" + strconv.Itoa(z) + "/"
	fileName := strconv.Itoa(x) + "_" + strconv.Itoa(y) + ".png"
	file, err := os.Open(filePath + fileName)
	if os.IsNotExist(err) {
		err = os.MkdirAll(filePath, 0755)
		if err != nil {
			return nil, err
		}
		file, err = os.Create(filePath + fileName)
		if err != nil {
			return nil, err
		}
		coef := math.Pow(2, float64(z))
		maxSize := min(layer.Width, layer.Height)
		xFloat := float64(x)
		yFloat := float64(y)
		maxSizeFloat := float64(maxSize)
		readSize := maxSizeFloat / coef
		if xFloat*readSize >= float64(layer.Width) || yFloat*readSize >= float64(layer.Height) {
			return nil, errors.New("Выход за границы")
		}
		cmd := exec.Command("gdal_translate", "-srcwin",
			fmt.Sprintf("%d", int(xFloat*readSize)),
			fmt.Sprintf("%d", int(yFloat*readSize)),
			fmt.Sprintf("%d", int(readSize)),
			fmt.Sprintf("%d", int(readSize)),
			"-outsize",
			fmt.Sprintf("%d", layer.TileSize),
			fmt.Sprintf("%d", layer.TileSize),
			"../resource/map/"+layer.SourcePath,
			filePath+fileName)
		err = cmd.Run()
		if err != nil {
			return nil, err
		}
	}

	imageRGBA, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return imageRGBA, nil
}
