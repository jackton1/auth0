package management

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	var err error

	domain := os.Getenv("AUTH0_DOMAIN")
	clientId := os.Getenv("AUTH0_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")

	token := GetToken(domain, clientId, clientSecret)

	clientParams, _ := json.Marshal(map[string]interface{}{
		"fields":         "initiate_login_uri",
		"include_fields": true,
	})

	initialClientResponse, err := Client(domain, token, clientId, clientParams)

	var clientResponse map[string]interface{}

	err = json.Unmarshal([]byte(initialClientResponse), &clientResponse)

	if err != nil {
		t.Fatal(err)
	}

	var initialResponse []byte

	if val, ok := clientResponse["initiate_login_uri"]; ok {
		initialResponse, err = json.Marshal(map[string]interface{}{
			"initiate_login_uri": val,
		})

		if err != nil {
			t.Fatal(err)
		}
	} else {
		initialResponse, err = json.Marshal(map[string]interface{}{
			"initiate_login_uri": "",
		})

		if err != nil {
			t.Fatal(err)
		}
	}

	t.Cleanup(func() {
		time.Sleep(time.Second)
		if _, err = ClientUpdate(domain, token, clientId, initialResponse); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		clientData, _ := json.Marshal(map[string]interface{}{
			"initiate_login_uri": "https://test.com",
		})

		_, err = ClientUpdate(domain, token, clientId, clientData)

		if err != nil {
			t.Errorf("Error: %v", err)
		}

		expectedError := fmt.Sprintf("exited with %d status code", http.StatusBadRequest)

		_, err = ClientUpdate(domain, token, clientId, []byte("Invalid"))

		if err != nil && err.Error() != expectedError {
			t.Errorf("Error: %v != %s", err, expectedError)
		}
	})
}
