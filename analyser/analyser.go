//Package analyser provides implementations for mobile app package analysis
package analyser

import (
	"archive/zip"
	"errors"
	"log"
	"regexp"
	"strings"
)

const platformAndroid = "android"
const PlatformIOS = "ios"

type IAnalyser interface {
	//getAppInfo returns populated AppInfo
	getAppInfo() *AppInfo
	//PrintAppInfo fetches and prints information about the mobile app package
	PrintAppInfo()
}

//AppInfo is a struct, holding the recognisable information of the mobile app package
type AppInfo struct {
	AppId    string
	Version  string
	Build    string
	IconPath string
}

//GetAnalyser returns the appropriate analyser implementation, based on the contents of the package
func GetAnalyser(rc *zip.ReadCloser, originalFilePath *string) (IAnalyser, error) {
	for _, f := range rc.File {
		platform := detectPlatform(&f.Name)

		if platform == platformAndroid {
			return &apkAnalyser{originalFilePath}, nil
		}

		if platform == PlatformIOS {
			return &ipaAnalyser{f}, nil
		}
	}

	err := errors.New("no platform detected")
	return nil, err
}

//detectPlatform returns one of the platform identifier string constants or an empty string
func detectPlatform(fileName *string) string {
	if strings.EqualFold(*fileName, "AndroidManifest.xml") {
		//fmt.Println("Android app recognized")
		return platformAndroid
	}

	plistPathRegex, err := regexp.Compile(`^Payload/[^/]+/Info\.plist$`)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	if plistPathRegex.MatchString(*fileName) {
		//fmt.Println("iOS app recognized")
		return PlatformIOS
	}

	return ""
}
