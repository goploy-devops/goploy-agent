package model


type AgentModel struct {}

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


func (AgentModel)Request(data RequestData) error {
	return Request("/agent/report", data)
}
