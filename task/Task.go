package task

import (
	"bytes"
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/zhenorzz/goploy-agent/config"
	"github.com/zhenorzz/goploy-agent/core"
	"github.com/zhenorzz/goploy-agent/model"
	"github.com/zhenorzz/goploy-agent/utils"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var task = gocron.NewScheduler(time.UTC)

var JobList = map[[sha1.Size]byte]struct {
	Cron model.Cron
	job  *gocron.Job
}{}

func Init() {
	if config.Toml.Goploy.ReportURL == "" {
		core.Log(core.WARNING, "no report url detect, turn to standalone mode")
		return
	}
	_, _ = task.Every(1).Minute().WaitForSchedule().SingletonMode().Do(reportCPUUsage)
	_, _ = task.Every(1).Minute().WaitForSchedule().SingletonMode().Do(reportRAMUsage)
	_, _ = task.Every(1).Minute().WaitForSchedule().SingletonMode().Do(reportLoadavg)
	_, _ = task.Every(1).Minute().WaitForSchedule().SingletonMode().Do(reportDisk)
	_, _ = task.Every(1).Minute().WaitForSchedule().SingletonMode().Do(reportDiskIO)
	_, _ = task.Every(1).Minute().WaitForSchedule().SingletonMode().Do(reportNet)
	_, _ = task.Every(1).Minute().WaitForSchedule().SingletonMode().Do(reportTcp)
	_, _ = task.Every(1).Minute().SingletonMode().Do(getCron)
	task.StartAsync()
}

func Add(cron model.Cron) error {
	s := task.CronWithSeconds(cron.Expression)
	if cron.SingleMode == 1 {
		s = s.SingletonMode()
	}
	job, err := s.WaitForSchedule().Do(func() {
		var message string
		var execErrorMessage string
		var execCode int
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		scriptArgs := strings.Fields(cron.Command)
		cmd := exec.Command(scriptArgs[0], scriptArgs[1:]...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				execCode = exitError.ExitCode()
			} else {
				execCode = 1 // for any unknown error
			}
			execErrorMessage = err.Error()
		}
		stderrMessage := utils.ClearNewline(stderr.String())
		stdoutMessage := utils.ClearNewline(stdout.String())

		if execErrorMessage != "" || cron.LogLevel > 0 {
			if execErrorMessage != "" {
				message += "Exec error:\n" + execErrorMessage + "\n"
			}
			message += "Stdout:\n" + stdoutMessage + "\n"
			if cron.LogLevel > 1 {
				message += "Stderr:\n" + stderrMessage + "\n"
			}
			err := model.CronLog{
				CronId:     cron.ID,
				Message:    message,
				ExecCode:   execCode,
				ReportTime: time.Now().Format("2006-01-02 15:04:05"),
			}.Report()
			if err != nil {
				core.Log(core.ERROR, err.Error())
			}
		}

	})
	if err != nil {
		return err
	}
	JobList[sha1.Sum([]byte(fmt.Sprintf("%v", cron)))] = struct {
		Cron model.Cron
		job  *gocron.Job
	}{cron, job}
	return nil
}

func Shutdown(ctx context.Context) error {
	task.Stop()
	ticker := time.NewTicker(10 * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if !task.IsRunning() {
				return nil
			}
		}
	}
}

func getCron() {
	if crons, err := (model.Cron{}).GetList(); err != nil {
		core.Log(core.ERROR, fmt.Sprintf("get cron list failed, error: %s", err.Error()))
	} else {
		cronHashMap := map[[sha1.Size]byte]int64{}
		for _, cron := range crons {
			cronHash := sha1.Sum([]byte(fmt.Sprintf("%v", cron)))
			cronHashMap[cronHash] = cron.ID
			if _, ok := JobList[cronHash]; !ok {
				if err := Add(cron); err != nil {
					core.Log(core.ERROR, fmt.Sprintf("id: %d, add to cron task error: %s", cron.ID, err.Error()))
				} else {
					core.Log(core.TRACE, fmt.Sprintf("id: %d(%x), add to cron task", cron.ID, cronHash))
				}
			}
		}

		for cronHash, o := range JobList {
			if _, ok := cronHashMap[cronHash]; !ok {
				task.RemoveByReference(o.job)
				delete(JobList, cronHash)
				core.Log(core.TRACE, fmt.Sprintf("id: %d(%x), delete from cron task", o.Cron.ID, cronHash))
			}
		}
	}
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
	if err := (model.Agent{
		Type:       model.TypeCPU,
		Item:       "cpu_usage",
		Value:      fmt.Sprintf("%.2f", cpuUsage),
		ReportTime: time.Now().Format("2006-01-02 15:04"),
	}).Request(); err != nil {
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
	if err := (model.Agent{
		Type:       model.TypeRAM,
		Item:       "ram_usage",
		Value:      fmt.Sprintf("%.2f", ramUsage),
		ReportTime: time.Now().Format("2006-01-02 15:04"),
	}).Request(); err != nil {
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
	if err := (model.Agent{
		Type:       model.TypeLoadavg,
		Item:       "loadavg_1m",
		Value:      procLoadavg[0],
		ReportTime: time.Now().Format("2006-01-02 15:04"),
	}).Request(); err != nil {
		core.Log(core.ERROR, err.Error())
	}
	if err := (model.Agent{
		Type:       model.TypeLoadavg,
		Item:       "loadavg_5m",
		Value:      procLoadavg[1],
		ReportTime: time.Now().Format("2006-01-02 15:04"),
	}).Request(); err != nil {
		core.Log(core.ERROR, err.Error())
	}
	if err := (model.Agent{
		Type:       model.TypeLoadavg,
		Item:       "loadavg_15m",
		Value:      procLoadavg[2],
		ReportTime: time.Now().Format("2006-01-02 15:04"),
	}).Request(); err != nil {
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

	if err := (model.Agent{
		Type:       model.TypeTcp,
		Item:       "tcp.total",
		Value:      strings.Fields(utils.ClearNewline(stdout.String()))[0],
		ReportTime: time.Now().Format("2006-01-02 15:04"),
	}).Request(); err != nil {
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

	if err := (model.Agent{
		Type:       model.TypeTcp,
		Item:       "tcp.established",
		Value:      utils.ClearNewline(stdout.String()),
		ReportTime: time.Now().Format("2006-01-02 15:04"),
	}).Request(); err != nil {
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

		if err := (model.Agent{
			Type:       logType,
			Item:       fields1[0][:len(fields1[0])-1] + ".in",
			Value:      strconv.Itoa(in),
			ReportTime: time.Now().Format("2006-01-02 15:04"),
		}).Request(); err != nil {
			core.Log(core.ERROR, err.Error())
		}

		if err := (model.Agent{
			Type:       logType,
			Item:       fields1[0][:len(fields1[0])-1] + ".out",
			Value:      strconv.Itoa(out),
			ReportTime: time.Now().Format("2006-01-02 15:04"),
		}).Request(); err != nil {
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
			if err := (model.Agent{
				Type:       model.TypeDiskUsage,
				Item:       diskName + ".usage",
				Value:      diskUsedPcent,
				ReportTime: time.Now().Format("2006-01-02 15:04"),
			}).Request(); err != nil {
				core.Log(core.ERROR, err.Error())
			}
		}

		if diskIUsedPcent != "" {
			if err := (model.Agent{
				Type:       model.TypeDiskUsage,
				Item:       diskName + ".inode_usage",
				Value:      diskIUsedPcent,
				ReportTime: time.Now().Format("2006-01-02 15:04"),
			}).Request(); err != nil {
				core.Log(core.ERROR, err.Error())
			}
		}
	}
}

func reportDiskIO() {
	disks, err := ioutil.ReadDir("/sys/block/")
	if err != nil {
		core.Log(core.ERROR, "err: "+err.Error())
		return
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	for _, disk := range disks {
		diskName := disk.Name()
		stdout.Reset()
		stderr.Reset()
		cmd := exec.Command("cat", fmt.Sprintf("/sys/block/%s/stat", diskName))
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
		cmd = exec.Command("cat", fmt.Sprintf("/sys/block/%s/stat", diskName))
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

		if err := (model.Agent{
			Type:       model.TypeDiskIO,
			Item:       diskName + ".read_iops",
			Value:      strconv.Itoa(rIOps),
			ReportTime: time.Now().Format("2006-01-02 15:04"),
		}).Request(); err != nil {
			core.Log(core.ERROR, err.Error())
		}

		if err := (model.Agent{
			Type:       model.TypeDiskIO,
			Item:       diskName + ".write_iops",
			Value:      strconv.Itoa(wIOps),
			ReportTime: time.Now().Format("2006-01-02 15:04"),
		}).Request(); err != nil {
			core.Log(core.ERROR, err.Error())
		}
	}
}
