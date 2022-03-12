package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhenorzz/goploy-agent/config"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ResponseBody struct
type ResponseBody struct {
	Code    int
	Message string
	Data    json.RawMessage
}

// Pagination struct
type Pagination struct {
	Page  uint64 `json:"page"`
	Rows  uint64 `json:"rows"`
	Total uint64 `json:"total"`
}

// ResponseSuccess response state type
const (
	ResponseSuccess = 0
)

// DB init when the program start
var goployURL string
var goployServerID int64
var gClient = &http.Client{Timeout: 5 * time.Second}

// Init -
func Init() {
	goployURL = config.Toml.Goploy.URL
	goployServerID = config.Toml.Goploy.ServerID
}

// PaginationFrom param return pagination struct
func PaginationFrom(param url.Values) (Pagination, error) {
	page, err := strconv.ParseUint(param.Get("page"), 10, 64)
	if err != nil {
		return Pagination{}, errors.New("invalid page")
	}
	rows, err := strconv.ParseUint(param.Get("rows"), 10, 64)
	if err != nil {
		return Pagination{}, errors.New("invalid rows")
	}
	pagination := Pagination{Page: page, Rows: rows}
	return pagination, nil
}

func Request(uri string, data interface{}) (ResponseBody, error) {
	_url := fmt.Sprintf("%s%s", goployURL, uri)
	u, err := url.Parse(_url)
	if err != nil {
		return ResponseBody{}, fmt.Errorf("%s, parse error: %s", _url, err.Error())
	}

	requestData := new(bytes.Buffer)
	err = json.NewEncoder(requestData).Encode(data)
	if err != nil {
		return ResponseBody{}, fmt.Errorf("%s, request data %+v, json encode error: %s", _url, data, err.Error())
	}
	requestStr := requestData.String()
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	unsignedStr := requestStr + timestamp + config.Toml.Goploy.Key
	h := sha256.New()
	h.Write([]byte(unsignedStr))
	q := u.Query()
	q.Set("timestamp", timestamp)
	q.Set("sign", base64.URLEncoding.EncodeToString(h.Sum(nil)))
	u.RawQuery = q.Encode()
	_url = u.String()
	resp, err := gClient.Post(_url, "application/json", requestData)
	if err != nil {
		return ResponseBody{}, fmt.Errorf("%s, request body: %s, requset err: %s", _url, requestStr, err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ResponseBody{}, fmt.Errorf("%s, request body: %s, http status code: %d", _url, requestStr, resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var responseBody ResponseBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return responseBody, fmt.Errorf("%s request body: %s, respond body: %s, decode json err: %s", _url, requestStr, string(body), err.Error())
	} else if responseBody.Code != ResponseSuccess {
		return responseBody, fmt.Errorf("%s request body: %s, respond body: %+v, respond code: %d", _url, requestStr, responseBody, responseBody.Code)
	}

	return responseBody, nil
}
