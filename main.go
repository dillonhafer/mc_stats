package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func statsDirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func readStats(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Command Must be a GET\n")
			return
		}
		w.Header().Set("Content-Type", "application/json")

		files, _ := ioutil.ReadDir(dir)
		players := make([]string, 0, len(files))

		for _, file := range files {
			path := filepath.Join(dir, file.Name())
			f, _ := ioutil.ReadFile(path)
			players = append(players, string(f))
		}

		player_json := strings.Join(players, ",")
		fmt.Fprint(w, player_json)
	}
}

func main() {
	var dir string
	staticFiles := http.FileServer(http.Dir("public"))
	flag.StringVar(&dir, "dir", "", "path to Minecraft world")
	flag.Parse()

	properStatsDir, _ := statsDirExists(dir)
	if !properStatsDir {
		fmt.Fprintln(os.Stderr, "You must provide a directory the world.")
		fmt.Println("Run `mc_stats -h` for more startup options")
		return
	}

	http.Handle("/", staticFiles)
	http.HandleFunc("/stats", readStats(dir))

	fmt.Println("Server running and listening on port 22334")
	fmt.Println("Run `mc_stats -h` for more startup options")
	fmt.Println("Ctrl-C to shutdown server")
	err := http.ListenAndServe(":22334", nil)
	fmt.Fprintln(os.Stderr, err)
}
