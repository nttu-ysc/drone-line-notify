package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

const (
	notifyUrl = "https://notify-api.line.me/api/notify"
	version   = "v1.2.2"
)

var accessTokens *[]string
var message *string
var help = `
    ____                              ___                              __  _ ____     
   / __ \_________  ____  ___        / (_)___  ___        ____  ____  / /_(_) __/_  __
  / / / / ___/ __ \/ __ \/ _ \______/ / / __ \/ _ \______/ __ \/ __ \/ __/ / /_/ / / /
 / /_/ / /  / /_/ / / / /  __/_____/ / / / / /  __/_____/ / / / /_/ / /_/ / __/ /_/ / 
/_____/_/   \____/_/ /_/\___/     /_/_/_/ /_/\___/     /_/ /_/\____/\__/_/_/  \__, /  
                                                                             /____/   

Author: Shun Cheng
GitHub: https://github.com/nttu-ysc/drone-line-notify
`
var rootCmd = cobra.Command{
	Use:     "drone-line-notify",
	Version: version,
	Long:    help,
	Run: func(cmd *cobra.Command, args []string) {
		if len(*accessTokens) == 0 || (*accessTokens)[0] == "" {
			cmd.Help()
			return
		}

		body := formatBody()

		fmt.Println(help)
		var wg sync.WaitGroup
		for _, token := range *accessTokens {
			wg.Add(1)
			go func(t string) {
				defer wg.Done()
				sendLineNotify(t, body)
			}(token)
		}

		wg.Wait()
	},
}

func main() {
	accessTokens = rootCmd.Flags().StringSliceP("line_access_token", "l", []string{os.Getenv("PLUGIN_LINE_ACCESS_TOKEN")}, "LINE access token")
	message = rootCmd.Flags().StringP("message", "m", "", "message")
	rootCmd.Execute()
}

func formatBody() io.Reader {
	data := url.Values{}
	if *message == "" {
		*message = fmt.Sprintf(`
Repo: %s
Branch: %s
Status: %s
Author: %s
Event: %s
Commit Message: %s
Drone Build number: %s
Drone Build status: %s
Build: %s
Changes: %s
Current time: %s`,
			os.Getenv("DRONE_REPO"),
			os.Getenv("DRONE_COMMIT_BRANCH"),
			os.Getenv("DRONE_BUILD_STATUS"),
			os.Getenv("DRONE_COMMIT_AUTHOR"),
			os.Getenv("DRONE_BUILD_EVENT"),
			os.Getenv("DRONE_COMMIT_MESSAGE"),
			os.Getenv("DRONE_BUILD_NUMBER"),
			os.Getenv("DRONE_BUILD_STATUS"),
			os.Getenv("DRONE_BUILD_LINK"),
			os.Getenv("DRONE_COMMIT_LINK"),
			time.Now().Local().Format("2006-01-02T15:04:05 -07:00:00"),
		)
	}
	data.Add("message", *message)
	return strings.NewReader(data.Encode())
}

// sendLineNotify sends a Line notification using the specified access token and request body.
func sendLineNotify(accessToken string, body io.Reader) {
	req, err := http.NewRequest(http.MethodPost, notifyUrl, body)
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Line notification response:", string(responseBody))
}
