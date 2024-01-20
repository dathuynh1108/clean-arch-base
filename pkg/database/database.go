package database

import (
	"fmt"

	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	"github.com/dathuynh1108/clean-arch-base/pkg/database/dbpool"
	"github.com/dathuynh1108/clean-arch-base/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase() error {
	// Init db connection here
	// Format database.master or database.replica
	config := config.GetConfig()
	for alias, databases := range config.DatabasesConfig {
		for _, database := range databases {
			logger.GetLogger().Infof("Init database alias: %v, engine: %v, instance: %v", alias, database.Engine, database.Instance)
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
				// dsn := fmt.Sprintf(
				// 	"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
				// 	database.Host,
				// 	database.Username,
				// 	database.Password,
				// 	database.Database,
				// 	database.Port,
				// )
				// dialector = postgres.Open(dsn)
			default:
				return fmt.Errorf("unsupported database engine: %v", database.Engine)
			}

			connection, err := gorm.Open(dialector, &gorm.Config{})
			if err != nil {
				return err
			}

			if database.Instance != dbpool.AliasMaster && database.Instance != dbpool.AliasReplica {
				return fmt.Errorf(`unsupported database instance, want "master" or "replica", got : "%v"`, database.Instance)
			}

			dbPool.SetDB(dbpool.BuildAlias(dbpool.DBAlias(alias), database.Instance), dbpool.NewDB(connection))
		}
	}
	return nil
}

var (
	dbPool = dbpool.NewDBPool()
)

func ProvideDBPool() dbpool.DBPool {
	return dbPool
}
