package database

import (
	"fmt"
	"time"

	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	"github.com/dathuynh1108/clean-arch-base/pkg/database/dbpool"
	"github.com/dathuynh1108/clean-arch-base/pkg/logger"

	"github.com/sirupsen/logrus"

	mysql "go.elastic.co/apm/module/apmgormv2/v2/driver/mysql"
	postgres "go.elastic.co/apm/module/apmgormv2/v2/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitDatabase() error {
	// Init db connection here
	// Format database.master or database.replica
	config := config.GetConfig()
	for alias, databases := range config.DatabasesConfig {
		for _, database := range databases {
			logger.GetLogger().
				WithFields(logrus.Fields{
					"alias":    alias,
					"engine":   database.Engine,
					"instance": database.Instance,
					"host":     database.Host,
					"port":     database.Port,
					"database": database.Database,
					"schema":   database.Schema,
				}).
				Infof("Init database")

			if database.Instance != dbpool.AliasMaster && database.Instance != dbpool.AliasReplica {
				return fmt.Errorf(`unsupported database instance, want "master" or "replica", got : "%v"`, database.Instance)
			}

			var dialector gorm.Dialector
			switch database.Engine {
			case EngineMySQL:
				dsn := fmt.Sprintf(
					"%s:%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
					database.Username,
					database.Password,
					"tcp",
					database.Host,
					database.Port,
					database.Database,
				)
				dialector = mysql.Open(dsn)

			case EngineSQLite:
			case EnginePostgres:
				dsn := fmt.Sprintf(
					"host=%v user=%v password=%v dbname=%v port=%v sslmode=%v",
					database.Host,
					database.Username,
					database.Password,
					database.Database,
					database.Port,
					database.SSLMode,
				)
				dialector = postgres.Open(dsn)
			default:
				return fmt.Errorf("unsupported database engine: %v", database.Engine)
			}

			connection, err := gorm.Open(dialector, &gorm.Config{
				Logger: gormlogger.Default.LogMode(convertLogMode(database.LogMode)),
				NamingStrategy: schema.NamingStrategy{
					TablePrefix:   database.Schema + ".",
					SingularTable: false,
				},
			})
			if err != nil {
				return err
			}

			sqlDB, _ := connection.DB()
			sqlDB.SetMaxIdleConns(2)
			sqlDB.SetMaxOpenConns(400)
			sqlDB.SetConnMaxLifetime(5 * time.Minute)
			sqlDB.SetConnMaxIdleTime(2 * time.Minute)

			dbPool.SetDB(dbpool.BuildAlias(dbpool.DBAlias(alias), database.Instance), dbpool.NewDB(connection))
		}
	}
	return nil
}

var (
	dbPool = dbpool.NewDBPool()
)

func GetDBPool() dbpool.DBPool {
	return dbPool
}

func convertLogMode(logMode string) gormlogger.LogLevel {
	switch logMode {
	case "silent":
		return gormlogger.Silent
	case "error":
		return gormlogger.Error
	case "warn":
		return gormlogger.Warn
	case "info":
		return gormlogger.Info
	default:
		return gormlogger.Silent
	}
}
