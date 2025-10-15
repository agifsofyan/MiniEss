package configs

import (
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var JwtSecret = []byte("replace-with-secure-secret")

type TokenClaim struct {
	Sub  float64 `json:"sub"` // jwt lib may unmarshal numeric as float64
	Role string  `json:"role"`
	jwt.RegisteredClaims
}

func Connection() (*sqlx.DB, error) {
	dbDriver := GetEnv("DB_DRIVER")
	dbHost := GetEnv("DB_HOST")
	dbPort := GetEnv("DB_PORT")
	dbUser := GetEnv("DB_USER")
	dbName := GetEnv("DB_NAME")
	dbPass := GetEnv("DB_PASS")

	var dsn, driver string

	if dbDriver == "psql" || dbDriver == "postgres" || dbDriver == "mysql" {
		if dbDriver == "psql" || dbDriver == "postgres" {
			dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
			driver = "postgres"
		} else {
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
			driver = "mysql"
		}

		db, err := sqlx.Open(driver, dsn)
		if err != nil {
			return nil, err
		}

		db.SetConnMaxIdleTime(time.Second * 30)
		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(20)
		db.SetConnMaxLifetime(time.Minute * 1)

		// Test the connection to the database
		if err := db.Ping(); err != nil {
			return nil, err
		} else {
			return db, nil
		}
	}

	errMsg := fmt.Sprintf("unsupported database type: %s", dbDriver)
	return nil, errors.New(errMsg)

}
