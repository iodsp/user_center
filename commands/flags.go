package commands

import "github.com/urfave/cli"

// Global CLI flags
var GlobalFlags = []cli.Flag{
	cli.BoolFlag{
		Name:   "debug",
		Usage:  "run in debug mode",
		EnvVar: "USER_CENTER_DEBUG",
	},
	cli.StringFlag{
		Name:   "config-file, c",
		Usage:  "load configuration from `FILENAME`",
		Value:  "/etc/user_center/uc.yml",
		EnvVar: "USER_CENTER_CONFIG_FILE",
	},
	cli.StringFlag{
		Name:   "database-driver",
		Usage:  "database `DRIVER` (internal or mysql)",
		Value:  "mysql",
		EnvVar: "USER_CENTER_DATABASE_DRIVER",
	},
	cli.StringFlag{
		Name:   "database-dsn",
		Usage:  "database data source name (`DSN`)",
		Value:  "root:root@tcp(192.168.33.11:3306)/iodsp?parseTime=true&charset=utf8&loc=Asia%2FShanghai",
		EnvVar: "USER_CENTER_DATABASE_DSN",
	},
	cli.IntFlag{
		Name:   "http-port, p",
		Usage:  "HTTP server port",
		Value:  8081,
		EnvVar: "USER_CENTER_HTTP_PORT",
	},
	cli.StringFlag{
		Name:   "http-host, i",
		Usage:  "HTTP server host",
		Value:  "",
		EnvVar: "USER_CENTER_HTTP_HOST",
	},
	cli.StringFlag{
		Name:   "http-mode, m",
		Usage:  "debug, release or test",
		Value:  "",
		EnvVar: "USER_CENTER_HTTP_MODE",
	},
}
