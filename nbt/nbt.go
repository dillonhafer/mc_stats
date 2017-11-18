// NBT pretty printer.
package nbt

import (
	"compress/gzip"
	"flag"
	"io"
	"log"
	"os"

	"github.com/minero/minero/proto/nbt"
)

type Player struct {
	Name string  `json:"name"`
	UUID string  `json:"uuid"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Z    float64 `json:"z"`
	Xp   int32   `json:"level"`
}

type Players []Player

func ParsePlayer(f io.Reader) Player {
	var r io.Reader
	r = f

	c, err := nbt.Read(r)
	if err != nil {
		log.Fatalln("nbt.Read:", err)
	}

	player := Player{}
	d := c.Value["Pos"].(*nbt.List)
	player.Xp = int32(c.Value["XpLevel"].(*nbt.Int32).Int32)
	for i, f := range d.Value {
		val := f.(*nbt.Float64).Float64
		switch i {
		case 0:
			player.X = float64(val)
		case 1:
			player.Y = float64(val)
		case 2:
			player.Z = float64(val)
		}
	}
	return player
}

func Load(file string) Player {
	var f io.ReadCloser
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Couldn't open file: %q.\n", flag.Arg(0))
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Close()
	return ParsePlayer(r)
}
