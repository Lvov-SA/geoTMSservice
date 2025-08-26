package translator

import "math"

func WebMercarator(x, y, z int) (minX, minY, maxX, maxY float64) {
	worldSize := 20037508.342789244 * 2
	tileSize := worldSize / math.Pow(2, float64(z))

	minX = -20037508.342789244 + float64(x)*tileSize
	maxX = minX + tileSize
	maxY = 20037508.342789244 - float64(y)*tileSize
	minY = maxY - tileSize

	return minX, minY, maxX, maxY
}
