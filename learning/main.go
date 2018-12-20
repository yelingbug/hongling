package main

import (
	"fmt"
	"os"
	"flag"
	"github.com/lestrrat-go/file-rotatelogs"
	"time"
	"github.com/sirupsen/logrus"
	"github.com/rifflock/lfshook"
)

var count = flag.Int("count", 1, "for a test.")

func main() {
	r, _ := rotatelogs.New(
		"/tmp/log1.%Y%m%d%H%M",
		rotatelogs.WithLinkName("/home/yeling/logxx1"),
		rotatelogs.WithRotationTime(1*time.Minute),
		rotatelogs.WithMaxAge(5*time.Minute))
	r1, _ := rotatelogs.New(
		"/tmp/log2.%Y%m%d%H%M",
		rotatelogs.WithLinkName("/home/yeling/logxx2"),
		rotatelogs.WithRotationTime(1*time.Minute),
		rotatelogs.WithMaxAge(5*time.Minute))
	fmt.Println(r)

	gg := logrus.New()
	gg.SetFormatter(&logrus.JSONFormatter{})
	gg.SetLevel(logrus.DebugLevel)
	gg.SetOutput(os.Stdout)

	gg.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.DebugLevel: r,
			logrus.InfoLevel:  r1,
		},
		&logrus.JSONFormatter{},
	))
	fmt.Println("---------------")
	gg.WithFields(logrus.Fields{"abc": 123}).Info("What the fuck")
	gg.WithFields(logrus.Fields{"fre": 123}).Debug("What the fuck")
	/*flag.Parse()
	gg.Debug(*count)*/

	/*app := cli.NewApp()
	app.Name = "WOW"
	app.Author = "Yelin.G"
	app.Usage = "That's for a test."
	app.Version = "99.99.99"

	cli.VersionFlag = cli.BoolFlag{
		Name: "print-version, V",
		Usage: "print only the version",
	}*/

	/*var language string
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name:        "lang",
			Value:       "english",
			Usage:       "language `LLL` for the greeting",
			Destination: &language,
		},
	}

	app.Action = func(c *cli.Context) error {
		name := "someone"
		if c.NArg() > 0 {
			name = c.Args()[0]
		}
		fmt.Println(c.NArg())
		if language == "spanish" {
			fmt.Println("Hola", name)
		} else {
			fmt.Println("Hello", name)
		}
		return nil
	}*/

	/*var language, config string
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "lang, l",
			Value:       "english",
			Usage:       "Language for the greeting",
			Destination: &language,
			EnvVar:      "LANG,GOPATH",
		},
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Load configuration from `FILE`",
			Destination: &config,
		},
		altsrc.NewIntFlag(cli.IntFlag{Name: "test"}),
		cli.StringFlag{Name: "load"},
	}

	app.Commands = []cli.Command{
		{
			Name:    "complete",
			Aliases: []string{"c"},
			Usage:   "complete a task on the list",
			Category:"Action Go",
			Action: func(c *cli.Context) error {
				fmt.Printf("complete [%s] [%d]", language, c.GlobalInt("test"))
				return nil
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a task to the list",
			Category:"Action Go",
			Action: func(c *cli.Context) error {
				fmt.Printf("add [%s] [%s]", language, config)
				return nil
			},
		},
		{
			Name:    "template",
			Aliases: []string{"t"},
			Usage:   "options for task templates",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a new template",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing template",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}

	app.Before = altsrc.InitInputSourceWithContext(app.Flags, altsrc.NewYamlSourceFromFlagFunc("load"))

	sort.Sort(cli.FlagsByName(app.Flags))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}*/
}
