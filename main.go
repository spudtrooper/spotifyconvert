package main

import (
	"flag"
	"fmt"
	"path"
	"regexp"

	"github.com/pkg/errors"
	"github.com/spudtrooper/goutil/check"
	goutillog "github.com/spudtrooper/goutil/log"
	"github.com/spudtrooper/goutil/net"
	"github.com/spudtrooper/spotifydown/api"
)

var (
	url     = flag.String("url", "", "url")
	track   = flag.String("track", "", "track")
	verbose = flag.Bool("verbose", false, "verbose")
	outDir  = flag.String("out_dir", ".", "output directory")

	// https://open.spotify.com/track/50M7nY1oQuNHecs0ahWAtI?si=4e77457217d24e5e
	idRE = regexp.MustCompile(`https://open.spotify.com/track/([^\?]+)`)
)

func realMain() error {
	argsTrack := *track
	var track string
	if argsTrack != "" {
		track = argsTrack
	} else {
		m := idRE.FindStringSubmatch(*url)
		if len(m) != 2 {
			return errors.Errorf("invalid url")
		}
		track = m[1]
	}

	logger := goutillog.MakeLog("spotifyconvert", goutillog.MakeLogColor(true))
	c := api.NewClient(api.NewClientLogger(logger))
	convert, err := c.Convert(api.ConvertTrack(track), api.ConvertVerbose(*verbose))
	if err != nil {
		return err
	}

	uri := convert.Download
	dir := *outDir
	if dir == "" {
		dir = "."
	}
	outFile := path.Join(dir, fmt.Sprintf("%s.mp3", track))

	if *verbose {
		logger.Printf("downloading %s -> %s", uri, outFile)
	}

	if err := net.DownloadFile(outFile, uri); err != nil {
		return err
	}

	if *verbose {
		logger.Printf("downloaded %s -> %s", uri, outFile)
	}

	return nil
}

func main() {
	flag.Parse()
	check.Err(realMain())
}
