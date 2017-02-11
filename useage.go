package main

import "fmt"

func usage() {
	fmt.Print(`
	GoFlow Usage:
	Flags:
		-dir	Parse a complete directory 
			example: 	-dir= ../src/appname/models/
			default: 	"./"

		-file	Parse a single go file 
			example: 	-file= ../src/appname/models/app.go
			overrides 	-dir and -recursive

		-out	Saves content to folder
			example: 	-out= ../src/appname/models/
						-out= ../src/appname/models/customname.js
			default: 	"./models". 
		-r	Transcends directories
			example:	-recursive= false
			default:	"true"
`)
}
