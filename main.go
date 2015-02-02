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

func readPlayers(userCache string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Command Must be a GET\n")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		player_json, _ := ioutil.ReadFile(userCache)
		fmt.Fprint(w, string(player_json))
	}
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
		AllStats := make([]string, 0, len(files))

		for _, file := range files {
			path := filepath.Join(dir, file.Name())
			filename := file.Name()
			extension := filepath.Ext(filename)
			name := filename[0 : len(filename)-len(extension)]

			f, _ := ioutil.ReadFile(path)
			playerStats := fmt.Sprintf("{\"UUID\": \"%s\", \"data\": %s}", name, f)
			AllStats = append(AllStats, playerStats)
		}

		statsJson := "["
		statsJson += strings.Join(AllStats, ",")
		statsJson += "]"
		fmt.Fprint(w, statsJson)
	}
}

func main() {
	var world string
	staticFiles := http.FileServer(http.Dir("public"))
	flag.StringVar(&world, "world", "", "path to Minecraft world")
	flag.Parse()

	statsPath := filepath.Join(world, "stats")
	userCache := filepath.Join(world, "..", "usercache.json")
	properStatsDir, _ := statsDirExists(statsPath)

	if !properStatsDir {
		fmt.Fprintln(os.Stderr, "You must provide a directory to the world.")
		return
	}

	http.Handle("/", staticFiles)
	http.HandleFunc("/stats", readStats(statsPath))
	http.HandleFunc("/players", readPlayers(userCache))

	fmt.Println("Server running and listening on port 22334")
	fmt.Println("Run `mc_stats -h` for more startup options")
	fmt.Println("Ctrl-C to shutdown server")
	err := http.ListenAndServe(":22334", nil)
	fmt.Fprintln(os.Stderr, err)
}
