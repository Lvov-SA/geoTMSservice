package models

import "gorm.io/gorm"

type Layer struct {
	gorm.Model

	// Идентификаторы
	Name  string `gorm:"size:100;uniqueIndex;not null"` // Уникальный ключ (например: "osm_roads")
	Title string `gorm:"size:255;not null"`             // Отображаемое название ("Дороги OSM")

	// Источник данных
	SourcePath string `gorm:"not null"` // Путь к файлу в папке resource/map (имя с расшрирением)

	// Настройки видимости
	MinZoom  int  `gorm:"default:0"`    // Минимальный zoom
	MaxZoom  int  `gorm:"default:10"`   // Максимальный zoom
	IsActive bool `gorm:"default:true"` // Доступен для запросов

	Width      int //Ширина исходника
	Height     int //Высота исходника
	TileSize   int `gorm:"default:256"` //размер тайла для слоя
	Projection string

	//Границы слоя
	UpperLeftX  float64
	UpperLeftY  float64
	LowerRightX float64
	LowerRightY float64
}
