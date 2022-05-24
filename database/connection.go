package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
)

func init() {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "swift_" + defaultTableName
	}
}

// OpenDatabase Open Connection to database with specific settings
func OpenDatabase(driver, url string) *gorm.DB {
	db, err := gorm.Open(driver, url)
	if err != nil {
		panic("failed to connect to database")
	}

	db.SetLogger(log.New(os.Stdout, "\r\n", 0))
	return db
}

var dbFactory *DbFactory
var db *gorm.DB

// Open create new connection to DB with default config
func Open() *gorm.DB {
	if db == nil {
		db = OpenDatabase(dbFactory.Driver, dbFactory.URL)
	}

	return db
}

// SetupFactory setup DatabaseFactory
func SetupFactory(driver, url string) *DbFactory {
	dbFactory = &DbFactory{
		Driver: driver,
		URL:    url,
	}

	return dbFactory
}

// GetFactory get DatabaseFactory instance
func GetFactory() *DbFactory {
	return dbFactory
}

// Close close existing connection
func Close() {
	if db != nil {
		db.Close()
	}
}
