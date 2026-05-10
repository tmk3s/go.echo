package repository_test

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mysqlDSN(dbName string) string {
	user := getenv("MYSQL_ID", "root")
	pass := getenv("MYSQL_PASSWORD", "pass")
	host := getenv("MYSQL_HOST", "db")
	port := getenv("MYSQL_PORT", "3306")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, dbName)
}

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	baseName := getenv("MYSQL_DATABASE", "develop")
	testDBName := baseName + "_test"

	admin, err := gorm.Open(mysql.Open(mysqlDSN("")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect to MySQL: %v", err)
	}
	if err := admin.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", testDBName)).Error; err != nil {
		t.Fatalf("failed to create test database %s: %v", testDBName, err)
	}

	db, err := gorm.Open(mysql.Open(mysqlDSN(testDBName)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test database %s: %v", testDBName, err)
	}
	return db
}

func truncate(t *testing.T, db *gorm.DB, tables ...string) {
	t.Helper()
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", table)).Error; err != nil {
			t.Errorf("failed to truncate %s: %v", table, err)
		}
	}
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}
