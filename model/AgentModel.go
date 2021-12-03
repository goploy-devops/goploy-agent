package model

type Agent struct {
	ServerId   int64  `json:"serverId"`
	Type       int    `json:"type"`
	Item       string `json:"item"`
	Value      string `json:"value"`
	ReportTime string `json:"reportTime"`
}

const (
	_ = iota
	TypeCPU
	TypeRAM
	TypeLoadavg
	TypeTcp
	TypePubNet
	TypeLoNet
	TypeDiskUsage
	TypeDiskIO
)

func (a Agent) Request() error {
	a.ServerId = goployServerID
	_, err := Request("/agent/report", a)

	return err
}
