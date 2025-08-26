package render

import (
	"errors"
	"geoserver/internal/db/models"
	"image"
	"os"
	"strconv"
	"sync"
	"time"
)

func Tiler(layer models.Layer, z, x, y int) (image.Image, error) {
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
			wg:       nil,
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
		file.Close()
		return nil, err
	}
	file.Close()
	return imageRGBA, nil
}

func MakeTask(layer models.Layer, z, x, y int, resultChan chan Result, wg *sync.WaitGroup) {
	filePath := "../resource/cache/" + layer.Name + "/" + strconv.Itoa(z) + "/"
	fileName := strconv.Itoa(x) + "_" + strconv.Itoa(y) + ".png"
	Tasks <- Task{
		layer:    layer,
		filePath: filePath,
		fileName: fileName,
		x:        x,
		y:        y,
		z:        z,
		result:   resultChan,
		wg:       wg,
	}
}
