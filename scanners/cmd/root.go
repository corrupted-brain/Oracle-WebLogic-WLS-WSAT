// Copyright © 2018 Kevin Kirsche <kev.kirsche[at]gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kkirsche/CVE-2017-10271/scanners/libcve201710271"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	config libcve201710271.Config
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cve-2017-10271",
	Short: "Scan for the CVE-2017-10271 vulnerability",
	Long: fmt.Sprintf(`A purpose built scanner for detecting CVE-2017-10271. Starts a web
server on the LPORT and then logs any host which contacts it, as they are
vulnerable.

Example usage:
./CVE-2017-10271.release.%s.amd64.linux -s "10.10.10.10" -t "$(pwd)/targets.txt -o output_file.txt -v --all-urls"

Example targets.txt:
http://pwned.com:7001/
https://pwnedalso.com:8002/
`, BuildVersion),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if config.Verbose {
			logrus.SetLevel(logrus.InfoLevel)
		} else {
			logrus.SetLevel(logrus.WarnLevel)
		}

		if config.Lport < 1 || config.Lport > 65535 {
			logrus.Errorln("Listening port must be greater than 0 and less than 65536. Exiting...")
			return
		}

		if config.Lhost == "" {
			logrus.Errorln("Listening host IP address or hostname is required. Exiting...")
			return
		}

		if config.TargetFile == "" {
			logrus.Errorln("Target file is required. Exiting...")
			return
		}

		if config.OutputFile != "" {
			f, err := os.OpenFile(config.OutputFile, os.O_WRONLY|os.O_CREATE, 0755)
			if err != nil {
				logrus.WithError(err).Errorln("Failed to open file for writing")
			}
			logrus.SetOutput(f)
		}

		libcve201710271.Banner(config)

		logrus.Infof("Starting webserver on port %d to catch vulnerable hosts", config.Lport)
		go func() {
			http.HandleFunc("/cve-2017-10271", vulnHandler)
			http.ListenAndServe(fmt.Sprintf(":%d", config.Lport), vulnLog(http.DefaultServeMux))
		}()

		f, err := os.Open(config.TargetFile)
		if err != nil {
			logrus.WithError(err).Errorln("Failed to open target file.")
			return
		}
		defer f.Close()

		targetCh := make(chan libcve201710271.TargetHost)

		m := &sync.Mutex{}
		for w := 1; w <= config.Threads; w++ {
			go libcve201710271.Worker(w, m, targetCh)
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			rhost := strings.TrimSpace(scanner.Text())
			rhost = strings.TrimRight(rhost, "/")

			var urls []string
			urls = libcve201710271.DefaultURLs
			if config.AllURLs {
				urls = libcve201710271.AllURLs
			}

			for _, url := range urls {
				xmlPayload := libcve201710271.GenerateCheckPayload(config.Lhost, config.Lport, rhost, url)
				th := libcve201710271.TargetHost{
					R: rhost,
					P: xmlPayload,
					U: url,
				}
				targetCh <- th
			}
		}

		if err := scanner.Err(); err != nil {
			close(targetCh)
			logrus.Fatal(err)
		}
		close(targetCh)

		logrus.Infoln("Sleeping for 10 seconds in case we have any stragglers...")
		time.Sleep(time.Duration(config.WaitTime) * time.Second)
	},
}

func vulnLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, _ := url.QueryUnescape(r.URL.Query().Get("target"))
		if config.OutputFile != "" {
			fmt.Printf("[VULNERABLE] Remote Address: %s | From Target: %s | Method: %s\n", r.RemoteAddr, t, r.Method)
		}
		logrus.Warnf("[VULNERABLE] Remote Address: %s | From Target: %s | Method: %s", r.RemoteAddr, t, r.Method)
		handler.ServeHTTP(w, r)
	})
}

func vulnHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "WARNING! You are vulnerable to CVE-2017-10271")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().StringVarP(&config.Lhost, "listening-host", "s", "", "The IP of this machine's public interface")
	RootCmd.Flags().IntVarP(&config.Lport, "listening-port", "l", 4444, "The port to listen for vulnerable responses")
	RootCmd.Flags().StringVarP(&config.TargetFile, "target-file", "t", "", "File with list of targets in http(s)://HOSTNAME:PORT format")
	RootCmd.Flags().BoolVarP(&config.Verbose, "verbose", "v", false, "Enable verbose mode (Print who is being scanned")
	RootCmd.Flags().BoolVarP(&config.AllURLs, "all-urls", "u", false, "Check for all possible vulnerable URL suffixes")
	RootCmd.Flags().StringVarP(&config.OutputFile, "output-file", "o", "", "File to output results to")
	RootCmd.Flags().IntVarP(&config.Threads, "threads", "a", 10, "Number of threads to use while scanning")
	RootCmd.Flags().IntVarP(&config.WaitTime, "wait-time", "w", 20, "Seconds to wait after we complete sending payloads")
}
