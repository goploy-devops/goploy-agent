package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/zhenorzz/goploy-agent/config"
	"github.com/zhenorzz/goploy-agent/core"
	"github.com/zhenorzz/goploy-agent/model"
	"github.com/zhenorzz/goploy-agent/route"
	"github.com/zhenorzz/goploy-agent/task"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"
)

var (
	help bool
	v    bool
	s    string
)

const appVersion = "1.3.0"

func init() {
	flag.StringVar(&core.AssetDir, "asset-dir", "", "default: ./")
	flag.StringVar(&s, "s", "", "stop")
	flag.BoolVar(&help, "help", false, "list available subcommands and some concept guides")
	flag.BoolVar(&v, "version", false, "show goploy-agent version")
	// 改变默认的 Usage
	flag.Usage = usage
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	if v {
		println(appVersion)
		return
	}
	handleClientSignal()
	println(`
   ______            __           
  / ____/___  ____  / /___  __  __
 / / __/ __ \/ __ \/ / __ \/ / / /
/ /_/ / /_/ / /_/ / / /_/ / /_/ / 
\____/\____/ .___/_/\____/\__, /  
          /_/            /____/   ` + appVersion + "\n")
	install()
	config.Create(core.GetConfigFile())
	pid := strconv.Itoa(os.Getpid())
	_ = ioutil.WriteFile(path.Join(core.GetAssetDir(), "goploy-agent.pid"), []byte(pid), 0755)
	println("Start at " + time.Now().String())
	println("goploy-agent -h for more help")
	println("Current pid:    " + pid)
	println("Config Loaded:  " + core.GetConfigFile())
	println("UID type:       " + config.Toml.Goploy.UIDType)
	println("UID:            " + config.Toml.Goploy.UID)
	println("Env:            " + config.Toml.Env)
	println("Report to:      " + config.Toml.Goploy.ReportURL)
	println("Log:            " + config.Toml.Log.Path)
	if config.Toml.Web.Port != "" {
		println("Listen:         " + config.Toml.Web.Port)
	}
	core.CreateValidator()
	model.Init()
	route.Init()
	task.Init()
	go checkUpdate()

	// server
	srv := http.Server{
		Addr: ":" + config.Toml.Web.Port,
	}
	core.Gwg.Add(1)
	go func() {
		defer core.Gwg.Done()
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		println("Received the signal: " + (<-c).String())
		println("Server is trying to shutdown, wait for a minute")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			fmt.Printf("Task shutdown failed, err: %s\n", err.Error())
		}
		println("Server shutdown gracefully")

		println("Task is trying to shutdown, wait for a minute")
		if err := task.Shutdown(ctx); err != nil {
			fmt.Printf("Task shutdown failed, err: %s\n", err.Error())
		}
		println("Task shutdown gracefully")

		println("SQLite is trying to shutdown, wait for a minute")
		if err := model.Shutdown(); err != nil {
			fmt.Printf("SQLite shutdown failed, err: %s\n", err.Error())
		}
		println("SQLite shutdown gracefully")
	}()
	if config.Toml.Web.Port != "" {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("ListenAndServe: ", err.Error())
		}
	}
	core.Gwg.Wait()
	_ = os.Remove(path.Join(core.GetAssetDir(), "goploy-agent.pid"))
	println("Goroutine shutdown gracefully")
	println("Success")
	return
}

func install() {
	_, err := os.Stat(core.GetConfigFile())
	if err == nil || os.IsExist(err) {
		println("The configuration file already exists, no need to reinstall (if you need to reinstall, please delete the .env file, then restart.)")
		return
	} else {
		println("The configuration file is not exists, please copy goploy-agent.example.toml to goploy-agent.toml")
		os.Exit(1)
	}
}

func handleClientSignal() {
	switch s {
	case "stop":
		pidStr, err := ioutil.ReadFile(path.Join(core.GetAssetDir(), "goploy-agent.pid"))
		if err != nil {
			log.Fatal("handle stop, ", err.Error(), ", may be the server not start")
		}
		pid, _ := strconv.Atoi(string(pidStr))
		process, err := os.FindProcess(pid)
		if err != nil {
			log.Fatal("handle stop, ", err.Error(), ", may be the server not start")
		}
		err = process.Signal(syscall.SIGTERM)
		if err != nil {
			log.Fatal("handle stop, ", err.Error())
		}
		os.Exit(1)
	}
}

func checkUpdate() {
	resp, err := http.Get("https://api.github.com/repos/zhenorzz/goploy-agent/releases/latest")
	if err != nil {
		println("Check failed")
		println(err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println("Check failed")
		println(err.Error())
		return
	}
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		println("Check failed")
		println(err.Error())
		return
	}

	if _, ok := result["tag_name"]; ok {
		tagName := result["tag_name"].(string)
		tagVer, err := version.NewVersion(tagName)
		if err != nil {
			println("Check version error")
			println(err.Error())
			return
		}
		currentVer, _ := version.NewVersion(appVersion)
		if tagVer.GreaterThan(currentVer) {
			println("New release available")
			println(result["html_url"].(string))
		}
	}
}
