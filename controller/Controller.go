package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/schema"
	"github.com/zhenorzz/goploy-agent/core"
	"github.com/zhenorzz/goploy-agent/model"
	"github.com/zhenorzz/goploy-agent/task"
	"github.com/zhenorzz/goploy-agent/utils"
	"gopkg.in/go-playground/validator.v9"
	"os/exec"
	"strconv"
	"strings"
)

// Controller struct
type Controller struct{}

var decoder = schema.NewDecoder()

func decodeJson(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	if err := core.Validate.Struct(v); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Translate(core.Trans))
		}
	}
	return nil
}

func decodeQuery(data map[string][]string, v interface{}) error {
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(v, data)
	if err != nil {
		return err
	}
	if err := core.Validate.Struct(v); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Translate(core.Trans))
		}
	}
	return nil
}

// General info
// uname -mrs
// cat /etc/os-release
// nproc --all
// hostname
// uptime -p
func (Controller) General(*core.Goploy) *core.Response {
	type GeneralInfo struct {
		KernelVersion string `json:"kernelVersion"`
		OS            string `json:"os"`
		Cores         string `json:"cores"`
		Hostname      string `json:"hostname"`
		Uptime        string `json:"uptime"`
	}

	var generalInfo GeneralInfo
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("uname", "-rs")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		generalInfo.KernelVersion = "Unknown"
	} else {
		generalInfo.KernelVersion = utils.ClearNewline(stdout.String())
	}
	stdout.Reset()
	stderr.Reset()
	cmd = exec.Command("bash", "-c", "cat /etc/os-release | grep \"PRETTY_NAME\" | awk -F\\\" '{print $2}'")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		generalInfo.OS = "Unknown"
	} else {
		generalInfo.OS = utils.ClearNewline(stdout.String())
	}
	stdout.Reset()
	stderr.Reset()
	cmd = exec.Command("nproc", "--all")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		println(err.Error())
		generalInfo.Cores = "Unknown"
	} else {
		generalInfo.Cores = utils.ClearNewline(stdout.String())
	}
	stdout.Reset()
	stderr.Reset()
	cmd = exec.Command("hostname")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		generalInfo.Hostname = "Unknown"
	} else {
		generalInfo.Hostname = utils.ClearNewline(stdout.String())
	}
	stdout.Reset()
	stderr.Reset()
	cmd = exec.Command("uptime", "-p")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		generalInfo.Uptime = "Unknown"
	} else {
		generalInfo.Uptime = utils.ClearNewline(stdout.String())
	}

	return &core.Response{
		Data: generalInfo,
	}
}

// Loadavg info
// cat /proc/loadavg
func (Controller) Loadavg(*core.Goploy) *core.Response {
	type LoadavgInfo struct {
		Avg   string `json:"avg"`
		Avg5  string `json:"avg5"`
		Avg15 string `json:"avg15"`
		Cores string `json:"cores"`
	}

	var loadavgInfo LoadavgInfo
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("cat", "/proc/loadavg")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
	} else {
		procLoadavg := strings.Split(stdout.String(), " ")
		loadavgInfo.Avg = procLoadavg[0]
		loadavgInfo.Avg5 = procLoadavg[1]
		loadavgInfo.Avg15 = procLoadavg[2]
	}
	stdout.Reset()
	stderr.Reset()
	cmd = exec.Command("nproc", "--all")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		println(err.Error())
		loadavgInfo.Cores = ""
	} else {
		loadavgInfo.Cores = utils.ClearNewline(stdout.String())
	}

	return &core.Response{
		Data: loadavgInfo,
	}
}

// RAM info
// cat /proc/meminfo
func (Controller) RAM(*core.Goploy) *core.Response {
	type RAMInfo struct {
		Total int `json:"total"`
		Free  int `json:"free"`
	}

	var ramInfo RAMInfo
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", "head -n 2 /proc/meminfo | awk -F \" \" '{print $2}'")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
	} else {
		procMeminfo := strings.Split(utils.ClearNewline(stdout.String()), "\n")
		ramInfo.Total, err = strconv.Atoi(procMeminfo[0])
		if err == nil {
			ramInfo.Total *= 1000
		}

		ramInfo.Free, err = strconv.Atoi(procMeminfo[1])
		if err == nil {
			ramInfo.Free *= 1000
		}
	}

	return &core.Response{
		Data: ramInfo,
	}
}

// CPU info
// cat /proc/stat
func (Controller) CPU(*core.Goploy) *core.Response {
	var cpuList [][]string
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("cat", "/proc/stat")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {

	} else {
		for _, line := range strings.Split(utils.ClearNewline(stdout.String()), "\n") {
			if strings.Contains(line, "cpu") {
				cpuList = append(cpuList, strings.Fields(line))
			}
		}
	}

	return &core.Response{
		Data: cpuList,
	}
}

// Net info
// cat /proc/net/dev
func (Controller) Net(*core.Goploy) *core.Response {
	var netList [][]string
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("cat", "/proc/net/dev")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
	} else {
		for _, line := range strings.Split(utils.ClearNewline(stdout.String()), "\n")[2:] {
			fields := strings.Fields(line)
			if !strings.HasPrefix(fields[0], "eth") && !strings.HasPrefix(fields[0], "lo") {
				continue
			}
			netList = append(netList, fields)
		}

	}
	return &core.Response{
		Data: netList,
	}
}

// DiskUsage info
// df -h  --output=source,size,used,avail,pcent,target,itotal,iused,iavail,ipcent,fstype
func (Controller) DiskUsage(*core.Goploy) *core.Response {
	// Filesystem            Size  Used Avail Use% Mounted on Inodes IUsed IFree IUse% Type
	type DiskUsageInfo struct {
		Filesystem string `json:"filesystem"`
		Size       string `json:"size"`
		Used       string `json:"used"`
		Avail      string `json:"avail"`
		UsedPcent  string `json:"usedPcent"`
		MountedOn  string `json:"mountedOn"`
		Inodes     string `json:"inodes"`
		IUsed      string `json:"iUsed"`
		IFree      string `json:"iFree"`
		IUsedPcent string `json:"iUsedPcent"`
		Type       string `json:"type"`
	}
	var diskUsageList []DiskUsageInfo

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("df", "-h", "--output=size,used,avail,pcent,target,itotal,iused,iavail,ipcent,fstype,source")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {

	} else {
		for _, line := range strings.Split(utils.ClearNewline(stdout.String()), "\n")[1:] {
			field := strings.Fields(line)
			if !strings.HasPrefix(field[10], "/dev/") {
				continue
			}
			diskUsageList = append(diskUsageList, DiskUsageInfo{
				Size:       field[0],
				Used:       field[1],
				Avail:      field[2],
				UsedPcent:  field[3],
				MountedOn:  field[4],
				Inodes:     field[5],
				IUsed:      field[6],
				IFree:      field[7],
				IUsedPcent: field[8],
				Type:       field[9],
				Filesystem: strings.Join(field[10:], " "),
			})
		}
	}

	return &core.Response{
		Data: diskUsageList,
	}
}

// DiskIOStat info
// iostat -xdk
func (Controller) DiskIOStat(*core.Goploy) *core.Response {
	var diskIOList [][]string
	var header []string
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("iostat", "-xdk")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
	} else {
		lines := strings.Split(utils.ClearNewline(stdout.String()), "\n")
		header = strings.Fields(lines[2])
		for _, line := range lines[3:] {
			diskIOList = append(diskIOList, strings.Fields(line))
		}
	}

	return &core.Response{
		Data: struct {
			Header []string   `json:"header"`
			List   [][]string `json:"list"`
		}{
			Header: header,
			List:   diskIOList,
		},
	}
}

func (Controller) CronList(*core.Goploy) *core.Response {
	var crons model.Crons

	for _, o := range task.JobList {
		crons = append(crons, o.Cron)
	}

	return &core.Response{Data: crons}
}

func (Controller) CronLogs(gp *core.Goploy) *core.Response {
	type ReqData struct {
		Page uint64 `schema:"page" validate:"gt=0"`
		Rows uint64 `schema:"rows" validate:"gt=0"`
	}

	var reqData ReqData
	if err := decodeQuery(gp.URLQuery, &reqData); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}

	id, err := strconv.ParseInt(gp.URLQuery.Get("id"), 10, 64)
	if err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	}
	if cronList, err := (model.CronLog{CronId: id}).GetList(reqData.Page, reqData.Rows); err != nil {
		return &core.Response{Code: core.Error, Message: err.Error()}
	} else {
		return &core.Response{Data: cronList}
	}
}
