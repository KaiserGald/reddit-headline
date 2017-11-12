package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
)

var apiMethods = []string{
	"/hot?limit=5",
	"/r/all/hot?limit=5",
}

type redditoauth struct {
	userInfo
	OauthTokenRequestURL string
	OauthAPIDomain       string
	ScopeID              string
	AccessTokenName      string
}

type userInfo struct {
	UsrName     string
	UsrPassword string
	AppID       string
	AppSecret   string
}

type redditresp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type postinfo struct {
	Title   string
	Sub     string
	Poster  string
	Score   string
	Created string
}

var home bool
var all bool

func handleError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func main() {

	args := os.Args[1:]
	handleArgs(args)

	client := &http.Client{}
	filePath := "./userinfo.json"

	file, err := ioutil.ReadFile(filePath)
	handleError(err)

	var user userInfo
	err = json.Unmarshal(file, &user)
	handleError(err)

	oauth := redditoauth{
		userInfo:             user,
		OauthTokenRequestURL: "https://www.reddit.com/api/v1/access_token",
		OauthAPIDomain:       "https://oauth.reddit.com",
	}

	token := getToken(client, oauth).AccessToken
	topSubs := []postinfo{}
	if home {
		topSubs = getPosts(client, oauth, token, 0)
	} else if all {
		topSubs = getPosts(client, oauth, token, 1)
	}

	printSubs(topSubs)
}

func handleArgs(args []string) {
	err := errors.New("arguments not understood")
	if len(args) == 0 {
		home = true
	} else if args[0] == "all" {
		all = true
	} else {
		log.Fatalln(err)
	}

}

func getToken(c *http.Client, oauth redditoauth) (r redditresp) {
	v := url.Values{}
	v.Set("grant_type", "password")
	v.Set("username", oauth.UsrName)
	v.Set("password", oauth.UsrPassword)
	s := v.Encode()

	req, err := http.NewRequest("POST", oauth.OauthTokenRequestURL, strings.NewReader(s))
	handleError(err)
	req.SetBasicAuth(oauth.AppID, oauth.AppSecret)
	req.Header.Add("User-Agent", "linux:redditHeadline:v0.0.0")

	resp, err := c.Do(req)
	handleError(err)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	handleError(err)

	err = json.Unmarshal(data, &r)
	handleError(err)

	return
}

func getPosts(c *http.Client, oauth redditoauth, token string, m int) (subs []postinfo) {
	req, err := http.NewRequest("GET", oauth.OauthAPIDomain+apiMethods[m], nil)
	handleError(err)

	req.Header.Add("Authorization", "bearer "+token)
	req.Header.Add("User-Agent", "linux:redditHeadline:v0.0.0")

	resp, err := c.Do(req)
	handleError(err)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	handleError(err)

	for i := 0; i < 5; i++ {
		index := strconv.Itoa(i)

		title, _, _, err := jsonparser.Get(data, "data", "children", "["+index+"]", "data", "title")
		handleError(err)

		sub, _, _, err := jsonparser.Get(data, "data", "children", "["+index+"]", "data", "subreddit_name_prefixed")
		handleError(err)

		poster, _, _, err := jsonparser.Get(data, "data", "children", "["+index+"]", "data", "author")
		handleError(err)

		score, _, _, err := jsonparser.Get(data, "data", "children", "["+index+"]", "data", "score")
		handleError(err)

		created, _, _, err := jsonparser.Get(data, "data", "children", "["+index+"]", "data", "created_utc")
		handleError(err)

		c := convertTime(created)
		post := postinfo{
			Title:   string(title),
			Sub:     string(sub),
			Poster:  string(poster),
			Score:   string(score),
			Created: string(c),
		}
		subs = append(subs, post)

	}
	return
}

func convertTime(ts []byte) string {
	t := strings.TrimSuffix(string(ts), ".0")
	i, err := strconv.ParseInt(t, 10, 64)
	handleError(err)

	tm := time.Unix(i, 0)
	td := time.Since(tm)
	return td.String()
}

func printSubs(subs []postinfo) {
	largeSpacer := strings.Repeat("-", 150)
	smallSpacer := strings.Repeat("-", 52)
	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()

	for _, v := range subs {
		fmt.Println(largeSpacer)
		fmt.Printf(" %s\n", v.Title)
		fmt.Println(smallSpacer)
		fmt.Printf(" Subreddit: %-25s Score: %s\n", v.Sub, v.Score)
		fmt.Println(smallSpacer)
		fmt.Printf(" Poster: %-25s  Created: %s ago.\n\n\n", v.Poster, v.Created)
	}

}
