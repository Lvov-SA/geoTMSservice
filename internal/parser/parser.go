package parser

import (
	"fmt"
	"regexp"
)

func ExtractProjectionFromWKT(wkt string) (string, error) {
	// Паттерн для поиска AUTHORITY["EPSG","XXXXX"]
	pattern := `AUTHORITY\["EPSG","(\d+)"\]`
	re := regexp.MustCompile(pattern)

	matches := re.FindAllStringSubmatch(wkt, -1)
	if len(matches) == 0 {
		return "", fmt.Errorf("EPSG code not found in WKT string")
	}

	// Берем последний AUTHORITY - это обычно код всей системы координат
	lastMatch := matches[len(matches)-1]
	if len(lastMatch) < 2 {
		return "", fmt.Errorf("invalid EPSG code format")
	}

	return "EPSG:" + lastMatch[1], nil
}
