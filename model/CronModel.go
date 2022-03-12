package model

import (
	"encoding/json"
	"fmt"
)

type Cron struct {
	ID          int64  `json:"id"`
	ServerId    int64  `json:"serverId"`
	Expression  string `json:"expression"`
	Command     string `json:"command"`
	SingleMode  uint8  `json:"singleMode"`
	LogLevel    uint8  `json:"logLevel"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
	Editor      string `json:"editor"`
	InsertTime  string `json:"insertTime"`
	UpdateTime  string `json:"updateTime"`
}

type Crons []Cron

func (c Cron) GetList() (Crons, error) {
	c.ServerId = goployServerID
	responseBody, err := Request("/agent/getCronList", c)
	if err != nil {
		return Crons{}, err
	}

	type Data struct {
		List Crons `json:"list"`
	}

	var data Data
	err = json.Unmarshal(responseBody.Data, &data)
	if err != nil {
		return Crons{}, fmt.Errorf("body: %v, err: %s", c, err.Error())
	}
	return data.List, nil

}
