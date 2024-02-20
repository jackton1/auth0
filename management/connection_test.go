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

func TestConnection(t *testing.T) {
	var err error

	domain := os.Getenv("AUTH0_DOMAIN")
	clientId := os.Getenv("AUTH0_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")

	token := GetToken(domain, clientId, clientSecret)

	connectionData, _ := json.Marshal(map[string]interface{}{
		"name":                 "Test",
		"display_name":         "Username-Password-Authentication",
		"strategy":             "auth0",
		"is_domain_connection": false,
		"enabled_clients":      []string{clientId},
	})

	connectionResponse, err := ConnectionCreate(domain, token, connectionData)

	if err != nil {
		t.Fatal(err)
	}

	connectionId := gjson.Get(connectionResponse, "id").String()

	t.Cleanup(func() {
		time.Sleep(time.Second)
		if _, err = ConnectionDelete(domain, token, connectionId); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("List", func(t *testing.T) {
		connectionParams, _ := json.Marshal(map[string]interface{}{
			"name":     "Test",
			"strategy": "auth0",
		})

		connectionResponse, err = Connections(domain, token, connectionParams)

		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		connectionIds := gjson.Get(connectionResponse, "#.id").Array()

		if len(connectionIds) > 0 {
			id := connectionIds[0].String()

			if id != connectionId {
				t.Errorf("%s != %s", id, connectionId)
			}
		}
	})

	t.Run("Update", func(t *testing.T) {
		updateConnectionData, _ := json.Marshal(map[string]interface{}{
			"display_name":         "Username-Password-Authentication",
			"is_domain_connection": false,
			"enabled_clients":      []string{clientId},
		})

		_, err = ConnectionUpdate(domain, token, connectionId, updateConnectionData)

		if err != nil {
			t.Errorf("Error: %v", err)
		}

		expectedError := fmt.Sprintf("exited with %d status code", http.StatusBadRequest)

		_, err = ConnectionUpdate(domain, token, connectionId, []byte("Invalid"))

		if err != nil && err.Error() != expectedError {
			t.Errorf("Error: %v != %s", err, expectedError)
		}
	})
}
