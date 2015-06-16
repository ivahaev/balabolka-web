package main

import (
	"encoding/json"
	"fmt"
	"github.com/ivahaev/balabolka-web/config"
	"github.com/ivahaev/balabolka-web/utils/hash"
	"github.com/ivahaev/go-logger"
	"net/http"
	"os/exec"
	"strings"
)

func synthHandler(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	if text == "" {
		http.Error(w, "No text param", 400)
		return
	}
	voice := r.URL.Query().Get("voice")
	if voice == "" {
		voice = config.Config.DefaultVoice
	}
	params := r.URL.Query().Get("params")
	fileName := config.Config.TmpDir + `\` + hash.New(text) + ".wav"
	params += strings.Join(config.Config.Params, " ")
	cmd := exec.Command("cmd", "/C "+config.Config.Exe+` -w `+fileName+` -n `+voice+` -t "`+text+`" `+params)
	res, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, "Server error. Can't exec: "+err.Error()+"\n "+res, 500)
		return
	}
	http.ServeFile(w, r, fileName)
}

func voicesHandler(w http.ResponseWriter, r *http.Request) {
	out, err := exec.Command("cmd", "/C", config.Config.Exe, "-l").Output()
	if err != nil {
		http.Error(w, "Server error. Can't exec: "+err.Error(), 500)
		return
	}
	resultString := string(out)
	resultArr := strings.Split(resultString, "\n")
	for i := len(resultArr) - 1; i >= 0; i-- {
		if strings.TrimSpace(resultArr[i]) == "" {
			resultArr = append(resultArr[:i], resultArr[i+1:]...)
		} else {
			resultArr[i] = strings.TrimSpace(resultArr[i])
		}
	}
	result, err := json.Marshal(resultArr)
	if err != nil {
		http.Error(w, "Server error. Can't read result: "+err.Error(), 500)
		return
	}
	fmt.Fprintf(w, string(result))
}

func main() {
	logger.Info("App started")
	logger.Debug(config.Config)
	http.HandleFunc("/voices", voicesHandler)
	http.HandleFunc("/synth", synthHandler)
	http.ListenAndServe(":"+config.Config.Port, nil)
}
