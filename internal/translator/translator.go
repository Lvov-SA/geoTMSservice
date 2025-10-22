package translator

import (
	"errors"
	"math"
)

func WebMercarator(x, y, z int) (minX, minY, maxX, maxY float64) {
	worldSize := 20037508.342789244 * 2
	tileSize := worldSize / math.Pow(2, float64(z))

	minX = -20037508.342789244 + float64(x)*tileSize
	maxX = minX + tileSize
	maxY = 20037508.342789244 - float64(y)*tileSize
	minY = maxY - tileSize

	return minX, minY, maxX, maxY
}

func WGS84(x, y, z int) (minLon, minLat, maxLon, maxLat float64) {
	minLon = float64(x)/math.Pow(2.0, float64(z))*360.0 - 180.0
	maxLon = float64(x+1)/math.Pow(2.0, float64(z))*360.0 - 180.0

	minLat = mercatorToLat(math.Pi * (1 - 2*float64(y+1)/math.Pow(2.0, float64(z))))
	maxLat = mercatorToLat(math.Pi * (1 - 2*float64(y)/math.Pow(2.0, float64(z))))

	return minLon, minLat, maxLon, maxLat
}

func mercatorToLat(y float64) float64 {
	return 180.0 / math.Pi * math.Atan(math.Sinh(y))
}

func ExtToOptionForTranslate(ext string) (string, error) {
	switch ext {
	case "png":
		return "PNG", nil
	case "jpg":
		return "JPG", nil
	case "webp":
		return "WEBP", nil
	default:
		return "", errors.New("Не поддерживаемый формат тайла")
	}
}
