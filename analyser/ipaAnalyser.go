package analyser

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
)

//ipaAnalyser is an IAnalyser implementation for iOS platform
type ipaAnalyser struct {
	plistFile *zip.File
}

//PrintAppInfo fetches and prints information about the iOS app package
func (analyser ipaAnalyser) PrintAppInfo() {
	appInfo := analyser.getAppInfo()
	fmt.Printf("Recognised app info:\n "+
		"\tBundle Identifier: %+v\n "+
		"\tVersion number: %+v\n "+
		"\tBuild number: %+v\n "+
		"\tApp icon file path: %+v\n",
		appInfo.AppId,
		appInfo.Version,
		appInfo.Build,
		appInfo.IconPath)
}

//getAppInfo returns populated AppInfo
func (analyser ipaAnalyser) getAppInfo() *AppInfo {
	rc, err := analyser.plistFile.Open()
	plistStr := readFile(rc)

	//Instead of regexes we could print every property from the plist file with PlistBuddy,
	//but to achieve this we need to extract the .ipa to the filesystem
	bundleIdRegex, err := regexp.Compile(`<key>CFBundleIdentifier</key>[^<]+<string>(?P<bundleId>[^<]+)</string`)
	versionRegex, err := regexp.Compile(`<key>CFBundleShortVersionString</key>[^<]+<string>(?P<bundleId>[^<]+)</string`)
	buildRegex, err := regexp.Compile(`<key>CFBundleVersion</key>[^<]+<string>(?P<bundleId>[^<]+)</string`)
	iconRegex, err := regexp.Compile(`<key>CFBundleIconName</key>[^<]+<string>(?P<bundleId>[^<]+)</string`) //TODO find app icon for all densities

	bundleId := bundleIdRegex.FindStringSubmatch(plistStr)[1]
	versionName := versionRegex.FindStringSubmatch(plistStr)[1]
	buildNumber := buildRegex.FindStringSubmatch(plistStr)[1]
	iconPath := iconRegex.FindStringSubmatch(plistStr)[1]

	if err != nil {
		log.Fatal(err)
	}

	return &AppInfo{bundleId, versionName, buildNumber, iconPath}
}

//readFile returns the ReadCloser's content as a string
func readFile(rc io.ReadCloser) string {
	var result string
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(rc)
	result += string(buf.Bytes())
	err = rc.Close()

	if err != nil {
		log.Fatal(err)
	}

	return result
}
