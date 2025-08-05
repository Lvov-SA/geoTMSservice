package render

import (
	"errors"
	"fmt"
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
		fmt.Println("Отправка задачи в канал")
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
		case <-time.After(5 * time.Second):
			return nil, errors.New("Истекло время ожидания воркера")
		}

	}

	imageRGBA, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	file.Close()
	return imageRGBA, nil
}
