package management

import (
	"bytes"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func uri(domain string, path ...string) string {
	return "https://" + domain + "/" + strings.Join(path, "/")
}

func apiURI(domain string, path ...string) string {
	return "https://" + domain + "/api/v2/" + strings.Join(path, "/")
}

func GetToken(domain string, clientId string, clientSecret string) string {
	q := url.Values{}
	q.Add("grant_type", "client_credentials")
	q.Add("client_id", clientId)
	q.Add("client_secret", clientSecret)
	q.Add("audience", apiURI(domain))

	payload := strings.NewReader(q.Encode())

	req, _ := http.NewRequest("POST", uri(domain, "oauth", "token"), payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}(res.Body)

	body, _ := ioutil.ReadAll(res.Body)
	value := gjson.Get(string(body), "access_token")

	return value.String()
}

func send(method string, uri string, token string, data []byte, params []byte, ignoreStatusCode ...int) (string, error) {
	payload := bytes.NewBuffer(data)

	req, _ := http.NewRequest(
		method,
		uri,
		payload,
	)

	if params != nil {
		q := req.URL.Query()
		result := gjson.ParseBytes(params)
		result.ForEach(func(key, value gjson.Result) bool {
			q.Add(key.String(), value.String())
			return true // keep iterating
		})
		req.URL.RawQuery = q.Encode()
	}

	bearer := fmt.Sprintf("Bearer %s", token)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", bearer)
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}(res.Body)

	isValid := res.StatusCode == http.StatusOK || res.StatusCode == http.StatusAccepted || res.StatusCode == http.StatusCreated || res.StatusCode == http.StatusNoContent

	if !isValid && len(ignoreStatusCode) > 0 {
		isValid = res.StatusCode == ignoreStatusCode[0]
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	if !isValid {
		log.Printf("URI: %s\n", req.URL.String())
		log.Printf("%s\n", bodyBytes)
		return "", fmt.Errorf("exited with %d status code", res.StatusCode)
	}

	log.Printf("exited with %d status code \n", res.StatusCode)
	return string(bodyBytes), nil
}
