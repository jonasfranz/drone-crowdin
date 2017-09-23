package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

var build = "0" //build number is set on compile time
var version = "1.0.0"

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
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.BoolFlag{
			Name:   "ignore-branch",
			Usage:  "if true it will not pass the branch to crowdin",
			EnvVar: "PLUGIN_IGNORE_BRANCH",
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
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
	if !c.Bool("ignore-branch") {
		plugin.Branch = c.String("commit.branch")
	}
	return plugin.Exec()
}
