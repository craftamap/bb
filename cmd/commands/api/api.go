package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/craftamap/bb/util/logging"

	"github.com/craftamap/bb/cmd/options"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/tidwall/pretty"
)

var (
	Method  string
	Headers []string
)

func Add(rootCmd *cobra.Command, globalOpts *options.GlobalOptions) {
	apiCmd := &cobra.Command{
		Use:   "api <url> [<body>]",
		Short: "Make an authenticated api.bitbucket.org request to the rest 2.0 api",
		Long:  "Make an authenticated api.bitbucket.org request to the rest 2.0 api",
		Annotations: map[string]string{
			"RequiresClient": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			client := http.Client{}

			url := ""
			if len(args) > 0 {
				url = args[0]
			}
			url = "https://api.bitbucket.org/2.0/" + url

			reqBody := ""
			if len(args) == 2 {
				reqBody = args[1]
			}

			fi, _ := os.Stdout.Stat()
			pipeMode := fi.Mode()&os.ModeCharDevice == 0 // if mode == blockdevice (= piped)

			req, err := http.NewRequest(Method, url, bytes.NewBufferString(reqBody))
			if err != nil {
				logging.Error(err)
				return
			}
			req.SetBasicAuth(globalOpts.Client.Username, globalOpts.Client.Password)

			for _, header := range Headers {
				splitted := strings.SplitN(header, "=", 2)
				if len(splitted) == 2 {
					req.Header.Add(splitted[0], splitted[1])
				}
			}
			response, err := client.Do(req)
			if err != nil {
				logging.Error(err)
				return
			}

			defer response.Body.Close()

			resBody, err := io.ReadAll(response.Body)
			if err != nil {
				logging.Error(err)
				return
			}

			if !pipeMode {
				fmt.Printf("%s%s%d\n", aurora.Yellow(":: "), aurora.Bold("Status Code: "), response.StatusCode)
			}

			if strings.Contains(response.Header["Content-Type"][0], "json") {
				if pipeMode { // if mode == blockdevice (= piped)
					fmt.Println(string(pretty.Pretty(resBody)))
				} else {
					fmt.Println(string(pretty.Color(pretty.Pretty(resBody), nil)))
				}
			} else {
				fmt.Println(string(resBody))
			}
		},
	}

	apiCmd.Flags().StringVarP(&Method, "method", "X", "GET", "The HTTP method for the request")
	apiCmd.Flags().StringArrayVarP(&Headers, "header", "H", []string{}, "Add an additional HTTP request header")
	rootCmd.AddCommand(apiCmd)
}
