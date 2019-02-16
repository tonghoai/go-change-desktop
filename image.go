package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/jasonlvhit/gocron"
)

/*
  API of unsplash.com
  get random image from server
*/
const accessKey = "bec1e68b3836c8bac251b043c243a2213fb95baaf77325ce20d38bf3faa89da9"
const endpoint = "https://api.unsplash.com/photos/random/"
const url = endpoint + "?client_id=" + accessKey
const PrefixPictureFolder = "pictures/"

var types = []string{"none", "wallpaper", "centered", "scaled", "stretched", "zoom", "spanned"}

// Response of get image api
type Response struct {
	ID   string `json:"id"`
	Urls struct {
		Raw string `json:"raw"`
	} `json:"urls"`
}

func main() {
	/*
	  Flag command
	  * s: seconds to change desktop background
	  * t: type to set desktop background
	*/
	secondsFlag := flag.Uint64("s", 3600, "time to change image")
	typeFlag := flag.String("t", "centered", "size")
	flag.Parse()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("error at main: ", r)
		}
	}()

	// Validate t command
	validType := false
	for _, t := range types {
		if t == *typeFlag {
			validType = true
			break
		}
	}
	if !validType {
		panic(errors.New("type is not validate"))
	}

	// To run task change desktop background
	job := func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("error at job: ", r)
			}
		}()

		var imageName string
		var err error
		image, err := GetImage()
		if err != nil {
			panic(err)
		}
		err = DownloadImage(image.ID+".jpg", image.Urls.Raw)
		if err != nil {
			panic(err)
		}
		imageName = image.ID + ".jpg"
		err = setBackgroundImage(imageName, *typeFlag)
		if err != nil {
			panic(err)
		}
	}

	gocron.Every(*secondsFlag).Seconds().Do(job)
	<-gocron.Start()
}

func setBackgroundImage(image string, typeFlag string) error {
	currentPath, err := GetCurrentPath()
	if err != nil {
		return err
	}
	if len(image) > 0 {
		_, err := exec.Command("/bin/bash", "-c", "gsettings set org.gnome.desktop.background picture-uri \"file://"+currentPath+PrefixPictureFolder+image+"\"").Output()
		if err != nil {
			return err
		}
		_, err = exec.Command("/bin/bash", "-c", "gsettings set org.gnome.desktop.background picture-options \""+typeFlag+"\"").Output()
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("no image to set")
}

func GetCurrentPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return dir + "/", nil
}

func GetImage() (Response, error) {
	response := new(Response)
	res, err := http.Get(url)
	if err != nil {
		return *response, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return *response, err
	}

	return *response, nil
}

func DownloadImage(filepath string, url string) error {
	out, err := os.Create(PrefixPictureFolder + filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
