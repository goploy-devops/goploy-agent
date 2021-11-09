package task

import (
	"bytes"
	"context"
	"fmt"
	model "github.com/zhenorzz/goploy-agent/Model"
	"github.com/zhenorzz/goploy-agent/core"
	"github.com/zhenorzz/goploy-agent/utils"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var counter int32
var stop = make(chan struct{})

func Init() {
	atomic.AddInt32(&counter, 1)
	go ticker()
}

func ticker() {
	defer atomic.AddInt32(&counter, -1)
	// create ticker
	minute := time.Tick(time.Minute)
	second := time.Tick(time.Second)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for {
			select {
			case <-second:
			case <-stop:
				wg.Done()
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-minute:
				reportCPUUsage()
				reportRAMUsage()
				reportLoadavg()
				reportDisk()
				reportNet()
				reportTcp()
			case <-stop:
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
}

func reportCPUUsage() {
	getCPUSample := func() (idle, total uint64) {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd := exec.Command("cat", "/proc/stat")
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			core.Log(core.ERROR, cmd.String()+" err: "+err.Error()+", detail: "+stderr.String())
			return
		}
		for _, line := range strings.Split(stdout.String(), "\n") {
			fields := strings.Fields(line)
			if fields[0] == "cpu" {
				numFields := len(fields)
				for i := 1; i < numFields; i++ {
					val, err := strconv.ParseUint(fields[i], 10, 64)
					if err == nil {
						total += val // tally up all the numbers to get total ticks
					}

					if i == 4 { // idle is the 5th field in the cpu line
						idle = val
					}
				}
				return
			}
		}
		return
	}
	idle0, total0 := getCPUSample()
	time.Sleep(3 * time.Second)
	idle1, total1 := getCPUSample()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks
	if err := model.Request(model.RequestData{
		Type:       model.TypeCPU,
		Item:       "cpu_usage",
		Value:      fmt.Sprintf("%.2f", cpuUsage),
		ReportTime: time.Now().Format("2006-01-02 15:04"),
	}); err != nil {
		core.Log(core.ERROR, err.Error())
	}
}

func reportRAMUsage() {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("head", "-n", "2", "/proc/meminfo")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		core.Log(core.ERROR, cmd.String()+" err: "+err.Error()+", detail: "+stderr.String())
		return
	}
	total := 0.0
	free := 0.0
	for i, line := range strings.Split(utils.ClearNewline(stdout.String()), "\n") {
		fields := strings.Fields(line)
		if i == 0 {
			val, err := strconv.ParseFloat(fields[1], 64)
			if err == nil {
				total += val // tally up all the numbers to get total ticks
			}
		} else if i == 1 {
			val, err := strconv.ParseFloat(fields[1], 64)
			if err == nil {
				free += val // tally up all the numbers to get total ticks
			}
		}
	}
	ramUsage := 100 * (total - free) / total
	if err := model.Request(model.RequestData{
		Type:       model.TypeRAM,
		Item:       "ram_usage",
		Value:      fmt.Sprintf("%.2f", ramUsage),
		ReportTime: time.Now().Format("2006-01-02 15:04"),
	}); err != nil {
		core.Log(core.ERROR, err.Error())
	}
}

func reportLoadavg() {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("cat", "/proc/loadavg")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		core.Log(core.ERROR, cmd.String()+" err: "+err.Error()+", detail: "+stderr.String())
		return
	}
	procLoadavg := strings.Split(stdout.String(), " ")
	if err := model.Request(model.RequestData{
		Type:       model.TypeLoadavg,
		Item:       "loadavg_1m",
		Value:      procLoadavg[0],
		ReportTime: time.Now().Format("2006-01-02 15:04"),
	}); err != nil {
		core.Log(core.ERROR, err.Error())
	}
	if err := model.Request(model.RequestData{
		Type:       model.TypeLoadavg,
		Item:       "loadavg_5m",
		Value:      procLoadavg[1],
		ReportTime: time.Now().Format("2006-01-02 15:04"),
	}); err != nil {
		core.Log(core.ERROR, err.Error())
	}
	if err := model.Request(model.RequestData{
		Type:       model.TypeLoadavg,
		Item:       "loadavg_15m",
		Value:      procLoadavg[2],
		ReportTime: time.Now().Format("2006-01-02 15:04"),
	}); err != nil {
		core.Log(core.ERROR, err.Error())
	}
}

func reportTcp() {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("wc", "-l", "/proc/net/tcp")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		core.Log(core.ERROR, cmd.String()+" err: "+err.Error()+", detail: "+stderr.String())
		return
	}

	if err := model.Request(model.RequestData{
		Type:       model.TypeTcp,
		Item:       "tcp.total",
		Value:      strings.Fields(utils.ClearNewline(stdout.String()))[0],
		ReportTime: time.Now().Format("2006-01-02 15:04"),
	}); err != nil {
		core.Log(core.ERROR, err.Error())
	}

	stdout.Reset()
	stderr.Reset()
	cmd = exec.Command("grep", "-c", "^ *[0-9]\\+: [0-9A-F: ]\\{27\\} 01 ", "/proc/net/tcp")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		core.Log(core.ERROR, cmd.String()+" err: "+err.Error()+", detail: "+stderr.String())
		return
	}

	if err := model.Request(model.RequestData{
		Type:       model.TypeTcp,
		Item:       "tcp.established",
		Value:      utils.ClearNewline(stdout.String()),
		ReportTime: time.Now().Format("2006-01-02 15:04"),
	}); err != nil {
		core.Log(core.ERROR, err.Error())
	}

}


func reportNet() {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("cat", "/proc/net/dev")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		core.Log(core.ERROR, cmd.String()+" err: "+err.Error()+", detail: "+stderr.String())
		return
	}
	net1 := strings.Split(utils.ClearNewline(stdout.String()), "\n")[2:]

	time.Sleep(1 * time.Second)

	stdout.Reset()
	stderr.Reset()
	cmd = exec.Command("cat", "/proc/net/dev")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		core.Log(core.ERROR, cmd.String()+" err: "+err.Error()+", detail: "+stderr.String())
		return
	}
	net2 := strings.Split(utils.ClearNewline(stdout.String()), "\n")[2:]

	for i, line := range net1 {
		fields1 := strings.Fields(line)
		logType := 0
		if strings.HasPrefix(fields1[0], "eth") {
			logType = model.TypePubNet
		} else if strings.HasPrefix(fields1[0], "lo") {
			logType = model.TypeLoNet
		} else {
			continue
		}

		fields2 := strings.Fields(net2[i])

		in1, _ := strconv.Atoi(fields1[1])
		in2, _ := strconv.Atoi(fields2[1])

		in := in2 - in1

		out1, _ := strconv.Atoi(fields1[9])
		out2, _ := strconv.Atoi(fields2[9])

		out := out2 - out1

		if err := model.Request(model.RequestData{
			Type:       logType,
			Item:       fields1[0][:len(fields1[0])-1] + ".in",
			Value:      strconv.Itoa(in),
			ReportTime: time.Now().Format("2006-01-02 15:04"),
		}); err != nil {
			core.Log(core.ERROR, err.Error())
		}

		if err := model.Request(model.RequestData{
			Type:       logType,
			Item:       fields1[0][:len(fields1[0])-1] + ".out",
			Value:      strconv.Itoa(out),
			ReportTime: time.Now().Format("2006-01-02 15:04"),
		}); err != nil {
			core.Log(core.ERROR, err.Error())
		}
	}
}

func reportDisk() {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("df", "--output=pcent,ipcent,source")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		core.Log(core.ERROR, cmd.String()+" err: "+err.Error()+", detail: "+stderr.String())
		return
	}
	for _, line := range strings.Split(utils.ClearNewline(stdout.String()), "\n")[1:] {
		fields := strings.Fields(line)
		diskName := strings.Join(fields[2:], " ")
		if !strings.HasPrefix(diskName, "/dev/") {
			continue
		}

		diskUsedPcent := fields[0][:len(fields[0])-1]
		diskIUsedPcent := fields[1][:len(fields[1])-1]

		if diskUsedPcent != "" {
			if err := model.Request(model.RequestData{
				Type:       model.TypeDiskUsage,
				Item:       diskName + ".usage",
				Value:      diskUsedPcent,
				ReportTime: time.Now().Format("2006-01-02 15:04"),
			}); err != nil {
				core.Log(core.ERROR, err.Error())
			}
		}

		if diskIUsedPcent != "" {
			if err := model.Request(model.RequestData{
				Type:       model.TypeDiskUsage,
				Item:       diskName + ".inode_usage",
				Value:      diskIUsedPcent,
				ReportTime: time.Now().Format("2006-01-02 15:04"),
			}); err != nil {
				core.Log(core.ERROR, err.Error())
			}
		}
		diskNameWithoutPrefix := strings.TrimPrefix(diskName, "/dev/")

		diskSuffixName := strings.Map(func(r rune) rune {
			if '0' <= r && r <= '9' {
				return -1
			}
			return r
		}, diskNameWithoutPrefix)

		stdout.Reset()
		stderr.Reset()
		cmd = exec.Command("cat", fmt.Sprintf("/sys/block/%s/%s/stat", diskSuffixName, diskNameWithoutPrefix))
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			core.Log(core.ERROR, cmd.String()+" err: "+err.Error()+", detail: "+stderr.String())
			continue
		}

		diskStats1 := strings.Fields(utils.ClearNewline(stdout.String()))

		time.Sleep(1 * time.Second)
		stdout.Reset()
		stderr.Reset()
		cmd = exec.Command("cat", fmt.Sprintf("/sys/block/%s/%s/stat", diskSuffixName, diskNameWithoutPrefix))
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			core.Log(core.ERROR, cmd.String()+" err: "+err.Error()+", detail: "+stderr.String())
			continue
		}

		diskStats2 := strings.Fields(utils.ClearNewline(stdout.String()))

		rIOpms1, _ := strconv.Atoi(diskStats1[0])
		rIOpms2, _ := strconv.Atoi(diskStats2[0])

		rIOps := rIOpms2 - rIOpms1

		wIOpms1, _ := strconv.Atoi(diskStats1[4])
		wIOpms2, _ := strconv.Atoi(diskStats2[4])

		wIOps := wIOpms2 - wIOpms1

		if err := model.Request(model.RequestData{
			Type:       model.TypeDiskIO,
			Item:       diskName + ".read_iops",
			Value:      strconv.Itoa(rIOps),
			ReportTime: time.Now().Format("2006-01-02 15:04"),
		}); err != nil {
			core.Log(core.ERROR, err.Error())
		}

		if err := model.Request(model.RequestData{
			Type:       model.TypeDiskIO,
			Item:       diskName + ".write_iops",
			Value:      strconv.Itoa(wIOps),
			ReportTime: time.Now().Format("2006-01-02 15:04"),
		}); err != nil {
			core.Log(core.ERROR, err.Error())
		}
	}
}


func Shutdown(ctx context.Context) error {
	close(stop)
	ticker := time.NewTicker(10 * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if atomic.LoadInt32(&counter) == 0 {
				return nil
			}
		}
	}
}
