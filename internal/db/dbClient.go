package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/plyovchev/sumup-assignment-notifications/internal/config"
	internal_logger "github.com/plyovchev/sumup-assignment-notifications/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Repository defines a interface for access the database.
type DbClient interface {
	Model(value interface{}) *gorm.DB
	Select(query interface{}, args ...interface{}) *gorm.DB
	Find(out interface{}, where ...interface{}) *gorm.DB
	Exec(sql string, values ...interface{}) *gorm.DB
	First(out interface{}, where ...interface{}) *gorm.DB
	Raw(sql string, values ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Updates(value interface{}) *gorm.DB
	Delete(value interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Preload(column string, conditions ...interface{}) *gorm.DB
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB
	ScanRows(rows *sql.Rows, result interface{}) error
	Transaction(fc func(tx DbClient) error) (err error)
	Close() error
	DropTableIfExists(value interface{}) error
	AutoMigrate(value interface{}) error
}

// repository defines a repository for access the database.
type dbClient struct {
	db *gorm.DB
}

func NewDbClient(
	schema string,
	logger *internal_logger.AppLogger,
	cfg *config.Config,
) DbClient {
	logger.Info().Msg("Try establish connection")
	db, err := connectDatabase(schema, cfg)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failure database connection")
		return nil
	}

	logger.Info().Msg("Successfully established database connection.")
	return &dbClient{db: db}
}

const (
	// POSTGRES represents PostgreSQL
	POSTGRES = "postgres"
)

func connectDatabase(schemaName string, config *config.Config) (*gorm.DB, error) {
	var dsn string
	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,         // Don't include params in the SQL log
			Colorful:                  false,        // Disable color
		},
	)
	gormConfig := &gorm.Config{
		Logger: logger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   schemaName,
			SingularTable: false,
		},
	}

	if config.Database.Dialect == POSTGRES {
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			config.Database.Host, config.Database.Port, config.Database.Username,
			config.Database.Dbname, config.Database.Password)

		return gorm.Open(postgres.Open(dsn), gormConfig)
	}
	return nil, errors.New("connection to the DB cannot be established")
}

// Model specify the model you would like to run db operations
func (rep *dbClient) Model(value interface{}) *gorm.DB {
	return rep.db.Model(value)
}

// Select specify fields that you want to retrieve from database when querying, by default, will select all fields;
func (rep *dbClient) Select(query interface{}, args ...interface{}) *gorm.DB {
	return rep.db.Select(query, args...)
}

// Find find records that match given conditions.
func (rep *dbClient) Find(out interface{}, where ...interface{}) *gorm.DB {
	return rep.db.Find(out, where...)
}

// Exec exec given SQL using by gorm.DB.
func (rep *dbClient) Exec(sql string, values ...interface{}) *gorm.DB {
	return rep.db.Exec(sql, values...)
}

// First returns first record that match given conditions, order by primary key.
func (rep *dbClient) First(out interface{}, where ...interface{}) *gorm.DB {
	return rep.db.First(out, where...)
}

// Raw returns the record that executed the given SQL using gorm.DB.
func (rep *dbClient) Raw(sql string, values ...interface{}) *gorm.DB {
	return rep.db.Raw(sql, values...)
}

// Create insert the value into database.
func (rep *dbClient) Create(value interface{}) *gorm.DB {
	return rep.db.Create(value)
}

// Save update value in database, if the value doesn't have primary key, will insert it.
func (rep *dbClient) Save(value interface{}) *gorm.DB {
	return rep.db.Save(value)
}

// Update update value in database
func (rep *dbClient) Updates(value interface{}) *gorm.DB {
	return rep.db.Updates(value)
}

// Delete delete value match given conditions.
func (rep *dbClient) Delete(value interface{}) *gorm.DB {
	return rep.db.Delete(value)
}

// Where returns a new relation.
func (rep *dbClient) Where(query interface{}, args ...interface{}) *gorm.DB {
	return rep.db.Where(query, args...)
}

// Preload preload associations with given conditions.
func (rep *dbClient) Preload(column string, conditions ...interface{}) *gorm.DB {
	return rep.db.Preload(column, conditions...)
}

// Scopes pass current database connection to arguments `func(*DB) *DB`,
// which could be used to add conditions dynamically
func (rep *dbClient) Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB {
	return rep.db.Scopes(funcs...)
}

// ScanRows scan `*sql.Rows` to give struct
func (rep *dbClient) ScanRows(rows *sql.Rows, result interface{}) error {
	return rep.db.ScanRows(rows, result)
}

// Close close current db connection. If database connection is not an io.Closer, returns an error.
func (rep *dbClient) Close() error {
	sqlDB, _ := rep.db.DB()
	return sqlDB.Close()
}

// DropTableIfExists drop table if it is exist
func (rep *dbClient) DropTableIfExists(value interface{}) error {
	return rep.db.Migrator().DropTable(value)
}

// AutoMigrate run auto migration for given models, will only add missing fields, won't delete/change current data
func (rep *dbClient) AutoMigrate(value interface{}) error {
	return rep.db.AutoMigrate(value)
}

// Transaction start a transaction as a block.
// If it is failed, will rollback and return error.
// If it is sccuessed, will commit.
// ref: https://github.com/jinzhu/gorm/blob/master/main.go#L533
func (rep *dbClient) Transaction(fc func(tx DbClient) error) (err error) {
	panicked := true
	tx := rep.db.Begin()
	defer func() {
		if panicked || err != nil {
			tx.Rollback()
		}
	}()

	txrep := &dbClient{}
	txrep.db = tx
	err = fc(txrep)

	if err == nil {
		err = tx.Commit().Error
	}

	panicked = false
	return
}
