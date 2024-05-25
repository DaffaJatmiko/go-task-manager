package migrations

import (
    "github.com/DaffaJatmiko/go-task-manager/model"
    "gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
    // Migrasi model User
    if err := db.AutoMigrate(&model.User{}); err != nil {
        return err
    }

    // Migrasi model Category
    if err := db.AutoMigrate(&model.Category{}); err != nil {
        return err
    }

    // Migrasi model Task
    if err := db.AutoMigrate(&model.Task{}); err != nil {
        return err
    }

    // Migrasi model Session
    if err := db.AutoMigrate(&model.Session{}); err != nil {
        return err
    }

    return nil
}
