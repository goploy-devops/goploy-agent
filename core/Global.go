package core

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

var (
	AssetDir string
	Gwg      sync.WaitGroup
)

// GetAssetDir if env = 'production' return absolute else return relative
func GetAssetDir() string {
	if AssetDir != "" {
		return AssetDir
	}
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		panic(err)
	}
	app, err := filepath.Abs(file)
	if err != nil {
		panic(err)
	}
	i := strings.LastIndex(app, "/")
	if i < 0 {
		i = strings.LastIndex(app, "\\")
	}
	if i < 0 {
		panic(err)
	}
	return app[0 : i+1]
}

func GetDBFile() string {
	return path.Join(GetAssetDir(), "goploy-agent.db")
}

func GetConfigFile() string {
	return path.Join(GetAssetDir(), "goploy-agent.toml")
}
