package management

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"log"
	"time"
)

func Actions(domain string, token string, params []byte) (string, error) {
	log.Printf("Retrieving '%s' action...", gjson.GetBytes(params, "actionName"))
	return send("GET", apiURI(domain, "actions", "actions"), token, nil, params)
}

func ActionCreate(domain string, token string, data []byte, ignoreStatusCode ...int) (string, error) {
	log.Print("Creating action...")
	return send("POST", apiURI(domain, "actions", "actions"), token, data, nil, ignoreStatusCode...)
}

func ActionUpdate(domain string, token string, id string, data []byte) (string, error) {
	log.Print("Updating action...")
	return send("PATCH", apiURI(domain, "actions", "actions", id), token, data, nil)
}

func ActionDeploy(domain string, token string, id string, actionName string, trigger string) (string, error) {
	log.Print("Checking if action is built...")

	// Check 10 times and sleep for 30 seconds
	for i := 0; i < 10; i++ {
		actionParams, _ := json.Marshal(map[string]string{
			"actionName": actionName,
			"triggerId":  trigger,
		})
		actions, err := Actions(domain, token, actionParams)

		if err != nil {
			return "", err
		}

		if gjson.Get(actions, "total").Int() == 0 {
			return "", fmt.Errorf("action '%s' not found", actionName)
		}

		// Filter the actions to get the one that matches the id
		action := gjson.Get(actions, "actions.#(id=="+id+")")

		if !action.Exists() {
			return "", fmt.Errorf("action '%s' not found", id)
		}

		if action.Get("status").String() == "built" {
			log.Printf("Action is '%s', deploying...", action.Get("status").String())
			break
		}
		log.Printf("%d of 10 attempts", i+1)
		log.Printf("Action is '%s', waiting 30 seconds for it to be built...", action.Get("status").String())
		delay := 30 * time.Second
		time.Sleep(delay)
	}

	return send("POST", apiURI(domain, "actions", "actions", id, "deploy"), token, nil, nil)
}

func ActionDelete(domain string, token string, id string) (string, error) {
	log.Print("Deleting action...")
	return send("DELETE", apiURI(domain, "actions", "actions", id), token, nil, nil)
}
