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
	agentLogs := AgentLogs{}
	stmt := DB.Prep("SELECT item, value, time FROM agent_log where type = $type and time BETWEEN $low AND $high;")

	stmt.SetInt64("$type", int64(a.Type))
	stmt.SetText("$low", low)
	stmt.SetText("$high", high)
	agent := Agent{}
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return agentLogs, err
		}
		if !hasRow {
			break
		}

		agent.Item = stmt.GetText("item")
		agent.Value = stmt.GetText("value")
		agent.ReportTime = stmt.GetText("time")

		agentLogs = append(agentLogs, agent)
	}
	stmt.Finalize()
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
	stmt, err := DB.Prepare("INSERT INTO agent_log (type, item, value, time) VALUES ($type, $item, $value, $time);")
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

	if err = stmt.Finalize(); err != nil {
		return err
	}

	return err
}
