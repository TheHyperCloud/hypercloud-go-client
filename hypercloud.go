package hypercloud

import (

    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "time"
    "strconv"
    "strings"

    Json "encoding/json"
)

type hypercloud struct{
    access, secret  string
    baseUrl         string

    tokenCache      map[string]interface{}
    client          *http.Client
}

func NewHypercloud(url string, access string, secret string) hypercloud {
    var ret = hypercloud{access, secret, url, nil, nil}
    ret.access = access
    ret.secret = secret
    ret.client = &http.Client{
        Timeout : 25 * time.Second,
    }
    ret.generateToken()
    return ret
}

func (h* hypercloud) Request(method string, url string, data interface{}) (rVal interface{}, err []string){
    //Normalize method
    method = strings.ToUpper(method)
    json, body, status := h._request(method, url, data)
    if status == 401 && strings.Contains(body, "invalid_token") {
        h.tokenCache = nil
        json, body, status = h._request(method, url, data)
    }
    rVal = json
    if 200 <= status && status < 300 {
        err = nil
        return
    } else if status == 401 {
        err = append(err, "Authentication error")
        err = append(err, body)
        return
    } else if status == 403 {
        err = append(err, "Unauthorized error")
        err = append(err, body)
    } else if status == 400 || status == 404 {
        err = append(err, "Invalid request error")
        err = append(err, body)
    } else if status == 422 {
        err = append(err, "Validation error")
        err = append(err, body)
    } else {
        err = append(err, fmt.Sprintf("API Error: %s", strconv.Itoa(status)))
        err = append(err, body)
    }
    return
}

func (h* hypercloud) _request(method string, url string, data interface{}) (json interface{}, body string, status int) {
    if h.tokenCache == nil {
        h.generateToken()
    } else if h.tokenCache["expires"].(int64) <= time.Now().Unix() {
        h.generateToken()
    }

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

    req.Header["Authorization"] = []string{"Bearer " + h.tokenCache["access_token"].(string)}
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

func (h* hypercloud) generateToken() {
    h.tokenCache = nil
    login_url := strings.Join([]string{h.baseUrl, "/oauth/token"}, "")
    var invalid map[string]interface{}
    err := Json.Unmarshal([]byte("{\"access_token\" : \"0000000000000000000000000000000000000000000000000000000000000000\", \"token_tpe\" : \"bearer\", \"expires_in\" : 2, \"refresh_token\" : null, \"scope\" : \"\", \"expires\" : 5}"), &invalid) //Always expires

    form := url.Values{}
    form.Add("grant_type", "client_credentials")
    form.Add("client_id", h.access)
    form.Add("client_secret", h.secret)
    req, err := http.NewRequest("POST", login_url, strings.NewReader(form.Encode()))

    if err != nil {
        h.tokenCache = invalid
        return
    }

    login_response, err := h.client.Do(req)

    if err != nil {
        h.tokenCache = invalid
        return
    }

    defer login_response.Body.Close()
    body, err := ioutil.ReadAll(login_response.Body)

    if err != nil {
        h.tokenCache = invalid
        return
    }

    var rVal map[string]interface{}
    err = Json.Unmarshal(body, &rVal)
    if err != nil {
        h.tokenCache = invalid
        return
    }
    if rVal["access_token"] == nil || len(rVal["access_token"].(string)) != 64 {
        h.tokenCache = invalid
        return
    }
    h.tokenCache = rVal
    h.tokenCache["expires"] = time.Now().Unix() + int64(h.tokenCache["expires_in"].(float64))-60 //Anti-time-skew thingy
    return
}
