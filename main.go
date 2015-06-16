package main

import (
	"encoding/json"
	"fmt"
	"github.com/ivahaev/balabolka-web/config"
	"github.com/ivahaev/balabolka-web/utils/hash"
	"github.com/ivahaev/go-logger"
	"io/ioutil"
	"net/http"
	"os"
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
	baseFileName := config.Config.TmpDir + `\` + hash.New(text)
	wavFileName := baseFileName + ".wav"
	txtFileName := baseFileName + ".txt"
	err := ioutil.WriteFile(txtFileName, []byte(text), 0644)
	if err != nil {
		http.Error(w, "Server error. Can't write file: "+err.Error(), 500)
		return
	}
	params += strings.Join(config.Config.Params, " ")
	command := "/C " + config.Config.Exe + ` -w ` + wavFileName + ` -n ` + voice + ` -f ` + txtFileName + ` ` + params
	cmd := exec.Command("cmd", command)
	res, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, "Server error. Can't exec: "+err.Error()+"\n "+string(res), 500)
		return
	}
	http.ServeFile(w, r, wavFileName)
	err = os.Remove(txtFileName)
	if err != nil {
		logger.Error("Can't remove txt file")
	}
	err = os.Remove(wavFileName)
	if err != nil {
		logger.Error("Can't remove wav file")
	}
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
	http.HandleFunc("/voices", voicesHandler)
	http.HandleFunc("/synth", synthHandler)
	http.ListenAndServe(":"+config.Config.Port, nil)
}
