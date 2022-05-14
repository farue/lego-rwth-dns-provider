package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	//zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02T15:04:05-0700",
	})

	app := &cli.App{
		Name:  "lego-rwth-dns-provider",
		Usage: "lego provider for RWTH DNS",
		Before: func(c *cli.Context) error {
			quiet := c.Bool("quiet")
			if quiet {
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:      "present",
				Aliases:   []string{"p"},
				Usage:     "present a DNS record",
				ArgsUsage: "[FQDN] [recordTxt]",
				Before: func(c *cli.Context) error {
					fqdn := c.Args().Get(0)
					if fqdn == "" {
						return fmt.Errorf("FQDN is required")
					}
					recordTxt := c.Args().Get(1)
					if recordTxt == "" {
						return fmt.Errorf("recordTxt is required")
					}
					return nil
				},
				Action: func(c *cli.Context) error {
					fqdn := c.Args().Get(0)
					recordTxt := c.Args().Get(1)
					token := c.String("token")
					api := NewApiClient(&http.Client{})
					return api.present(fqdn, recordTxt, token)
				},
			},
			{
				Name:      "cleanup",
				Aliases:   []string{"c"},
				Usage:     "cleanup a DNS record",
				ArgsUsage: "[FQDN] [recordTxt]",
				Before: func(c *cli.Context) error {
					fqdn := c.Args().Get(0)
					if fqdn == "" {
						return fmt.Errorf("FQDN is required")
					}
					recordTxt := c.Args().Get(1)
					if recordTxt == "" {
						return fmt.Errorf("recordTxt is required")
					}
					return nil
				},
				Action: func(c *cli.Context) error {
					fqdn := c.Args().Get(0)
					recordTxt := c.Args().Get(1)
					token := c.String("token")
					api := NewApiClient(&http.Client{})
					return api.cleanup(fqdn, recordTxt, token)
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "token",
				Aliases:  []string{"t"},
				Usage:    "RWTH DNS API token",
				EnvVars:  []string{"RWTH_DNS_API_TOKEN"},
				Required: true,
			},
			&cli.BoolFlag{
				Name:    "quiet",
				Aliases: []string{"q"},
				Usage:   "disable debug messages",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Msg("A fatal error occurred")
	}
}
