package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
)

func httpRAMUsage(w http.ResponseWriter, r *http.Request) {
	out, err := exec.Command("sh", "-c", " awk '/MemAvailable/ { printf \"%.3f\", $2/1024/1024 }' /proc/meminfo").Output()
	if err != nil {
		panic(err)
	}
	w.Write(out)
}

func httpRAMTotal(w http.ResponseWriter, r *http.Request) {
	out, err := exec.Command("sh", "-c", " awk '/MemTotal/ { printf \"%.3f\", $2/1024/1024 }' /proc/meminfo").Output()
	if err != nil {
		panic(err)
	}
	w.Write(out)
}

func httpCPUUsage(w http.ResponseWriter, r *http.Request) {
	out, err := exec.Command("sh", "-c", "grep 'cpu ' /proc/stat | awk '{usage=($2+$4)*100/($2+$4+$5)} END {print usage}'").Output()
	if err != nil {
		panic(err)
	}
	w.Write(out)
}

func main() {
	fmt.Println("goCade v0.1 running on " + runtime.GOOS)
	fmt.Println("Starting web server and defining routes")
	http.Handle("/", http.FileServer(http.Dir("./pages")))
	http.HandleFunc("/getRAMUsage", httpRAMUsage)
	http.HandleFunc("/getRAMTotal", httpRAMTotal)
	http.HandleFunc("/getCPUUsage", httpCPUUsage)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
