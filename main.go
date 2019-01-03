package main

import (
	"fmt"
	"hongling/utility"
	"os"
	"gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{
		Name:  "hl",
		Usage: "简易命令行工具.",
		Description: `一个内部使用的简易命令行工具.
   1. 支持在测试和生产环境的数据库上执行SQL.`,
		HelpName:  "hl",
		UsageText: "hl [global options] command [command options] [arguments...]",
		Version:   "1.0.beta",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "environment",
				Aliases: []string{"env", "e"},
				Usage:   "选择环境" + utility.DEV + "/" + utility.TEST + "/" + utility.PROD + ", 缺省为" + utility.DEV + ".",
				Value:   utility.DEV,
			},
		},
		Action: utility.Menu,
	}

	app.Commands = []*cli.Command{utility.DrdsCommand, utility.ArchetypeCommand, utility.ExecjavaCommand}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
