package main

import (
	"fmt"
	"flag"
	dbus "github.com/guelfey/go.dbus"
)

func main() {

	flag.Parse()
	action := flag.Arg(0)

	if action == "" {
		fmt.Println("Specify an action")
		return
	}

	switch action {
	case "next":
		Next()
	case "prev":
		Previous()
	case "pause":
		PlayPause()
	case "cur":
		CurSong()
	default:
		fmt.Printf("Invalid action %s\n", action)
		return
	}
}

func connDbus() *dbus.Object {
	conn, err := dbus.SessionBus()

	// couldnt connect to session bus
	if err != nil {
		panic(err)
	}

	return conn.Object("org.mpris.MediaPlayer2.spotify", "/org/mpris/MediaPlayer2")
}

func Next() {
	connDbus().Call("Next", 0)
}

func Previous() {
	connDbus().Call("Previous", 0)
}

func PlayPause() {
	connDbus().Call("PlayPause", 0)
}

func CurSong() {
	sdata := new(dbus.Variant)
	// playing status
	pstatus := new(dbus.Variant)
	connDbus().Call("Get", 0, "org.mpris.MediaPlayer2.Player","Metadata").Store(sdata)
	err := connDbus().Call("Get", 0, "org.mpris.MediaPlayer2.Player","PlaybackStatus").Store(pstatus)

	if err != nil {
		// most likely spotify not running
		return
	}

	songData := sdata.Value().(map[string]dbus.Variant)

	title := songData["xesam:title"]
	// buggy spotify dbus only sends a single artist
	artist := songData["xesam:artist"].Value().([]string)
	rating := int(songData["xesam:autoRating"].Value().(float64) * 100)

	if songStatus := pstatus.Value().(string); songStatus == "Paused" {
		fmt.Printf("(paused) %s %s (paused)", artist[0], title)
	} else {
		fmt.Printf("%s %s (%d)", artist[0], title, rating)
	}

}
