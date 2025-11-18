package repository

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLmode  string
}

func NewPostgresDB(cfg DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLmode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func DBConfigFromViper() DBConfig {
	return DBConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLmode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	}
}
