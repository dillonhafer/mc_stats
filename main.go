package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/dillonhafer/mc_stats/gziphandler"
	"github.com/dillonhafer/mc_stats/nbt"
)

func searchForWorlds() (string, error) {
	var path string
	var world string

	switch runtime.GOOS {
	case "darwin":
		path = os.Getenv("HOME") + "/Library/Application Support/minecraft/saves"
	default:
		return "", errors.New("Unknown OS")
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	number := 1
	var worlds []string
	for _, f := range files {
		if f.IsDir() {
			println(fmt.Sprintf("  %d. %s", number, f.Name()))
			worlds = append(worlds, f.Name())
			number = number + 1
		}
	}

	if len(worlds) <= 0 {
		return "", errors.New("Could not find any worlds")
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please select a world: ")
	w, _ := reader.ReadString('\n')
	s, err := strconv.ParseInt(strings.TrimSpace(w), 10, 32)
	if err != nil {
		return "", errors.New("You did not select a world.")
	}

	world = path + "/" + worlds[int(s)-1]
	fmt.Println("You chose: ", world)

	return world, nil
}
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

func readPlayers(userCache, playerData string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		println(fmt.Sprintf("[%s] GET /players", time.Now().String()))
		if r.Method != "GET" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Command Must be a GET\n")
			return
		}

		w.Header().Set("Content-Type", "application/json")

		users := []nbt.Player{}
		usersCache, err := os.Open(userCache)
		if err != nil {
			return
		}
		jsonParser := json.NewDecoder(usersCache)
		if err = jsonParser.Decode(&users); err != nil {
			return
		}
		for i, _ := range users {
			user := &users[i]
			player := nbt.Load(filepath.Join(playerData, user.UUID+".dat"))
			user.X = player.X
			user.Y = player.Y
			user.Z = player.Z
			user.Xp = player.Xp
		}

		if err != nil {
			fmt.Fprint(w, "[]")
		} else {
			json.NewEncoder(w).Encode(users)
		}
	}
}

func readStats(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		println(fmt.Sprintf("[%s] GET /stats", time.Now().String()))
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
	if os.Getenv("WORLD") != "" {
		world = os.Getenv("WORLD")
	} else {
		var err error
		world, err = searchForWorlds()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not find any local worlds and none were provided.", err)
			return
		}
	}
	staticFiles := http.FileServer(http.Dir("frontend/build"))

	statsPath := filepath.Join(world, "stats")
	userCache := filepath.Join(world, "..", "usercache.json")
	playerData := filepath.Join(world, "playerdata")
	properStatsDir, _ := statsDirExists(statsPath)

	if !properStatsDir {
		fmt.Fprintln(os.Stderr, "You must provide a directory to the world.")
		return
	}

	stats := gziphandler.GzipHandler(readStats(statsPath))
	players := gziphandler.GzipHandler(readPlayers(userCache, playerData))

	http.Handle("/", staticFiles)
	http.Handle("/stats", stats)
	http.Handle("/players", players)

	port := "22334"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	addr := "127.0.0.1"
	if os.Getenv("ADDR") != "" {
		addr = os.Getenv("ADDR")
	}

	cert := ""
	key := ""
	if os.Getenv("TLS_CERT") != "" {
		cert = os.Getenv("TLS_CERT")
	}
	if os.Getenv("TLS_KEY") != "" {
		key = os.Getenv("TLS_KEY")
	}

	fmt.Println("Server running and listening on port " + port)
	fmt.Println("Run `mc_stats -h` for more startup options")
	fmt.Println("Ctrl-C to shutdown server")

	var err error
	if cert != "" && key != "" {
		err = http.ListenAndServeTLS(addr+":"+port, cert, key, nil)
	} else {
		err = http.ListenAndServe(addr+":"+port, nil)
	}
	fmt.Fprintln(os.Stderr, err)
}
