package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ivahaev/go-logger"
	"io/ioutil"
)

var (
	confFile *string
	C        Conf
)

type Conf struct {
	Exe          string   `json:"exe" binding:"required"`
	Params       []string `json:"params`
	DefaultVoice string   `json:"defaultVoice" binding:"required"`
	TmpDir       string   `json:"tmpDir" binding:"required"`
	Port         string   `json:"port" binding:"required"`
}

func mustReadConfig() {
	C = Conf{}
	logger.Info("Will load config from: " + *confFile)
	file, err := ioutil.ReadFile(*confFile)
	if err != nil {
		panic(fmt.Sprintf("Config file read error: %v\n", err.Error()))
	}

	err = json.Unmarshal(file, &C)
	if err != nil {
		panic(fmt.Sprintf("Can't decode config: %v\n", err.Error()))
	}
}

func init() {
	confFile = flag.String("config", "./config.json", "Path to conf file")
	flag.Parse()
	mustReadConfig()
}
