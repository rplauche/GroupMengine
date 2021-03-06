package GroupMengine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type ImgurImage struct {
	Link     string
	Nsfw     bool
	Is_album bool
}

type ImgurData struct {
	Images []ImgurImage
}

type ImgurResponse struct {
	Data []ImgurImage
}

func selectImage(images []ImgurImage) string {
	rand.Seed(time.Now().Unix())
	numImgs := len(images)
	if numImgs > 0 {
		for i := 0; i < 20; i++ {
			randimg := images[rand.Intn(numImgs)]
			if randimg.Is_album != true && randimg.Nsfw != true {
				return randimg.Link
			}
		}
	}
	return "Error No Image Found"
}

func imgurSearch(client *http.Client, term string, w http.ResponseWriter) string {
	// curl -H "Authorization: Client-ID CLIENT_ID_HERE" https://api.imgur.com/3/gallery/hot/viral/0.json
	searchUrl := fmt.Sprintf("https://api.imgur.com/3/gallery/search?q=%s&q_any", term)
	req, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil {
		fmt.Fprint(w, "get setup failed")
		return "Error No Image Found"
	}
	req.Header.Add("Authorization", imgurkey)

	res, err := client.Do(req)
	if err != nil {
		fmt.Fprint(w, "get failed")
		return "Error No Image Found"
	}

	resp, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Fprint(w, "bad resp body")
		return "Error No Image Found"
	}

	var imgResponse ImgurResponse
	err = json.Unmarshal(resp, &imgResponse)

	return selectImage(imgResponse.Data)
}

func imgurRandom(client *http.Client, term string, w http.ResponseWriter) string {
	// curl -H "Authorization: Client-ID CLIENT_ID_HERE" https://api.imgur.com/3/gallery/hot/viral/0.json
	searchUrl := fmt.Sprintf("https://api.imgur.com/3/gallery/hot/viral/0.json")
	req, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil {
		fmt.Fprint(w, "get setup failed")
		return ""
	}
	req.Header.Add("Authorization", imgurkey)

	res, err := client.Do(req)
	if err != nil {
		fmt.Fprint(w, "get failed")
		return ""
	}

	resp, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Fprint(w, "bad resp body")
		return ""
	}

	var imgResponse ImgurResponse
	err = json.Unmarshal(resp, &imgResponse)

	return selectImage(imgResponse.Data)
}
