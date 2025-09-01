package render

import (
	"errors"
	"fmt"
	"geoserver/internal/translator"
	"os"
	"os/exec"
	"strconv"

	"github.com/Lvov-SA/gdal"
)

func CliWarpRender(task Task) error {
	var minX, minY, maxX, maxY float64

	switch task.layer.Projection {
	//WebMercarator
	case "EPSG:3857":
		minX, minY, maxX, maxY = translator.WebMercarator(task.x, task.y, task.z)
		intersects := !(maxX < task.layer.UpperLeftX || minX > task.layer.LowerRightX ||
			maxY < task.layer.LowerRightY || minY > task.layer.UpperLeftY)
		if !intersects {
			return errors.New("Ошибка границ слоя")
		}
	case "EPSG:4326":
		minX, minY, maxX, maxY = translator.WGS84(task.x, task.y, task.z)
		intersects := !(maxX < task.layer.UpperLeftX || minX > task.layer.LowerRightX ||
			maxY < task.layer.LowerRightY || minY > task.layer.UpperLeftY)
		if !intersects {
			return errors.New("Ошибка границ слоя")
		}
	default:
	}
	err := os.MkdirAll(task.filePath, 0755)
	if err != nil {
		return err
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
	err = cmd.Run()
	return err
}

func WarpRender(task Task) error {
	var minX, minY, maxX, maxY float64

	switch task.layer.Projection {
	//WebMercarator
	case "EPSG:3857":
		minX, minY, maxX, maxY = translator.WebMercarator(task.x, task.y, task.z)
		intersects := !(maxX < task.layer.UpperLeftX || minX > task.layer.LowerRightX ||
			maxY < task.layer.LowerRightY || minY > task.layer.UpperLeftY)
		if !intersects {
			return errors.New("Ошибка границ слоя")
		}
	case "EPSG:4326":
		minX, minY, maxX, maxY = translator.WGS84(task.x, task.y, task.z)
		intersects := !(maxX < task.layer.UpperLeftX || minX > task.layer.LowerRightX ||
			maxY < task.layer.LowerRightY || minY > task.layer.UpperLeftY)
		if !intersects {
			return errors.New("Ошибка границ слоя")
		}
	default:
	}
	err := os.MkdirAll(task.filePath, 0755)
	if err != nil {
		return err
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
		"-co", "BIGTIFF=YES",
		"-wm", "2048",
		"-r", "near",
		"-of", "PNG",
		"-co", "ZLEVEL=1",
		"-overwrite",
	}
	ds, err := gdal.Warp(task.filePath+task.fileName, nil, []gdal.Dataset{task.layer.Gd}, args)
	if ds.RasterXSize() == 0 || ds.RasterYSize() == 0 {
		return errors.New("Ошибка генерации тайла")
	}
	ds.Close()
	return err
}

func TranslateRender(task Task) error {
	var minX, minY, maxX, maxY float64

	switch task.layer.Projection {
	//WebMercarator
	case "EPSG:3857":
		minX, minY, maxX, maxY = translator.WebMercarator(task.x, task.y, task.z)
		intersects := !(maxX < task.layer.UpperLeftX || minX > task.layer.LowerRightX ||
			maxY < task.layer.LowerRightY || minY > task.layer.UpperLeftY)
		if !intersects {
			return errors.New("Ошибка границ слоя")
		}
	case "EPSG:4326":
		minX, minY, maxX, maxY = translator.WGS84(task.x, task.y, task.z)
		intersects := !(maxX < task.layer.UpperLeftX || minX > task.layer.LowerRightX ||
			maxY < task.layer.LowerRightY || minY > task.layer.UpperLeftY)
		if !intersects {
			return errors.New("Ошибка границ слоя")
		}
	default:
	}
	err := os.MkdirAll(task.filePath, 0755)
	if err != nil {
		return err
	}
	options := []string{
		"-projwin",
		fmt.Sprintf("%f", minX),
		fmt.Sprintf("%f", maxY),
		fmt.Sprintf("%f", maxX),
		fmt.Sprintf("%f", minY),
		"-a_srs", task.layer.Projection,
		"-outsize",
		fmt.Sprintf("%d", task.layer.TileSize),
		fmt.Sprintf("%d", task.layer.TileSize)}
	err = gdal.ConvertTile(
		"../resource/map/"+task.layer.SourcePath,
		task.filePath+task.fileName,
		options,
	)
	return err
}
