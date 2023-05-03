package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const notifyUrl = "https://notify-api.line.me/api/notify"

var accessTokens = flag.String("line_access_token", os.Getenv("PLUGIN_LINE_ACCESS_TOKEN"), "line access token")

func main() {
	flag.Parse()

	if *accessTokens == "" {
		log.Fatalln("Missing LINE access token")
	}

	accessTokensArr := strings.Split(*accessTokens, ",")
	wg := sync.WaitGroup{}
	body := formatBody()

	for _, v := range accessTokensArr {
		wg.Add(1)
		go func(accessToken string, body io.Reader) {
			callLineNotify(accessToken, body)
			defer wg.Done()
		}(v, body)
	}

	wg.Wait()

	fmt.Println("Line notify done.")
}

func formatBody() io.Reader {
	data := url.Values{}
	data.Add("message", fmt.Sprintf(`
Repo: %s
Brach: %s
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
		os.Getenv("DRONE_COMMIT_AUTHOR"),
		os.Getenv("DRONE_BUILD_EVENT"),
		os.Getenv("DRONE_COMMIT_MESSAGE"),
		os.Getenv("DRONE_BUILD_NUMBER"),
		os.Getenv("DRONE_BUILD_STATUS"),
		os.Getenv("DRONE_BUILD_LINK"),
		os.Getenv("DRONE_COMMIT_LINK"),
		time.Now().Local().Format("2006-01-02T15:04:05 -07:00:00"),
	))
	return bytes.NewBufferString(data.Encode())
}

func callLineNotify(accessToken string, body io.Reader) {
	req, err := http.NewRequest(http.MethodPost, notifyUrl, body)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))
}
