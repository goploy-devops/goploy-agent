package route

import (
	"github.com/zhenorzz/goploy-agent/controller"
	router "github.com/zhenorzz/goploy-agent/core"
	"net/http"
)

// Init router
func Init() *router.Router {
	var rt = new(router.Router)
	// rt.Middleware(example)

	rt.Add("/general", http.MethodGet, controller.Controller{}.General)
	rt.Add("/loadavg", http.MethodGet, controller.Controller{}.Loadavg)
	rt.Add("/ram", http.MethodGet, controller.Controller{}.RAM)
	rt.Add("/cpu", http.MethodGet, controller.Controller{}.CPU)
	rt.Add("/net", http.MethodGet, controller.Controller{}.Net)
	rt.Add("/diskUsage", http.MethodGet, controller.Controller{}.DiskUsage)
	rt.Add("/diskIOStat", http.MethodGet, controller.Controller{}.DiskIOStat)
	rt.Add("/cronList", http.MethodGet, controller.Controller{}.CronList)
	rt.Add("/cronLogs", http.MethodGet, controller.Controller{}.CronLogs)

	rt.Start()
	return rt
}
