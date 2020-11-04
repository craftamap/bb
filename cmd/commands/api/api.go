package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Headers)

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

			req, err := http.NewRequest(Method, url, bytes.NewBufferString(reqBody))
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}
			req.SetBasicAuth(globalOpts.Username, globalOpts.Password)

			for _, header := range Headers {
				splitted := strings.SplitN(header, "=", 2)
				if len(splitted) == 2 {
					req.Header.Add(splitted[0], splitted[1])
				}
			}

			response, err := client.Do(req)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}

			defer response.Body.Close()

			resBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occured: "), err)
				return
			}
			if response.StatusCode <= 200 || response.StatusCode > 299 {
				fmt.Printf("%s%s%d\n", aurora.Yellow(":: "), aurora.Bold("Status Code: "), response.StatusCode)
			}

			if strings.Contains(response.Header["Content-Type"][0], "json") {
				fmt.Println(string(pretty.Color(pretty.Pretty(resBody), nil)))
			} else {
				fmt.Println(string(resBody))
			}
		},
	}

	apiCmd.Flags().StringVarP(&Method, "method", "X", "GET", "The HTTP method for the request")
	apiCmd.Flags().StringArrayVarP(&Headers, "header", "H", []string{}, "Add an additional HTTP request header")
	rootCmd.AddCommand(apiCmd)
}
