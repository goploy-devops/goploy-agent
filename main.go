package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/joho/godotenv"
	"github.com/zhenorzz/goploy-agent/core"
	"github.com/zhenorzz/goploy-agent/model"
	"github.com/zhenorzz/goploy-agent/route"
	"github.com/zhenorzz/goploy-agent/task"
	"github.com/zhenorzz/goploy-agent/utils"
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

const appVersion = "1.0.0"

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
	_ = godotenv.Load(core.GetEnvFile())
	pid := strconv.Itoa(os.Getpid())
	port := os.Getenv("PORT")
	_ = ioutil.WriteFile(path.Join(core.GetAssetDir(), "goploy-agent.pid"), []byte(pid), 0755)
	println("Start at " + time.Now().String())
	println("goploy-agent -h for more help")
	println("Current pid   : " + pid)
	println("Server id     : " + os.Getenv("GOPLOY_SERVER_ID"))
	println("Config Loaded : " + core.GetEnvFile())
	println("Report to     : " + os.Getenv("GOPLOY_URL"))
	println("Log           : " + os.Getenv("LOG_PATH"))
	if port != "" {
		println("Server running at http://localhost:" + port)
	}
	core.CreateValidator()
	model.Init()
	route.Init()
	task.Init()
	go checkUpdate()

	// server
	srv := http.Server{
		Addr: ":" + os.Getenv("PORT"),
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
	}()
	if port != "" {
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
	_, err := os.Stat(core.GetEnvFile())
	if err == nil || os.IsExist(err) {
		println("The configuration file already exists, no need to reinstall (if you need to reinstall, please delete the .env file, then restart.)")
		return
	}
	println("Installation guide ↓")
	inputReader := bufio.NewReader(os.Stdin)
	println("Installation guidelines (Enter to confirm input)")
	println("Please enter the absolute path of the log directory(default stdout):")
	logPath, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	logPath = utils.ClearNewline(logPath)
	if len(logPath) == 0 {
		logPath = "stdout"
	}
	println("Please enter the goploy server id (number):")
	serverID, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	serverID = utils.ClearNewline(serverID)
	if len(serverID) == 0 {
		log.Fatal("You must enter the goploy server id.")
	}

	println("Please enter the goploy url (like http://localhost):")
	goployURL, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	goployURL = utils.ClearNewline(goployURL)
	if len(goployURL) == 0 {
		log.Fatal("You must enter the goploy url.")
	}

	println("Please enter the listening port(default turn off web ui):")
	port, err := inputReader.ReadString('\n')
	if err != nil {
		panic("There were errors reading, exiting program.")
	}
	port = utils.ClearNewline(port)
	envContent := "# when you edit its value, you need to restart\n"
	envContent += "ENV=production\n"
	envContent += fmt.Sprintf("GOPLOY_URL=%s\n", goployURL)
	envContent += fmt.Sprintf("GOPLOY_SERVER_ID=%s\n", serverID)
	envContent += fmt.Sprintf("LOG_PATH=%s\n", logPath)
	envContent += fmt.Sprintf("PORT=%s\n", port)
	println("Start writing configuration file...")
	file, err := os.Create(core.GetEnvFile())
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(envContent)
	println("Write configuration file completed")
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
