package management

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestAction(t *testing.T) {
	var err error

	domain := os.Getenv("AUTH0_DOMAIN")
	clientId := os.Getenv("AUTH0_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")

	token := GetToken(domain, clientId, clientSecret)

	actionName := "test-action"
	trigger := "post-login"

	actionData, _ := json.Marshal(map[string]interface{}{
		"name": actionName,
		"supported_triggers": []map[string]string{
			{
				"id":      trigger,
				"version": "v2",
			},
		},
		"code": "console.log(\"Test\")",
		"dependencies": []map[string]string{
			{
				"name":    "auth0",
				"version": "2.36.1",
			},
		},
		"runtime": "node16",
		"secrets": []map[string]string{
			{
				"name":  "DOMAIN",
				"value": domain,
			},
			{
				"name":  "CLIENT_ID",
				"value": clientId,
			},
			{
				"name":  "CLIENT_SECRET",
				"value": clientSecret,
			},
		},
	})

	actionResponse, err := ActionCreate(domain, token, actionData)

	if err != nil {
		t.Fatal(err)
	}

	actionId := gjson.Get(actionResponse, "id").String()

	t.Cleanup(func() {
		time.Sleep(time.Second)
		if _, err = ActionDelete(domain, token, actionId); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Create Error", func(t *testing.T) {
		_, err = ActionCreate(domain, token, actionData, http.StatusConflict)

		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		actionParams, _ := json.Marshal(map[string]string{
			"actionName": actionName,
			"triggerId":  trigger,
		})

		actionResponse, err = Actions(domain, token, actionParams)

		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		actionIds := gjson.Get(actionResponse, "actions.#.id").Array()

		if len(actionIds) > 0 {
			id := actionIds[0].String()

			if id != actionId {
				t.Errorf("%s != %s", id, actionId)
			}
		}

		_, err = ActionUpdate(domain, token, actionId, actionData)

		if err != nil {
			t.Errorf("Error: %v", err)
		}

		expectedError := fmt.Sprintf("exited with %d status code", http.StatusBadRequest)

		_, err = ActionUpdate(domain, token, actionId, []byte("Invalid"))

		if err != nil && err.Error() != expectedError {
			t.Errorf("Error: %v != %s", err, expectedError)
		}
	})

	t.Run("Deploy", func(t *testing.T) {
		if _, err = ActionDeploy(domain, token, actionId, actionName, trigger); err != nil {
			t.Errorf("Error: %v", err)
		}
	})
}
