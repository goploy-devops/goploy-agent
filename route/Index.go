package route

import (
	router "github.com/zhenorzz/goploy-agent/core"
)

// Init router
func Init() *router.Router {
	var rt = new(router.Router)
	// rt.Middleware(example)
	// no need to check login
	rt.RegisterWhiteList(map[string]struct{}{
		"/user/login":      {},
	})

	rt.Start()
	return rt
}
