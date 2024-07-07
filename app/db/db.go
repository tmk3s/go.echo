package db

import (
  "gorm.io/gorm"
  "gorm.io/gorm/logger"
  // "gorm.io/driver/sqlite"
  "gorm.io/driver/mysql"
)

var db *gorm.DB

func init() {
    var err error
    // newLogger := logger.New(
    //   log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
    //   logger.Config{
    //     SlowThreshold:              time.Second,   // Slow SQL threshold
    //     LogLevel:                   logger.Silent, // Log level
    //     IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
    //     ParameterizedQueries:      true,           // Don't include params in the SQL log
    //     Colorful:                  false,          // Disable color
    //   },
    // )
    dsn := "root:pass@tcp(db:3306)/develop?charset=utf8mb4&parseTime=True&loc=Local"
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{ Logger: logger.Default.LogMode(logger.Info) })
    if err != nil {
      panic("failed to connect database")
    }
    db.AutoMigrate(&User{}, &Todo{}, &UserInfo{}, &UserAddress{}, &Tag{}, &PostCode{})
}