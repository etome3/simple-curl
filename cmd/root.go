package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

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

			&cli.StringFlag{
				Name:    "data",
				Aliases: []string{"d"},
				Usage:   "HTTP POST data",
			},

			&cli.StringSliceFlag{
				Name:    "header",
				Aliases: []string{"H"},
				Usage:   "Pass custom header(s) to server",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			url := c.Args().First()
			if url == "" {
				return fmt.Errorf("error: you must provide a URL")
			}

			method := c.String("method")
			data := c.String("data")

			if data != "" && !c.IsSet("method") {
				method = "POST"
			}

			var bodyReader io.Reader

			if data != "" {
				bodyReader = strings.NewReader(data)
			}

			req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
			if err != nil {
				return fmt.Errorf("failed to create request: %w", err)
			}

			for _, h := range c.StringSlice("header") {
				parts := strings.SplitN(h, ":", 2)

				if len(parts) != 2 {
					return fmt.Errorf("header '%s' has wrong format, expect 'Key: Value'", h)
				}

				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				req.Header.Add(key, value)
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
