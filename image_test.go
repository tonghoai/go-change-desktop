package main

import (
	"os"
	"testing"
)

func TestGetImage(t *testing.T) {
	res, err := GetImage()
	if err != nil {
		t.Error("cannot get api ", err)
	}

	if res == (Response{}) {
		t.Error("response is empty")
	}
}

func TestDownloadImage(t *testing.T) {
	res, err := GetImage()
	if err != nil {
		t.Error("cannot get api ", err)
	}

	if res == (Response{}) {
		t.Error("response is empty")
	}

	err = DownloadImage(res.ID+".jpg", res.Urls.Raw)
	if err != nil {
		t.Error("cannot download image", err)
	}

	currentPath, err := GetCurrentPath()
	if err != nil {
		t.Error("cannot get current path", err)
	}

	if _, err := os.Stat(currentPath + PrefixPictureFolder + res.ID + ".jpg"); os.IsNotExist(err) {
		t.Error("download complete but file do not exist", err)
	}
}
