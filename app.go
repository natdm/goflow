package main

import (
	"flag"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/briandowns/spinner"
	"github.com/natdm/goflow/parse"
)

func main() {
	start := time.Now()

	spin := spinner.New(spinner.CharSets[35], time.Second/3)
	spin.Color("green")

	in := flag.String("dir", "./", "dir is to specify what folder to parse types from")
	file := flag.String("file", "-", "file is to parse a single file. Will override a directory")
	out := flag.String("out", "./", "dir is to specify what folder to parse types to")
	recursive := flag.Bool("r", true, "to recursively ascend all folders in dir")
	flag.Usage = usage
	flag.Parse()

	// Try to be smart about where to save
	var outFile string
	if strings.HasSuffix(*out, ".js") {
		outFile = *out
	} else if strings.HasSuffix(*out, "/") {
		outFile = *out + "models.js"
	} else {
		outFile = *out + "/models.js"
	}

	fi, err := os.Create(outFile)
	if err != nil {
		log.WithError(err).Fatalln("error creating file")
	}
	defer fi.Close()

	p := parse.New(*recursive, fi)
	spin.Start()

	if *file != "-" {
		if !strings.HasSuffix(*file, ".go") {
			log.Error("the file passed in is not a go file.")
			os.Exit(1)
		}
		p.Files = append(p.Files, *file)
		if err := p.ParseFiles(); err != nil {
			log.WithError(err).Fatalln("error parsing file")
		}
	} else {
		if err := p.ParseDir(*in); err != nil {
			log.WithError(err).Fatalln("error parsing directory")
		}
		if err := p.ParseFiles(); err != nil {
			log.WithError(err).Fatalln("error parsing directory")
		}
	}

	p.WriteDocument()

	spin.Stop()
	log.WithField("save_location", outFile).Info("saved")
	log.WithField("duration", time.Now().Sub(start)).Info("completed code generation")
}
