package main

import (
	"fmt"
	"net/http"

	"github.com/capnm/sysinfo"
)

func httpRAMUsage(w http.ResponseWriter, r *http.Request) {
	si := sysinfo.Get()
	w.Write([]byte(fmt.Sprintf("%d", si.TotalRam)))
}

func main() {
	fmt.Println("goCade v0.1")
	fmt.Println("Starting web server and defining routes")
	http.Handle("/", http.FileServer(http.Dir("./pages")))
	http.HandleFunc("/getRAMUsage", httpRAMUsage)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
