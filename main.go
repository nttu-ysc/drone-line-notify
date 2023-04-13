package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const notifyUrl = "https://notify-api.line.me/api/notify"

func main() {
	var accessToken = os.Getenv("line_access_token")
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
	body := bytes.NewBufferString(data.Encode())

	req, err := http.NewRequest(http.MethodPost, notifyUrl, body)
	if err != nil {
		panic(err)
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
