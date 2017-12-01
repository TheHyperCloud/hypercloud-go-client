package hypercloud

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	Json "encoding/json"
)

type hypercloud struct {
	token   string
	baseUrl string

	client *http.Client
}

func ToHypercloud(data interface{}) hypercloud {
	return data.(hypercloud)
}

func NewHypercloud(url string, token string) (hc hypercloud, erro []error) {
	var ret = hypercloud{token, url, nil}
	ret.client = &http.Client{
		Timeout: 25 * time.Second,
	}
	hc = ret
	return
}

func (h *hypercloud) Request(method string, url string, data interface{}) (rVal interface{}, err []error) {
	//Normalize method
	method = strings.ToUpper(method)
	json, body, status := h._request(method, url, data)

	rVal = json
	if 200 <= status && status < 300 {
		err = nil
		return
	} else if status == 401 {
		err = append(err, fmt.Errorf("Authentication error: %s", body))
		return
	} else if status == 403 {
		err = append(err, fmt.Errorf("Unauthorized error: %s", body))
	} else if status == 400 || status == 404 {
		err = append(err, fmt.Errorf("Invalid request error: %s", body))
	} else if status == 422 {
		err = append(err, fmt.Errorf("Validation error: %s", body))
	} else {
		err = append(err, fmt.Errorf("API Error: %s \n%s", strconv.Itoa(status), body))
	}
	return
}

func (h *hypercloud) _request(method string, url string, data interface{}) (json interface{}, body string, status int) {
	url = h.baseUrl + "/api/v1" + url
	var req *http.Request
	if data != nil {
		sendData, err := Json.Marshal(data)
		if err != nil {
			err = Json.Unmarshal([]byte("{\"error\" : \"Invalid data\", \"error_description\" : \"data failed to be marshalled to json\"}"), &json)
			body = data.(string)
			status = 400
			return
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(sendData))
		if err != nil {
			err = Json.Unmarshal([]byte("{\"error\" : \"Invalid data\", \"error_description\" : \"unable to create a new request\"}"), &json)
			body = data.(string)
			status = 400
			return
		}
	} else {
		var err error
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			err = Json.Unmarshal([]byte("{\"error\" : \"Invalid data\", \"error_description\" : \"unable to create a new request\"}"), &json)
			body = data.(string)
			status = 400
			return
		}
	}

	req.Header["Authorization"] = []string{"Bearer " + h.token}
	req.Header["User-agent"] = []string{"Generated Client (golang)"}
	req.Header["Content-type"] = []string{"application/json"}
	req.Header["Accept"] = []string{"application/json"}

	resp, err := h.client.Do(req)
	if err != nil {
		Json.Unmarshal([]byte("{\"error\" : \"Invalid data\", \"error_description\" : \"request failed to complete. Refer to body for details\"}"), &json)
		body = err.Error()
		status = 503
		return
	}
	defer resp.Body.Close()
	mData, err := ioutil.ReadAll(resp.Body)
	err = Json.Unmarshal(mData, &json)
	status = resp.StatusCode
	if err != nil {
		Json.Unmarshal([]byte("{\"error\" : \"Invalid data\", \"error_description\" : \"Unable to decode json\"}"), &json)
		body = err.Error()
		status = 503
		return
	}
	return
}
