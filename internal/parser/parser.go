package parser

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

func ExtractProjectionFromWKT(wkt string) (string, error) {
	pattern := `AUTHORITY\["EPSG","(\d+)"\]`
	re := regexp.MustCompile(pattern)

	matches := re.FindAllStringSubmatch(wkt, -1)
	if len(matches) == 0 {
		return "", fmt.Errorf("EPSG code not found in WKT string")
	}
	lastMatch := matches[len(matches)-1]
	if len(lastMatch) < 2 {
		return "", fmt.Errorf("invalid EPSG code format")
	}

	return "EPSG:" + lastMatch[1], nil
}

func ExtractBounds(sourcePath string) (upperLeftX float64, upperLeftY float64, lowerRightX float64, lowerRightY float64, err error) {

	cmd := exec.Command("gdalinfo", sourcePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("gdalinfo failed: %v", err)
	}

	info := string(output)

	// Извлекаем границы исходного файла
	re := regexp.MustCompile(`Upper Left\s*\(\s*([\d\.-]+),\s*([\d\.-]+)`)
	matches := re.FindStringSubmatch(info)
	if len(matches) < 3 {
		return 0, 0, 0, 0, errors.New("could not parse Upper Left coordinates")
	}

	upperLeftX, _ = strconv.ParseFloat(matches[1], 64)
	upperLeftY, _ = strconv.ParseFloat(matches[2], 64)

	re = regexp.MustCompile(`Lower Right\s*\(\s*([\d\.-]+),\s*([\d\.-]+)`)
	matches = re.FindStringSubmatch(info)
	if len(matches) < 3 {
		return 0, 0, 0, 0, errors.New("could not parse Lower Right coordinates")
	}

	lowerRightX, _ = strconv.ParseFloat(matches[1], 64)
	lowerRightY, _ = strconv.ParseFloat(matches[2], 64)

	return upperLeftX, upperLeftY, lowerRightX, lowerRightY, nil
}
