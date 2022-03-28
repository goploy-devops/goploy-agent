package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhenorzz/goploy-agent/config"
	"github.com/zhenorzz/goploy-agent/core"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"zombiezen.com/go/sqlite/sqlitex"
)

// ResponseBody struct
type ResponseBody struct {
	Code    int
	Message string
	Data    json.RawMessage
}

// ResponseSuccess response state type
const (
	ResponseSuccess = 0
)

// DB init when the program start
var goployURL string
var goployServerID int64
var gClient = &http.Client{Timeout: 5 * time.Second}
var DB *sqlitex.Pool

// Init -
func Init() {
	c, err := sqlitex.Open(core.GetDBFile(), 0, 10)
	if err != nil {
		log.Fatal(err)
	}
	DB = c
	if err := createTable(); err != nil {
		log.Fatal(err)
	}
	goployURL = config.Toml.Goploy.ReportURL
	goployServerID = getServerID()
	core.Log(core.INFO, fmt.Sprintf("server id %d", goployServerID))
}

func Shutdown() error {
	return DB.Close()
}

func createTable() error {
	conn := DB.Get(nil)
	defer DB.Put(conn)
	stmt, _, err := conn.PrepareTransient(`
		CREATE TABLE IF NOT EXISTS agent_log (
		  type INTEGER,
		  item TEXT,
		  value TEXT,
		  time TEXT
		);
	`)
	if err != nil {
		return err
	}

	if _, err := stmt.Step(); err != nil {
		return err
	}

	if err := stmt.Finalize(); err != nil {
		return err
	}

	stmt, _, err = conn.PrepareTransient(`
		CREATE INDEX IF NOT EXISTS time_idx ON agent_log (time);
	`)
	if err != nil {
		return err
	}

	if _, err := stmt.Step(); err != nil {
		return err
	}

	if err := stmt.Finalize(); err != nil {
		return err
	}
	return nil
}

func getServerID() int64 {
	if config.Toml.Goploy.UIDType == "id" {
		serverID, err := strconv.ParseInt(config.Toml.Goploy.UID, 10, 64)
		if err != nil {
			core.Log(core.ERROR, fmt.Sprintf("Parse uid to server id error, %s", err.Error()))
			return 0
		}
		return serverID
	} else if config.Toml.Goploy.UIDType == "name" {
		responseBody, err := Request("/agent/getServerID", struct {
			Name string `json:"name"`
		}{Name: config.Toml.Goploy.UID})
		if err != nil {
			core.Log(core.ERROR, fmt.Sprintf("request error, %s", err.Error()))
			return 0
		}

		type Data struct {
			ID int64 `json:"id"`
		}

		var data Data
		err = json.Unmarshal(responseBody.Data, &data)
		if err != nil {
			core.Log(core.ERROR, fmt.Sprintf("Parse response body fail, %s", err.Error()))
			return 0
		}
		return data.ID

	} else if config.Toml.Goploy.UIDType == "host" {
		responseBody, err := Request("/agent/getServerID", struct {
			IP string `json:"ip"`
		}{IP: config.Toml.Goploy.UID})
		if err != nil {
			core.Log(core.ERROR, fmt.Sprintf("request error, %s", err.Error()))
			return 0
		}

		type Data struct {
			ID int64 `json:"id"`
		}

		var data Data
		err = json.Unmarshal(responseBody.Data, &data)
		if err != nil {
			core.Log(core.ERROR, fmt.Sprintf("Parse response body fail, %s", err.Error()))
			return 0
		}
		return data.ID
	}
	return 0
}

var ErrNoReportURL = errors.New("no report url in toml")

func Request(uri string, data interface{}) (ResponseBody, error) {
	responseBody := ResponseBody{}
	if config.Toml.Goploy.ReportURL == "" {
		return responseBody, ErrNoReportURL
	}
	_url := fmt.Sprintf("%s%s", goployURL, uri)
	u, err := url.Parse(_url)
	if err != nil {
		return responseBody, fmt.Errorf("%s, parse error: %s", _url, err.Error())
	}

	requestData := new(bytes.Buffer)
	err = json.NewEncoder(requestData).Encode(data)
	if err != nil {
		return responseBody, fmt.Errorf("%s, request data %+v, json encode error: %s", _url, data, err.Error())
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
		return responseBody, fmt.Errorf("%s, request body: %s, requset err: %s", _url, requestStr, err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return responseBody, fmt.Errorf("%s, request body: %s, http status code: %d", _url, requestStr, resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return responseBody, fmt.Errorf("%s request body: %s, respond body: %s, decode json err: %s", _url, requestStr, string(body), err.Error())
	} else if responseBody.Code != ResponseSuccess {
		return responseBody, fmt.Errorf("%s request body: %s, respond body: %+v, respond code: %d", _url, requestStr, responseBody, responseBody.Code)
	}

	return responseBody, nil
}
