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

func (cl CronLog) GetList(page, rows uint64) (CronLogs, error) {
	responseBody, err := Request("/agent/getCronLogs", struct {
		ServerID int64  `json:"serverId"`
		CronID   int64  `json:"cronId"`
		Page     uint64 `json:"page"`
		Rows     uint64 `json:"rows"`
	}{
		ServerID: goployServerID,
		CronID:   cl.CronId,
		Page:     page,
		Rows:     rows,
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
	if err == ErrNoReportURL {
		return nil
	}
	return err
}

func (cl CronLog) Insert() error {
	stmt, err := DB.Prepare("INSERT INTO cron_log (type, item, value, time) VALUES ($type, $item, $value, $item);")
	if err != nil {
		return err
	}
	stmt.SetInt64("$type", cl.CronId)
	stmt.SetText("$item", cl.Message)
	stmt.SetText("$value", cl.Message)
	stmt.SetText("$time", cl.ReportTime)

	if _, err = stmt.Step(); err != nil {
		return err
	}

	if err = stmt.Finalize(); err != nil {
		return err
	}

	return err
}
