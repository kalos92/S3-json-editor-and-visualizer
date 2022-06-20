package main

import (
	"log"
	"os"
	"s3-file-editor/appconfig"
	"s3-file-editor/server"
	"strconv"

	"github.com/webview/webview"
)

func main() {

	debug := true
	w := webview.New(debug)
	defer w.Destroy()

	config, err := appconfig.ParseConfig()

	if err != nil {
		log.Println("No config found starting on the port 3005")
		config.WebPort = 3005
	}
	srv := server.NewServer(config)
	os.Setenv("JSON_S3_SERVER_PORT", strconv.Itoa(config.WebPort))

	go srv.HandleRequests()

	w.SetTitle("S3 Json Editor")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("http://localhost:" + strconv.Itoa(config.WebPort) + "/build/index.html")
	w.Run()
}

func CheckOSVariables(param string, fnc func(string) (string, bool)) string {
	if p, check := fnc(param); !check {
		panic("The " + param + " is missing")
	} else {
		log.Println(p)
		return p
	}
}
