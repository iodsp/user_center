package context

import (
	"github.com/iodsp/user_center/fsutil"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kylelemons/go-gypsy/yaml"
	"github.com/urfave/cli"
	"log"
	"time"
)

const (
	DbMySQL = "mysql"
)

// Config provides a struct in which application configuration is stored.
// Application code must use functions to get config values, for two reasons:
//
// 1. Some values are computed and we don't want to leak implementation details (aims at reducing refactoring overhead).
//
// 2. Paths might actually be dynamic later (if we build a multi-user version).
//
// See https://github.com/photoprism/photoprism/issues/50#issuecomment-433856358
type Config struct {
	appName            string
	appVersion         string
	appCopyright       string
	debug              bool
	configFile         string
	httpServerHost     string
	httpServerPort     int
	httpServerMode     string
	databaseDriver     string
	databaseDsn        string
	db                 *gorm.DB
}

// NewConfig() creates a new configuration entity by using two methods:
//
// 1. SetValuesFromFile: This will initialize values from a yaml config file.
//
// 2. SetValuesFromCliContext: Which comes after SetValuesFromFile and overrides
//    any previous values giving an option two override file configs through the CLI.
func NewConfig(ctx *cli.Context) *Config {
	c := &Config{}
	c.appName = ctx.App.Name
	c.appCopyright = ctx.App.Copyright
	c.appVersion = ctx.App.Version
	c.SetValuesFromFile(fsutil.ExpandedFilename(ctx.GlobalString("config-file")))
	c.SetValuesFromCliContext(ctx)

	return c
}

// SetValuesFromFile uses a yaml config file to initiate the configuration entity.
func (c *Config) SetValuesFromFile(fileName string) error {
	yamlConfig, err := yaml.ReadFile(fileName)

	if err != nil {
		return err
	}

	c.configFile = fileName
	if debug, err := yamlConfig.GetBool("debug"); err == nil {
		c.debug = debug
	}

	if httpServerHost, err := yamlConfig.Get("http-host"); err == nil {
		c.httpServerHost = httpServerHost
	}

	if httpServerPort, err := yamlConfig.GetInt("http-port"); err == nil {
		c.httpServerPort = int(httpServerPort)
	}

	if httpServerMode, err := yamlConfig.Get("http-mode"); err == nil {
		c.httpServerMode = httpServerMode
	}

	if databaseDriver, err := yamlConfig.Get("database-driver"); err == nil {
		c.databaseDriver = databaseDriver
	}

	if databaseDsn, err := yamlConfig.Get("database-dsn"); err == nil {
		c.databaseDsn = databaseDsn
	}

	return nil
}

// SetValuesFromCliContext uses values from the CLI to setup configuration overrides
// for the entity.
func (c *Config) SetValuesFromCliContext(ctx *cli.Context) error {
	if ctx.GlobalBool("debug") {
		c.debug = ctx.GlobalBool("debug")
	}

	if ctx.GlobalIsSet("database-driver") || c.databaseDriver == "" {
		c.databaseDriver = ctx.GlobalString("database-driver")
	}

	if ctx.GlobalIsSet("database-dsn") || c.databaseDsn == "" {
		c.databaseDsn = ctx.GlobalString("database-dsn")
	}

	if ctx.GlobalIsSet("http-host") || c.httpServerHost == "" {
		c.httpServerHost = ctx.GlobalString("http-host")
	}

	if ctx.GlobalIsSet("http-port") || c.httpServerPort == 0 {
		c.httpServerPort = ctx.GlobalInt("http-port")
	}

	if ctx.GlobalIsSet("http-mode") || c.httpServerMode == "" {
		c.httpServerMode = ctx.GlobalString("http-mode")
	}

	return nil
}

// connectToDatabase establishes a database connection.
// When used with the internal driver, it may create a new database server instance.
// It tries to do this 12 times with a 5 second sleep interval in between.
func (c *Config) connectToDatabase() error {
	dbDriver := c.DatabaseDriver()
	dbDsn := c.DatabaseDsn()

	db, err := gorm.Open(dbDriver, dbDsn)

	if err != nil || db == nil {
		for i := 1; i <= 12; i++ {
			time.Sleep(5 * time.Second)

			db, err = gorm.Open(dbDriver, dbDsn)

			if db != nil && err == nil {
				break
			}
		}

		if err != nil || db == nil {
			log.Fatal(err)
		}
	}

	c.db = db

	return err
}

// AppName returns the application name.
func (c *Config) AppName() string {
	return c.appName
}

// AppVersion returns the application version.
func (c *Config) AppVersion() string {
	return c.appVersion
}

// AppCopyright returns the application copyright.
func (c *Config) AppCopyright() string {
	return c.appCopyright
}

// Debug returns true if debug mode is on.
func (c *Config) Debug() bool {
	return c.debug
}

// ConfigFile returns the config file name.
func (c *Config) ConfigFile() string {
	return c.configFile
}

// HttpServerHost returns the built-in HTTP server host name or IP address (empty for all interfaces).
func (c *Config) HttpServerHost() string {
	return c.httpServerHost
}

// HttpServerPort returns the built-in HTTP server port.
func (c *Config) HttpServerPort() int {
	return c.httpServerPort
}

// HttpServerMode returns the server mode.
func (c *Config) HttpServerMode() string {
	return c.httpServerMode
}

// DatabaseDriver returns the database driver name.
func (c *Config) DatabaseDriver() string {
	return c.databaseDriver
}

// DatabaseDsn returns the database data source name (DSN).
func (c *Config) DatabaseDsn() string {
	return c.databaseDsn
}

// Db returns the db connection.
func (c *Config) Db() *gorm.DB {
	if c.db == nil {
		c.connectToDatabase()
	}

	return c.db
}
