package fargatedeploy

import "github.com/urfave/cli"

// FlagSet ...flagを設定
func FlagSet() *cli.App {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration *.toml",
		},
		cli.StringFlag{
			Name:  "env, e",
			Usage: "Choose Env",
		},
		cli.StringFlag{
			Name:  "id, i",
			Value: "",
			Usage: "Docker build /push from ECR ID",
		},
		cli.BoolFlag{
			Name:   "push, p",
			Hidden: false,
			Usage:  "Push an image or a repository to a registry (default:false)",
		},
		cli.BoolFlag{
			Name:   "lifecycle, l",
			Hidden: false,
			Usage:  "Delete images of unused tags (default:false)",
		},
	}
	return app
}
