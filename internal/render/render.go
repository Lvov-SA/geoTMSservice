package render

import (
	"errors"
	"geoserver/internal/db/models"
	"image"
	"os"
	"strconv"
	"time"
)

func CliRender(layer models.Layer, z, x, y int) (image.Image, error) {
	filePath := "../resource/cache/" + layer.Name + "/" + strconv.Itoa(z) + "/"
	fileName := strconv.Itoa(x) + "_" + strconv.Itoa(y) + ".png"
	file, err := os.Open(filePath + fileName)
	if os.IsNotExist(err) {
		resultChan := make(chan Result, 1)
		Tasks <- Task{
			layer:    layer,
			filePath: filePath,
			fileName: fileName,
			x:        x,
			y:        y,
			z:        z,
			result:   resultChan,
		}
		select {
		case result := <-resultChan:
			{
				if result.err != nil {
					return nil, err
				}
			}
		case <-time.After(10 * time.Minute):
			return nil, errors.New("Истекло время ожидания воркера")
		}
		file, err = os.Open(filePath + fileName)
		if err != nil {
			return nil, err
		}
	}

	imageRGBA, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	file.Close()
	return imageRGBA, nil
}
