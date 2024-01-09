package database

import (
	"fmt"
	"log"

	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	"github.com/dathuynh1108/clean-arch-base/pkg/database/dbpool"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase() {
	// Init db connection here
	// Format database.master or database.replica
	config := config.GetConfig()
	for alias, databases := range config.DatabasesConfig {
		for _, database := range databases {
			log.Println("--> Init database:", alias)
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
			default:
				panic(fmt.Sprintf("Unsupported database engine: %v", database.Engine))
			}

			connection, err := gorm.Open(dialector, &gorm.Config{})
			if err != nil {
				panic(err)
			}

			dbPool.SetDB(alias, connection)
		}
	}

}

var (
	dbPool = dbpool.NewDBPool[*gorm.DB]()
)

func ProvideDBPool() dbpool.DBPool[*gorm.DB] {
	return dbPool
}
