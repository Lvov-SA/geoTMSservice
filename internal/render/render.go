package render

import (
	"errors"
	"fmt"
	"geoserver/internal/translator"
	"math"
	"os/exec"
	"strconv"

	"github.com/Lvov-SA/gdal"
)

func CliWarpRender(task Task) error {
	minX, minY, maxX, maxY := translator.WebMercarator(task.x, task.y, task.z)

	intersects := !(maxX < task.layer.UpperLeftX || minX > task.layer.LowerRightX ||
		maxY < task.layer.LowerRightY || minY > task.layer.UpperLeftY)
	if !intersects {
		return errors.New("Ошибка границ слоя")
	}
	args := []string{
		"-s_srs", task.layer.Projection,
		"-t_srs", "EPSG:3857",
		"-te",
		fmt.Sprintf("%f", minX),
		fmt.Sprintf("%f", minY),
		fmt.Sprintf("%f", maxX),
		fmt.Sprintf("%f", maxY),
		"-ts",
		strconv.Itoa(task.layer.TileSize),
		strconv.Itoa(task.layer.TileSize),
		"--config", "GDAL_CACHEMAX", "512",
		"-wm", "2048",
		"-r", "near",
		"-of", "PNG",
		"-co", "COMPRESS=DEFLATE",
		"-co", "ZLEVEL=6",
		"-overwrite",
		"../resource/map/" + task.layer.SourcePath,
		task.filePath + task.fileName,
	}

	cmd := exec.Command("gdalwarp", args...)
	err := cmd.Run()
	return err
}

func TranslateRender(task Task) error {
	coef := math.Pow(2, float64(task.z))
	maxSize := min(task.layer.Width, task.layer.Height)
	xFloat := float64(task.x)
	yFloat := float64(task.y)
	maxSizeFloat := float64(maxSize)
	readSize := maxSizeFloat / coef
	if xFloat*readSize >= float64(task.layer.Width) || yFloat*readSize >= float64(task.layer.Height) {

		return errors.New("Выход за границы")
	}
	options := []string{"-srcwin",
		fmt.Sprintf("%d", int(xFloat*readSize)),
		fmt.Sprintf("%d", int(yFloat*readSize)),
		fmt.Sprintf("%d", int(readSize)),
		fmt.Sprintf("%d", int(readSize)),
		"-outsize",
		fmt.Sprintf("%d", task.layer.TileSize),
		fmt.Sprintf("%d", task.layer.TileSize)}
	err := gdal.ConvertTile(
		"../resource/map/"+task.layer.SourcePath,
		task.filePath+task.fileName,
		options,
	)
	return err
}
