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

	inFlag := flag.String("dir", "./", "dir is to specify what folder to parse types from")
	fileFlag := flag.String("file", "-", "file is to parse a single file. Will override a directory")
	outFlag := flag.String("out", "./", "dir is to specify what folder to parse types to")
	recursiveFlag := flag.Bool("r", true, "to recursively ascend all folders in dir")
	flag.Usage = usage
	flag.Parse()

	// Try to be smart about where to save
	var out string
	if strings.HasSuffix(*outFlag, ".js") {
		out = *outFlag
	} else if strings.HasSuffix(*outFlag, "/") {
		out = *outFlag + "models.js"
	} else {
		out = *outFlag + "/models.js"
	}

	fi, err := os.Create(out)
	if err != nil {
		log.WithError(err).Fatalln("error creating file")
	}
	defer fi.Close()

	p := parse.New(*recursiveFlag, fi)
	spin.Start()

	if *fileFlag != "-" {
		if !strings.HasSuffix(*fileFlag, ".go") {
			log.Error("the file passed in is not a go file.")
			os.Exit(1)
		}
		p.Files = append(p.Files, *inFlag)
		if err := p.ParseFiles(); err != nil {
			log.WithError(err).Fatalln("error parsing file")
		}
	} else {
		if err := p.ParseDir(*inFlag); err != nil {
			log.WithError(err).Fatalln("error parsing directory")
		}
		if err := p.ParseFiles(); err != nil {
			log.WithError(err).Fatalln("error parsing directory")
		}
	}

	p.WriteDocument()

	spin.Stop()
	log.WithField("save_location", out).Info("saved")
	log.WithField("duration", time.Now().Sub(start)).Info("completed code generation")
}
