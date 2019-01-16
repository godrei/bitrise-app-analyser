package analyser

import (
	"testing"
)

var androidManifestPath = "AndroidManifest.xml"
var invalidAndroidManifestPath = "InvalidAndroidManifest.xml"
var plistFilePath = "Payload/Appname.app/Info.plist"
var invalidPlistFilePath = "InvalidPayload/Appname.app/Info.plist"
var empty = ""

//Unit test for the detectPlatform function

func TestGetAndroidPlatform(t *testing.T) {
	platform := detectPlatform(&androidManifestPath)
	if platform != platformAndroid {
		t.Error("Platform should be " + platformAndroid)
	}
}

func TestNotGetAndroidPlatform(t *testing.T) {
	platform := detectPlatform(&invalidAndroidManifestPath)
	if platform == platformAndroid {
		t.Error("Platform should NOT be " + platformAndroid)
	}

	platform = detectPlatform(&plistFilePath)
	if platform == platformAndroid {
		t.Error("Platform should NOT be " + platformAndroid)
	}
}

func TestGetIOSPlatform(t *testing.T) {
	platform := detectPlatform(&plistFilePath)
	if platform != PlatformIOS {
		t.Error("Platform should be " + PlatformIOS)
	}
}

func TestNotGetIOSPlatform(t *testing.T) {
	platform := detectPlatform(&invalidPlistFilePath)
	if platform == PlatformIOS {
		t.Error("Platform should NOT be " + PlatformIOS)
	}

	platform = detectPlatform(&androidManifestPath)
	if platform == PlatformIOS {
		t.Error("Platform should NOT be " + PlatformIOS)
	}
}

func TestWithEmptyFilePath(t *testing.T) {
	platform := detectPlatform(&empty)
	if platform == platformAndroid {
		t.Error("Platform should NOT be " + platformAndroid)
	}

	if platform == PlatformIOS {
		t.Error("Platform should NOT be " + PlatformIOS)
	}
}
