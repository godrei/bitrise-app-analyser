package main

import (
	"archive/zip"
	"bitrise-app-analyser/analyser"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		inputFilePath := os.Args[1]
		handleInputFile(&inputFilePath)
	} else {
		fmt.Println("Please specify an APK/IPA file path as input argument")
	}
}

//handleInputFile unzips input file and prints detected app info
func handleInputFile(path *string) {
	rc, err := zip.OpenReader(*path)
	if err != nil {
		fmt.Println("Unable to unzip file")
		return
	}
	defer rc.Close()

	appAnalyser, err := analyser.GetAnalyser(rc, path)
	if err != nil {
		fmt.Println(err)
		return
	}

	appAnalyser.PrintAppInfo()
}
