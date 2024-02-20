package management

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestUsers(t *testing.T) {
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
	connectionName := gjson.Get(connectionResponse, "name").String()

	userData, _ := json.Marshal(map[string]interface{}{
		"blocked":        false,
		"email":          "test@test.com",
		"user_id":        "1",
		"connection":     connectionName,
		"email_verified": true,
		"password":       "##random-password-123##",
	})

	userResponse, err := UserCreate(domain, token, userData)

	if err != nil {
		t.Fatal(err)
	}

	userId := gjson.Get(userResponse, "user_id").String()
	userEmail := gjson.Get(userResponse, "email").String()

	t.Cleanup(func() {
		time.Sleep(time.Second)
		if _, err = ConnectionDelete(domain, token, connectionId); err != nil {
			t.Fatal(err)
		}
		if _, err = UserDelete(domain, token, userId); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("List", func(t *testing.T) {
		userParams, _ := json.Marshal(map[string]interface{}{
			"page":           1,
			"per_page":       100,
			"include_totals": true,
			"search_engine":  "v3",
			"q":              "email.domain:\"govfetch.com\"",
		})

		_, err := Users(domain, token, userParams)

		if err != nil {
			t.Errorf("Error: %v", err)
		}
	})

	t.Run("ByEmail", func(t *testing.T) {
		userParams, _ := json.Marshal(map[string]string{
			"email": userEmail,
		})

		_, err = UsersByEmail(domain, token, userParams)

		if err != nil {
			t.Errorf("Error: %v", err)
		}

		expectedError := fmt.Sprintf("exited with %d status code", http.StatusBadRequest)

		_, err = UsersByEmail(domain, token, []byte("Invalid"))

		if err != nil && err.Error() != expectedError {
			t.Errorf("Error: %v != %s", err, expectedError)
		}
	})
}
