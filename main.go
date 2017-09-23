package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

var build = "0" //build number is set on compile time
var version = "0.1.0"

func main() {
	app := cli.NewApp()
	app.Name = "crowdin plugin"
	app.Usage = "upload translations"
	app.Version = fmt.Sprintf("%s+%s", version, build) // 0.1.0+3
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "project-key",
			Usage:  "project key",
			EnvVar: "PLUGIN_PROJECT_KEY,CROWDIN_KEY",
		},
		cli.StringFlag{
			Name:   "project-identifier",
			Usage:  "project identifier",
			EnvVar: "PLUGIN_PROJECT_IDENTIFIER,CROWDIN_IDENTIFIER",
		},
		cli.GenericFlag{
			Name:   "files",
			Usage:  "files for updating translations",
			EnvVar: "PLUGIN_FILES",
			Value:  &StringMapFlag{},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Files: c.Generic("files").(*StringMapFlag).Get(),
		Config: Config{
			Identifier: c.String("project-identifier"),
			Key:        c.String("project-key"),
		},
	}
	return plugin.Exec()
}
