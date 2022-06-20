package apilib

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

type CognitoResp struct {
	Token string `json:"access_token"`
}

func doRequest(method, url string, headersMap, queryMap, body map[string]string, bodyEncoding string) (*http.Response, error) {
	client := http.Client{}
	var req *http.Request
	var err error
	query := ""
	bodyStr := ""
	header := ""

	switch method {
	case "GET":
		query = createQueryURL(queryMap)
	case "POST":
		header, bodyStr, err = createBody(body, bodyEncoding)

		if err != nil {
			return nil, err
		}

		if header != "" {
			headersMap["content-type"] = header
		}
	}

	req, err = http.NewRequest(method, url+query, strings.NewReader(bodyStr))

	if err != nil {
		return nil, err
	}

	for key, value := range headersMap {
		req.Header.Add(key, value)
	}

	fmt.Println("Doing a request with method", method, "to", url+query, body)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func createQueryURL(queryMap map[string]string) string {
	query := "?"

	first := true
	for key, value := range queryMap {
		if first {
			query += key + "=" + value
			first = false
			continue
		}
		query += "&" + key + "=" + value
	}
	return query
}

func createBody(body map[string]string, encoding string) (string, string, error) {
	bodyStr := ""
	header := ""

	switch encoding {
	case "json":
		jsonString, err := json.Marshal(body)

		if err != nil {
			return "", "", err
		}

		bodyStr = string(jsonString)
		bodyStr = strings.Replace(bodyStr, "\"[", "[", -1)
		bodyStr = strings.Replace(bodyStr, "]\"", "]", -1)
	case "x-www-form-urlecoded":
		first := true
		for key, value := range body {
			if first {
				bodyStr += key + "=" + value
				first = false
				continue
			}
			bodyStr += "&" + key + "=" + value
		}
	case "raw":
		for _, value := range body {
			bodyStr += value
		}
	case "form-data":
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		for key, value := range body {
			err := writer.WriteField(key, value)

			if err != nil {
				return "", "", err
			}
		}
		err := writer.Close()
		if err != nil {
			return "", "", err
		}
		bodyStr = payload.String()
		header = writer.FormDataContentType()
	default:
		break
	}
	return header, bodyStr, nil
}

//Authenticator holds the definition for connect to an enpoint
type Authenticator struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	ClientID string `json:"clientID"`
	Secret   string `json:"secret"`
	URL      string `json:"URLToken"`
	Body     map[string]string
	Query    map[string]string
}

func (auth *Authenticator) SetParameters(usrname string, pwd string, clid string, scrt string, url string, b, q map[string]string) {

	auth.UserName = usrname
	auth.Password = pwd
	auth.ClientID = clid
	auth.Secret = scrt
	auth.URL = url
	auth.Body = b
	auth.Query = q
}

//GetToken from the OAuth2 bearer
func (auth *Authenticator) GetToken() (string, error) {
	//"Basic base64Encode(clid:secretID)"
	basicToken := "Basic " + b64.StdEncoding.EncodeToString([]byte(auth.ClientID+":"+auth.Secret))
	headers := map[string]string{
		"Authorization": basicToken,
		"Content-Type":  "application/x-www-form-urlencoded",
	}
	log.Println(headers)
	resp, err := doRequest("POST", auth.URL, headers, nil, auth.Body, "x-www-form-urlecoded")
	if err != nil {
		panic(err)
	}

	cognitoResp := &CognitoResp{}
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err = json.NewDecoder(resp.Body).Decode(cognitoResp)
	if err != nil {
		panic(err)
	}

	// Do something with the Person struct...

	return cognitoResp.Token, nil
}
