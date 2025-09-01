package db

import (
	"fmt"
	"geoserver/internal/config"
	"geoserver/internal/db/models"
	"geoserver/internal/db/seeds"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init() error {
	dbPath := getDbPath()

	_, err := os.Stat(dbPath)
	dbExists := !os.IsNotExist(err)

	// 2. Открываем/создаем БД
	db, err := GetConnection()
	if err != nil {
		return err
	}
	// 3. Если БД не существовала - применяем миграции
	if !dbExists {
		err = runMigrate(db)
		if err != nil {
			return err
		}
		err = runSeed(db)
		if err != nil {
			return err
		}
	}
	return nil
}

func getDbPath() string {
	return "../resource/" + config.Configs.DB_DATABASE
}

func GetConnection() (*gorm.DB, error) {

	dbName := getDbPath()
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Ошибка подклчения к базе данных: %w", err)
	}

	return db, nil
}

func runMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		return fmt.Errorf("Ошибка миграции: %w", err)
	}
	err = db.AutoMigrate(&models.Layer{})
	if err != nil {
		return fmt.Errorf("Ошибка миграции: %w", err)
	}
	return nil
}

func runSeed(db *gorm.DB) error {
	err := seeds.User(db)
	if err != nil {
		return fmt.Errorf("Ошибка сида: %w", err)
	}
	err = seeds.Layers(db)
	if err != nil {
		return fmt.Errorf("Ошибка сида: %w", err)
	}
	return nil
}
