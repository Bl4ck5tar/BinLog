package flag

import (
	"BinLog/server/global"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"
	"go.uber.org/zap"
)

var (
	sqlFlag = &cli.BoolFlag{
		Name:	"sql",
		Usage: 	"Initializes the structure of the MySQL database table.",
	}
	sqlExportFlag = &cli.BoolFlag{
		Name: 	"sql-export",
		Usage: 	"Exports SQL data from a specified file.",
	}
	sqlImportFlag = &cli.StringFlag{
		Name: 	"sql-import",
		Usage: 	"Imports SQL data from a specified file.",
	}
	esFlag = &cli.BoolFlag{
		Name: 	"es",
		Usage: 	"initializes the Elasticsearch index.",
	}
	esExportFlag = &cli.StringFlag{
		Name: 	"es-esport",
		Usage: 	"Exports data from Elasticsearch to a specified file.",
	}
	esImportFlag = &cli.StringFlag{
		Name: 	"es-import",
		Usage: 	"Imports data into Elasticsearch from a specified file.",
	}
	adminFlag = &cli.BoolFlag{
		Name: 	"admin",
		Usage: 	"Creates an administrator using the name, email and address specified in the config.yaml file",
	}
)

func Run(c *cli.Context) {
	//检查是否设置了多个标志
	if c.NumFlags() > 1 {
		err := cli.NewExitError("Only one command can be specified", 1)
		global.Log.Error("Invalid command usage:", zap.Error(err))
		os.Exit(1)
	}
	//根据不同的标志选择执行的操作
	switch {
		case c.Bool(sqlFlag.Name):
			if err := SQL(); err != nil {
				global.Log.Error("Failed to create table structure:", zap.Error(err))
				return
			}else {
				global.Log.Info("Successfully created table structure")
			}
		case c.Bool(sqlExportFlag.Name):
			if err := SQLExport(); err != nil {
				global.Log.Error("Failed to export SQL data:", zap.Error(err))
				return
			}else {
				global.Log.Info("Successfully exported SQL data")
			}
		case c.Bool(sqlImportFlag.Name):
			if errs := SQLImport(c.String(sqlImportFlag.Name)); len(errs) > 0 {
				var combinedErrors string
				for _, err := range errs {
					combinedErrors += err.Error() + "\n"
				}
				err := errors.New(combinedErrors)
				global.Log.Error("Failed to import SQL data:", zap.Error(err))
			}else {
				global.Log.Info("Successfully imported SQL data")
			}
		case c.Bool(esFlag.Name):
			if err := ElasticsearchExport(); err != nil {
				global.Log.Error("Failed to export ES data:", zap.Error(err))
			}else {
				global.Log.Info("Successfully exported ES data")
			}
		case c.IsSet(esImportFlag.Name):
			if num, err := ElasticsearchImport(c.String(esImportFlag.Name)); err != nil {
				global.Log.Error("Failed to import ES data:", zap.Error(err))
			}else {
				global.Log.Info(fmt.Sprintf("Successfully imported ES data, totaling %d records", num))
			}
		case c.IsSet(esExportFlag.Name):
			if err := ElasticsearchExport(); err != nil {
				global.Log.Error("Failed to export ES data:", zap.Error(err))
			}else {
				global.Log.Info("Successfully exported ES data")
			}
		case c.Bool(adminFlag.Name):
			if err := Admin(); err != nil {
				global.Log.Error("Failed to create an administrator:", zap.Error(err))
			}else {
				global.Log.Info("Successfully created an administrator")
			}
		default:
			err := cli.NewExitError("unknown command", 1)
			global.Log.Error(err.Error(), zap.Error(err))
	}
}

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Go blog"
	app.Flags = []cli.Flag{
		sqlFlag,
		sqlExportFlag,
		sqlImportFlag,
		esFlag,
		esExportFlag,
		esImportFlag,
		adminFlag,
	}
	app.Action = Run
	return app
}

//initFlag 初始化并运行 CLI 应用程序
func InitFlag() {
	if len(os.Args) > 1 {
		app := NewApp()
		err := app.Run(os.Args)
		if err != nil {
			global.Log.Error("Application execution encountered an error:", zap.Error(err))
			os.Exit(1)
		}
		if os.Args[1] == "-h" || os.Args[1] == "-help" {
			fmt.Println("Displaying help message...")
		}
		os.Exit(0)
	}
}