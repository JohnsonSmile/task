package database

import (
	"fmt"
	"log"
	"os"
	"time"
	"twitter_task/infra/config"
	"twitter_task/model"

	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db          *gorm.DB
	initialized bool
)

func GetDatabase() *gorm.DB {
	if !initialized {
		panic("should initialize first!")
	}
	return db
}

func Initialize(isDebug bool) {
	var (
		err error
		dsn string
	)

	conf := config.GetServerConfig()

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.MySql.Username, conf.MySql.Password, conf.MySql.Address, conf.MySql.Port, conf.MySql.Database)
	var dbLogger logger.Interface
	if isDebug {
		dbLogger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 彩色打印
		})
	} else {
		dbLogger = logger.New(log.New(&lumberjack.Logger{
			// Filename:   "/var/log/mxshop/userweb/info.log",
			Filename:   "./log/database.log",
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
		}, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      false,       // 彩色打印
		})
	}
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_", // 表名前缀，`User`表为`t_users`
			SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
		},
		Logger:                                   dbLogger,
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
	})
	if err != nil {
		panic(err)
	}
	// 自动迁移
	db.AutoMigrate(&model.User{}, &model.Twitter{})
	initialized = true
}
