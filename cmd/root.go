/*
Copyright © 2021 Brandon Butler <bmbawb@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
	"gitlab.com/brandonbutler/chiabot/internal/release"
)

var (
	token     string
	channelID string
	interval  string
)

var (
	previousLatest string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chiabot",
	Short: "A release bot for the Chia blockchain",

	Run: func(cmd *cobra.Command, args []string) {
		//Check interval duration setup
		intervalWithSuffix := fmt.Sprintf("%ss", interval)
		checkInterval, err := time.ParseDuration(intervalWithSuffix)
		if err != nil {
			log.Fatal(fmt.Sprintf("[ERROR] Could not parse duration. Interval: %s\n%v"), intervalWithSuffix, err)
		}

		//HTTP client setup
		timeout := time.Duration(60 * time.Second)
		client := &http.Client{
			Timeout: timeout,
		}

		//Get latest release for the first time
		previousLatest, err = release.GetLatest(client)
		if err != nil {
			log.Fatal(err)
		}

		//Slack client setup
		bot := slack.New(token)

		for {
			//Get latest release
			latest, err := release.GetLatest(client)
			if err != nil {
				log.Println(err)
				continue
			}

			//Handle condition -- latest release unchanged
			if latest == previousLatest {
				time.Sleep(checkInterval)
				log.Println("[INFO] Iteration found same release")
				continue
			}

			//Get changes
			releaseURL := fmt.Sprintf("https://github.com/Chia-Network/chia-blockchain/releases/tag/%s", latest)
			cl, err := release.GetChanges(client, releaseURL)
			if err != nil {
				log.Println(err)
			}
			logs := compileChangelogs(releaseURL, latest, cl)

			//Post to Slack
			_, _, err = bot.PostMessage(channelID,
				slack.MsgOptionText(logs, true),
				slack.MsgOptionAsUser(true),
				slack.MsgOptionEnableLinkUnfurl(),
			)
			if err != nil {
				log.Printf("[ERROR] Could not post message to channel. Received error: /n%s", err)
			}

			//Update the old latest
			log.Println("[INFO] New latest release discovered: ", latest)
			previousLatest = latest
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	token = os.Getenv("TOKEN")
	channelID = os.Getenv("CHANNEL_ID")
	interval = os.Getenv("CHECK_INTERVAL")
}

func compileChangelogs(url string, ver string, cl release.Changelog) string {
	log := fmt.Sprintf(":rocket::rocket::rocket: New Chia release v%s! :rocket::rocket::rocket: \n%s\n", ver, url)
	log += "\n"
	log += "```"
	if len(cl.Added) > 0 {
		log += "Added\n"
		for _, line := range cl.Added {
			if line == "" {
				continue
			}
			log += fmt.Sprintf(" * %s\n", line)
		}
	}

	if len(cl.Changed) > 0 {
		log += "\n"
		log += "Changed\n"
		for _, line := range cl.Changed {
			if line == "" {
				continue
			}
			log += fmt.Sprintf(" * %s\n", line)
		}
	}

	if len(cl.Fixed) > 0 {
		log += "\n"
		log += "Fixed\n"
		for _, line := range cl.Fixed {
			if line == "" {
				continue
			}
			log += fmt.Sprintf(" * %s\n", line)
		}
	}
	log += "```"

	return log
}
