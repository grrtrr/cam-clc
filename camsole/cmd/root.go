package cmd

import (
	"context"
	"os"
	"path"
	"time"

	"github.com/grrtrr/clccam"
	"github.com/spf13/cobra"
)

var (
	// Global client
	client *clccam.Client

	// Flags:
	rootFlags struct {
		url     string        // REST endpoint URL
		token   string        // Bearer Token
		json    bool          // Print JSON response to stdout
		debug   bool          // Print request/response debug to stderr
		timeout time.Duration // Client timeout
	}

	// Top-level command
	Root = &cobra.Command{
		Use: path.Base(os.Args[0]),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var camToken clccam.Token
			var err error

			switch cmd.Name() {
			case "token": // Does not need client initialization
				return
			}

			if rootFlags.token == "" {
				if camToken, err = clccam.LoadToken(); err != nil {
					die("failed to load token: %s", err)
				}
			} else {
				if camToken, err = tokenFromStringOrFile(rootFlags.token); err != nil {
					die("%s", err)
				}
				// Save it in the default location ($HOME/.clc/cam.token)
				camToken.Save()
			}

			// Validate token
			if cl, err := camToken.Claims(); err != nil {
				die("token failed to decode: %s", err)
			} else if cl.Expired() {
				die("%s -- get a new one from https://cam.ctl.io", cl)
			}

			// Client initialization:
			client = camToken.NewClient(
				clccam.HostURL(rootFlags.url),
				clccam.Retryer(3, 1*time.Second, rootFlags.timeout),
				clccam.Context(context.Background()),
				clccam.Debug(rootFlags.debug),
				clccam.JsonResponse(rootFlags.json),
			)
		},
	}
)

func init() {
	var endpointUrl = "cam.ctl.io" // Default endpoint URL

	if u := os.Getenv("CAM_URL"); u != "" {
		endpointUrl = u
	}

	Root.PersistentFlags().StringVarP(&rootFlags.token, "token", "t", os.Getenv("CAM_TOKEN"), "Path or contents of CAM token")
	Root.PersistentFlags().StringVarP(&rootFlags.url, "url", "u", endpointUrl, "REST API endpoint URL")
	Root.PersistentFlags().BoolVarP(&rootFlags.debug, "debug", "d", false, "Print request/response debug output to stderr")
	Root.PersistentFlags().BoolVar(&rootFlags.json, "json", false, "Print JSON response to stdout")
	Root.PersistentFlags().DurationVar(&rootFlags.timeout, "timeout", 180*time.Second, "Client default timeout")
}
