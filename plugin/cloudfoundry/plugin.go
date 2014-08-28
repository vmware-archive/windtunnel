package cloudfoundry

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os/user"
)

type Plugin struct {
}

type Config struct {
	AccessToken string
}

type Instance struct {
	State string
}

func (this *Plugin) Authenticate() string {
	currentUser, _ := user.Current()
	configFile := currentUser.HomeDir + "/.cf/config.json"
	configJson, _ := ioutil.ReadFile(configFile)

	var config Config

	json.Unmarshal(configJson, &config)

	return config.AccessToken
}

func (this *Plugin) Status(token string, app string) []int {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.run.pivotal.io/v2/apps/"+app+"/instances", nil)
	req.Header.Add("Authorization", token)
	req.Header.Add("Host", "example.org")
	req.Header.Add("Cookie", "")
	resp, _ := client.Do(req)

	instanceJson, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var instances map[string]Instance
	err := json.Unmarshal(instanceJson, &instances)
	if err != nil {
    log.Println("Cannot marshall instance status:")
    log.Println("JSON: " + string(instanceJson))
		log.Fatal(err)
	}

	upCount := 0
	for _, instance := range instances {
		if instance.State == "RUNNING" {
			upCount++
		}
	}

	return []int{upCount, len(instances)}
}
