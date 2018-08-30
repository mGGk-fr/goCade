package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

// SYSTEM FUNCTIONS
func execUnixSysCommand(c string) []byte {
	out, err := exec.Command("sh", "-c", c).Output()
	if err != nil {
		log.Panic(err)
		return nil
	}
	return out
}

func execWinSysCommand(c string) []byte {
	return nil
}

// HTTP FUNCTIONS
func httpRAMUsage(w http.ResponseWriter, r *http.Request) {
	var co []byte
	switch runtime.GOOS {
	case "linux", "darwin":
		co = execUnixSysCommand(" awk '/MemAvailable/ { printf \"%.3f\", $2/1024/1024 }' /proc/meminfo")
	case "windows":
		co = execWinSysCommand("")
	}
	w.Write(co)
}

func httpRAMTotal(w http.ResponseWriter, r *http.Request) {
	var co []byte
	switch runtime.GOOS {
	case "linux", "darwin":
		co = execUnixSysCommand(" awk '/MemTotal/ { printf \"%.3f\", $2/1024/1024 }' /proc/meminfo")
	case "windows":
		co = execWinSysCommand("")
	}
	w.Write(co)
}

func httpCPUUsage(w http.ResponseWriter, r *http.Request) {
	var co []byte
	switch runtime.GOOS {
	case "linux", "darwin":
		co = execUnixSysCommand("grep 'cpu ' /proc/stat | awk '{usage=($2+$4)*100/($2+$4+$5)} END {print usage}'")
	case "windows":
		co = execWinSysCommand("")
	}
	w.Write(co)
}

func httpCurrentGame(w http.ResponseWriter, r *http.Request) {
	var co []byte
	switch runtime.GOOS {
	case "linux", "darwin":
		co = execUnixSysCommand("ps -ef | grep '[/]usr/bin/mame' | awk '{printf $NF}'")
	case "windows":
		co = execWinSysCommand("")
	}
	var gn = retriveGameName(string(co[:]))
	w.Write([]byte(gn))
}

// UTILS FUNCTIONS
func retriveGameName(g string) string {
	var co []byte
	if g == "" {
		return "Pas de jeu en cours"
	}
	switch runtime.GOOS {
	case "linux", "darwin":
		co = execUnixSysCommand("mame -listfull " + g + " | sed -n 2p | cut -d'\"' -f2")
		return string(co[:])
	case "windows":
		return ""
	}
	return "ERREUR"
}

func main() {
	fmt.Println("goCade v0.1 running on " + runtime.GOOS)
	fmt.Println("Starting web server and defining routes")
	http.Handle("/", http.FileServer(http.Dir("./pages")))
	http.HandleFunc("/getRAMUsage", httpRAMUsage)
	http.HandleFunc("/getRAMTotal", httpRAMTotal)
	http.HandleFunc("/getCPUUsage", httpCPUUsage)
	http.HandleFunc("/getCurrentGame", httpCurrentGame)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
