package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/urfave/cli/v3"
)

// Execute is the entry point called by main.go
func Execute() {
	cmd := &cli.Command{
		Name:  "simple-curl",
		Usage: "A simple curl clone in Go",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "method",
				Aliases: []string{"X"},
				Value:   "GET",
				Usage:   "HTTP method to use (GET, POST, etc.)",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			url := c.Args().First()
			if url == "" {
				return fmt.Errorf("error: you must provide a URL")
			}

			method := c.String("method")

			req, err := http.NewRequestWithContext(ctx, method, url, nil)
			if err != nil {
				return fmt.Errorf("failed to create request: %w", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("request failed: %w", err)
			}

			defer resp.Body.Close()

			if resp.StatusCode >= 400 {
				fmt.Fprintf(os.Stderr, "Warning: Server returned status %s\n", resp.Status)
			}

			_, err = io.Copy(os.Stdout, resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response: %w", err)
			}

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
