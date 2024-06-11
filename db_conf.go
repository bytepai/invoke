package invoke

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var (
	dbConfigPath = "db_conf.json"
	dbConfigLock sync.Mutex
	DBConfig     *DatabaseConfig
)

type PostgreSQLConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

type OracleConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

type SQLiteConfig struct {
	DBPath string `json:"dbpath"`
}

type MySQLConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

type MongoDBConfig struct {
	URI      string `json:"uri"`
	Database string `json:"database"`
}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type DatabaseConfig struct {
	PostgreSQL PostgreSQLConfig `json:"postgresql"`
	Oracle     OracleConfig     `json:"oracle"`
	SQLite     SQLiteConfig     `json:"sqlite"`
	MySQL      MySQLConfig      `json:"mysql"`
	MongoDB    MongoDBConfig    `json:"mongodb"`
	Redis      RedisConfig      `json:"redis"`
}

// LoadDBConfig loads configuration from a JSON file
func LoadDBConfig(filePath string) (*DatabaseConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config DatabaseConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// createDefaultDBConfig creates a default configuration file
func createDefaultDBConfig(filePath string) error {
	defaultConfig := DatabaseConfig{
		PostgreSQL: PostgreSQLConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			DBName:   "postgres",
			SSLMode:  "disable",
		},
		Oracle: OracleConfig{
			User:     "oracle_user",
			Password: "password",
			DBName:   "oracledb",
		},
		SQLite: SQLiteConfig{
			DBPath: "sqlite.db",
		},
		MySQL: MySQLConfig{
			Host:     "localhost",
			Port:     3306,
			User:     "mysql_user",
			Password: "password",
			DBName:   "mysqldb",
		},
		MongoDB: MongoDBConfig{
			URI:      "mongodb://localhost:27017",
			Database: "mydatabase",
		},
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(defaultConfig); err != nil {
		return err
	}
	return nil
}

func InitializeDBConfig() (*DatabaseConfig, error) {
	if _, err := os.Stat(dbConfigPath); os.IsNotExist(err) {
		log.Println("DB config file not found, creating default config")
		if err := createDefaultDBConfig(dbConfigPath); err != nil {
			return nil, err
		}
	}
	config, err := LoadDBConfig(dbConfigPath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func init() {
	var err error
	DBConfig, err = InitializeDBConfig()
	if err != nil {
		log.Fatalf("Failed to initialize DB config: %v", err)
	}
}
