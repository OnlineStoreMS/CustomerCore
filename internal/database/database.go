package database

import (
	"fmt"
	"os"
	"path/filepath"

	"customercore/internal/config"
	"customercore/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch cfg.Driver {
	case "postgres":
		dialector = postgres.Open(cfg.PostgresDSN)
	case "sqlite":
		if err := os.MkdirAll(filepath.Dir(cfg.SQLitePath), 0o755); err != nil {
			return nil, err
		}
		dialector = sqlite.Open(cfg.SQLitePath)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
	return gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.Customer{},
		&model.CustomerAddress{},
		&model.CustomerChannelBinding{},
	); err != nil {
		return err
	}
	if db.Dialector.Name() == "postgres" {
		return db.Exec(`
			CREATE INDEX IF NOT EXISTS idx_customers_tenant_name ON customers (tenant_id, display_name);
			CREATE INDEX IF NOT EXISTS idx_customers_tenant_created ON customers (tenant_id, created_at DESC);
			CREATE INDEX IF NOT EXISTS idx_customer_addresses_customer ON customer_addresses (tenant_id, customer_id);
			CREATE INDEX IF NOT EXISTS idx_customer_bindings_customer ON customer_channel_bindings (tenant_id, customer_id);
		`).Error
	}
	return nil
}
