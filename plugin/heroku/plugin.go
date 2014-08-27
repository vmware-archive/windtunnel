package heroku

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

type Plugin struct {
}

type dyno struct {
	State string
}

func (this dyno) String() string {
	return this.State
}

func (this *Plugin) Authenticate() string {
	out, err := exec.Command("heroku", "auth:token").Output()

	if err != nil {
		log.Fatal(err)
	}

	withColon := ":" + string(out)
	return base64.StdEncoding.EncodeToString([]byte(withColon))
}

func (this *Plugin) Status(token string, app string) []int {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.heroku.com/apps/"+app+"/dynos", nil)
	req.Header.Add("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Add("Authorization", token)
	resp, _ := client.Do(req)

	dynoJson, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var dynos []dyno
	err := json.Unmarshal(dynoJson, &dynos)
	if err != nil {
		log.Fatal(err)
	}

	upCount := 0
	for _, dyno := range dynos {
		if dyno.State == "up" {
			upCount++
		}
	}

	return []int{upCount, len(dynos)}
}
