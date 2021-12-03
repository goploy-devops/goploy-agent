package model

type CronModel struct {
	ServerId   int64  `json:"serverId"`
	ExecRes    int    `json:"status"`
	Message    string `json:"message"`
	ReportTime string `json:"reportTime"`
}

func (c CronModel) GetList() error {
	return Request("/cron/getList", c)
}

func (c CronModel) Report() error {
	c.ServerId = goployServerID
	return Request("/cron/report", c)
}
