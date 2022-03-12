package model

import (
	"encoding/json"
	"fmt"
)

type CronLog struct {
	ServerId   int64  `json:"serverId"`
	CronId     int64  `json:"cronId"`
	ExecCode   int    `json:"execCode"`
	Message    string `json:"message"`
	ReportTime string `json:"reportTime"`
}

type CronLogs []CronLog

func (cl CronLog) GetList(pagination Pagination) (CronLogs, error) {
	responseBody, err := Request("/agent/getCronLogs", struct {
		ServerID int64  `json:"serverId"`
		CronID   int64  `json:"cronId"`
		Page     uint64 `json:"page"`
		Rows     uint64 `json:"rows"`
	}{
		ServerID: goployServerID,
		CronID:   cl.CronId,
		Page:     pagination.Page,
		Rows:     pagination.Rows,
	})
	if err != nil {
		return CronLogs{}, err
	}
	type Data struct {
		List CronLogs `json:"list"`
	}
	var data Data
	err = json.Unmarshal(responseBody.Data, &data)
	if err != nil {
		return CronLogs{}, fmt.Errorf("body: %v, err: %s", cl, err.Error())
	}
	return data.List, nil
}

func (cl CronLog) Report() error {
	cl.ServerId = goployServerID
	_, err := Request("/agent/cronReport", cl)
	return err
}
