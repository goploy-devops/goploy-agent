package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

// ResponseBody struct
type ResponseBody struct {
	Code    int
	Message string
}

type RequestData struct {
	ServerId   int64  `json:"serverId"`
	Type       int    `json:"type"`
	Item       string `json:"item"`
	Value      string `json:"value"`
	ReportTime string `json:"reportTime"`
}

// ResponseFail response state type
const (
	ResponseFail = iota
)

const (
	_ = iota
	TypeCPU
	TypeRAM
	TypeLoadavg
	TypeDiskUsage
	TypeDiskIO
	TypeNet
	TypeTcp
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

func Request(data RequestData) error {
	data.ServerId = goployServerID
	_url := fmt.Sprintf("%s%s", goployURL, "/agent/report")
	requestData := new(bytes.Buffer)
	_ = json.NewEncoder(requestData).Encode(data)
	requestStr := requestData.String()
	resp, err := gClient.Post(_url, "application/json", requestData)
	if err != nil {
		return errors.New(_url + " " + err.Error())
	}
	defer resp.Body.Close()
	var responseBody ResponseBody
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return errors.New(_url + ", body: " + requestStr + " err: " + err.Error())
	} else if responseBody.Code > 0 {
		return errors.New(_url + ", body: " + requestStr + " message: " + responseBody.Message)
	}
	return nil
}
