package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/zhenorzz/goploy-agent/core"
	"github.com/zhenorzz/goploy-agent/utils"
	"gopkg.in/go-playground/validator.v9"
	"os/exec"
	"strconv"
	"strings"
)

// Controller struct
type Controller struct {
}

func verify(data []byte, v interface{}) error {
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

// General info
// uname -mrs
// cat /etc/os-release
// nproc --all
// hostname
// uptime -p
func (Controller) General(gp *core.Goploy) *core.Response {
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
func (Controller) Loadavg(gp *core.Goploy) *core.Response {
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
func (Controller) RAM(gp *core.Goploy) *core.Response {
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
func (Controller) CPU(gp *core.Goploy) *core.Response {
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
func (Controller) Net(gp *core.Goploy) *core.Response {
	var netList [][]string
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("cat", "/proc/net/dev")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
	} else {
		for _, line := range strings.Split(utils.ClearNewline(stdout.String()), "\n")[2:] {
			netList = append(netList, strings.Fields(line))
		}

	}
	return &core.Response{
		Data: netList,
	}
}

// DiskUsage info
// df -h  --output=source,size,used,avail,pcent,target,itotal,iused,iavail,ipcent,fstype
func (Controller) DiskUsage(gp *core.Goploy) *core.Response {
	var diskUsageInfo string
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("df", "-h", "--output=source,size,used,avail,pcent,target,itotal,iused,iavail,ipcent,fstype")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		diskUsageInfo = err.Error()
	} else {
		diskUsageInfo = utils.ClearNewline(stdout.String())
	}

	return &core.Response{
		Data: diskUsageInfo,
	}
}

// DiskIOStat info
// iostat -xdk
func (Controller) DiskIOStat(gp *core.Goploy) *core.Response {
	var diskIOStatInfo string
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("iostat", "-xdk")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		diskIOStatInfo = err.Error()
	} else {
		diskIOStatInfo = utils.ClearNewline(stdout.String())
	}

	return &core.Response{
		Data: diskIOStatInfo,
	}
}
