package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
	goployURL = os.Getenv("GOPLOY_URL")
	goployServerID, _ = strconv.ParseInt(os.Getenv("GOPLOY_SERVER_ID"), 10, 64)
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
	requestData := new(bytes.Buffer)
	_ = json.NewEncoder(requestData).Encode(data)
	requestStr := requestData.String()
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
