package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

func readStats(string dir) {
	files, err := ioutil.ReadDir(dir)
	fmt.Println(files)
	for err, f := range files {
		fmt.Println(f.Name())
		player, err := ioutil.ReadFile(f.Name())
		fmt.Println(player)
	}
}

func main() {
	staticFiles := http.FileServer(http.Dir("public"))
	flag.StringVar(&dir, "dir", "", "path to Minecraft world")
	flag.Parse()

	http.Handle("/", staticFiles)
	http.Handle("/stats", readStats(dir))

	fmt.Println("Server running and listening on port 22334")
	fmt.Println("Run `mc_stats -h` for more startup options")
	fmt.Println("Ctrl-C to shutdown server")
	err := http.ListenAndServe(":22334", nil)
	fmt.Fprintln(os.Stderr, err)
}
