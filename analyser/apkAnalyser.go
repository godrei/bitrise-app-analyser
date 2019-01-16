package analyser

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

//apkAnalyser is an IAnalyser implementation for Android platform
type apkAnalyser struct {
	filePath *string
}

//PrintAppInfo fetches and prints information about the Android app package
func (analyser apkAnalyser) PrintAppInfo() {
	appInfo := analyser.getAppInfo()
	fmt.Printf("Recognised app info:\n "+
		"\tPackage Name: %+v\n "+
		"\tVersion name: %+v\n "+
		"\tVersion code: %+v\n "+
		"\tApp icon file path: %+v\n",
		appInfo.AppId,
		appInfo.Version,
		appInfo.Build,
		appInfo.IconPath)
}

//getAppInfo returns populated AppInfo
func (analyser apkAnalyser) getAppInfo() *AppInfo {
	buildToolsVersion := getBuildToolsVersion()
	buildToolPath := "/build-tools/" + buildToolsVersion
	androidSDKPath := os.ExpandEnv("$ANDROID_HOME")
	if len(androidSDKPath) == 0 {
		log.Fatal("$ANDROID_HOME environment variable not set")
		return nil
	}
	cmd := androidSDKPath + buildToolPath + "/aapt"
	output, err := exec.Command(cmd, "dump", "badging", *analyser.filePath).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	//Instead of regexes we could print some properties with apkanalyzer, except the app icons.
	//Also apkanalyzer is not able to fetch version code from a release apk
	appIdRegex, err := regexp.Compile(`package: name='(?P<package>[^']+)'`)
	versionCodeRegex, err := regexp.Compile(`versionCode='(?P<code>[^']+)'`)
	versionNameRegex, err := regexp.Compile(`versionName='(?P<name>[^']+)'`)
	appIconRegex, err := regexp.Compile(`icon='(?P<name>[^']+)'`)

	if err != nil {
		log.Fatal(err)
	}

	packageName := string(appIdRegex.FindSubmatch(output)[1])
	versionCode := string(versionCodeRegex.FindSubmatch(output)[1])
	versionName := string(versionNameRegex.FindSubmatch(output)[1])
	appIcon := string(appIconRegex.FindSubmatch(output)[1]) //TODO find app icon for all densities

	return &AppInfo{packageName, versionName, versionCode, appIcon}
}

//getBuildToolsVersion returns the first available build tools version
//available in the current Android SDK installation
func getBuildToolsVersion() string {
	cmd := "ls"
	param := os.ExpandEnv("$ANDROID_HOME") + "/build-tools"
	output, err := exec.Command(cmd, param).CombinedOutput()

	if err != nil {
		log.Fatal(err)
	}

	return strings.Trim(string(output), "\n")
}
