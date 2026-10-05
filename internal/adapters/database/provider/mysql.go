package database_provider

import (
	"account-service/config"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQLClient(cfg config.Config) (*gorm.DB, error) {
	dbCfg := cfg.Database
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=UTC",
		dbCfg.Username,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true, // cache prepared sql statements
	})
	if err != nil {
		return nil, err
	}

	// Configure underlying *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(dbCfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dbCfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(dbCfg.ConnMaxLifetime)

	return db, nil
}
