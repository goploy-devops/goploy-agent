package model

type Agent struct {
	ServerId   int64  `json:"serverId"`
	Type       int    `json:"type"`
	Item       string `json:"item"`
	Value      string `json:"value"`
	ReportTime string `json:"reportTime"`
}

type AgentLogs []Agent

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

func (a Agent) GetListBetweenTime(low, high string) (AgentLogs, error) {
	conn := DB.Get(nil)
	defer DB.Put(conn)
	agentLogs := AgentLogs{}
	stmt, _, err := conn.PrepareTransient("SELECT item, value, time FROM agent_log where type = $type and time BETWEEN $low AND $high;")
	if err != nil {
		return agentLogs, err
	}
	stmt.SetInt64("$type", int64(a.Type))
	stmt.SetText("$low", low)
	stmt.SetText("$high", high)

	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return agentLogs, err
		}
		if !hasRow {
			break
		}
		agent := Agent{
			Item:       stmt.GetText("item"),
			Value:      stmt.GetText("value"),
			ReportTime: stmt.GetText("time"),
		}
		agentLogs = append(agentLogs, agent)
	}
	if err := stmt.Finalize(); err != nil {
		return agentLogs, err
	}
	return agentLogs, nil
}

func (a Agent) Report() error {
	a.ServerId = goployServerID
	_, err := Request("/agent/report", a)
	if err == ErrNoReportURL {
		return nil
	}
	return err
}

func (a Agent) Insert() error {
	conn := DB.Get(nil)
	defer DB.Put(conn)
	stmt, err := conn.Prepare("INSERT INTO agent_log (type, item, value, time) VALUES ($type, $item, $value, $time);")
	if err != nil {
		return err
	}
	stmt.SetInt64("$type", int64(a.Type))
	stmt.SetText("$item", a.Item)
	stmt.SetText("$value", a.Value)
	stmt.SetText("$time", a.ReportTime)
	if _, err = stmt.Step(); err != nil {
		return err
	}

	return err
}
