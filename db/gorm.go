/**
 * @Time: 2022/2/24 11:09
 * @Author: yt.yin
 */

package db

import (
	"os"

	"github.com/goworkeryyt/go-core/global"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

)

// Gorm 初始化数据库并产生数据库全局变量
// @return: *gorm.DB
func Gorm() *gorm.DB {
	switch global.CONFIG.Server.DataDriver {
	case "mysql":
		return GormMySQL()
	case "postgre":
		return GormPostgreSQL()
	case "sqlite":
		return GormSQLite()
	default:
		return GormMySQL()
	}
}

// GormMySQL 初始化Mysql数据库
// @return: *gorm.DB
func GormMySQL() *gorm.DB {
	config := global.CONFIG.MySQL
	if config.Host == "" {
		return nil
	}
	dsn := config.Username + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.Dbname + "?" + config.Config
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig(config.LogLevel)); err != nil {
		global.LOG.Error("MySQL启动异常", zap.Any("err", err))
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
		return db
	}
}

//gormConfig 根据配置决定是否开启日志
//@param: mod bool
//@return: *gorm.Config
func gormConfig(logLevel string) *gorm.Config {
	var config = &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}}
	switch logLevel {
	case "silent", "Silent":
		config.Logger = logger.Default.LogMode(logger.Silent)
	case "error", "Error":
		config.Logger = logger.Default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = logger.Default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = logger.Default.LogMode(logger.Info)
	default:
		config.Logger = logger.Default.LogMode(logger.Error)
	}
	return config
}


// GormPostgreSQL 初始化PostgreSQL数据库
// @return: *gorm.DB
func GormPostgreSQL() *gorm.DB {
	config := global.CONFIG.PostgreSQL
	dsn := "host=" + config.Host + " user=" + config.Username + " password=" + config.Password + " dbname=" + config.Dbname + " port=" + config.Port + " " + config.Config
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // 禁用隐式 prepared statement
	}), gormConfig(config.LogLevel))
	if err != nil {
		global.LOG.Error("PostgreSQL启动异常", zap.Any("err", err))
		os.Exit(0)
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
		return db
	}
}

// GormSQLite 连接sqlite 数据库
func GormSQLite() *gorm.DB {
	config := global.CONFIG.SQLite
	db, err := gorm.Open(sqlite.Open(config.DbPath), gormConfig(config.LogLevel))
	if err != nil {
		global.LOG.Error("SQLite数据库连接失败：", zap.Any("err", err))
		os.Exit(0)
		return nil
	}else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
		return db
	}
}

