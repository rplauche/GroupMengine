package GroupMengine

import (
	"appengine"
	"appengine/urlfetch"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	http.HandleFunc("/newmsg", sendMessage)
}

type NewMessage struct {
	Id          string `json:"id"`
	Source_guid string `json:"source_guid"`
	Created_at  int    `json:"created_at"`
	User_id     string `json:"user_id"`
	Groupd_id   string `json:"group_id"`
	Name        string `json:"name"`
	Avatar_url  string `json:"avatar_url"`
	Text        string `json:"text"`
}

type HandleFunc func(client *http.Client, term string, w http.ResponseWriter) string

func randomWiki(client *http.Client, term string, w http.ResponseWriter) string {
	res, err := client.Get("http://en.wikipedia.org/wiki/Special:Random")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(w, string(resp))
	return string(resp)
}

func generic(client *http.Client, term string, w http.ResponseWriter) string {
	return "I LOVE YOU"
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)

	var handlers = map[string]HandleFunc{
		"/music":        spotifySearch,
		"/groupmengine": generic,
		"/saywhat":      redditSearch,
	}

	p := make([]byte, r.ContentLength)
	_, err := r.Body.Read(p)
	if err == nil {
		var msg NewMessage
		err1 := json.Unmarshal(p, &msg)
		if err1 == nil {
			cmd := msg.Text
			if strings.Index(cmd, "/") == 0 {
				cmdType := strings.Split(cmd, " ")[0]
				cmdBody := strings.TrimSpace(strings.Replace(cmd, cmdType, "", -1))
				cmdBody = strings.Replace(cmdBody, " ", "+", -1)
				handler, ok := handlers[cmdType]
				if ok {
					form := make(url.Values)
					form.Add("bot_id", bots[msg.Groupd_id])
					form.Add("text", handler(client, cmdBody, w))
					client.PostForm("https://api.groupme.com/v3/bots/post", form)
				}
			}
		}
	}
	r.Body.Close()
}