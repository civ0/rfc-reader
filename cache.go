package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

const XMLIndexURL = "https://www.rfc-editor.org/in-notes/rfc-index.xml"
const RFCBaseURL = "https://www.rfc-editor.org/rfc/"

func GetAppCacheDir() string {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatalln("os.UserCacheDir() failed: %v", err)
	}
	appCacheDir := path.Join(userCacheDir, "rfc-reader")
	return appCacheDir
}

func UpdateCache() error {
	log.Println("Setting up cache dir...")
	appCacheDir := GetAppCacheDir()
	// also create "rfc" sub dir for caching RFCs
	err := os.MkdirAll(path.Join(appCacheDir, "rfc"), 0755)
	if err != nil {
		log.Fatalln("Failed to create cache dir: %v", err)
		return err
	}

	indexFilePath := path.Join(appCacheDir, "rfc-index.xml")
	if _, err = os.Stat(indexFilePath); err == nil {
		log.Println("RFC index is present")
		// File exists
		// TODO: Update based on timestamp
	} else if os.IsNotExist(err) {
		log.Println("RFC index not present")
		err = DownloadIndex(indexFilePath)
	} else {
		// TODO: Proper error handling
		// https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return err
	}
	log.Println("Cache dir setup done")
	return err
}

func DownloadIndex(path string) error {
	log.Println("Downloading RFC index")
	res, err := http.Get(XMLIndexURL)
	if err != nil {
		return err
	}

	index, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	res.Body.Close()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(index)

	return err
}

func ReadRFC(canonicalName string) (string, error) {
	rfcFilePath := path.Join(GetAppCacheDir(), "rfc", canonicalName+".txt")
	if _, err := os.Stat(rfcFilePath); err == nil {
		data, err := ioutil.ReadFile(rfcFilePath)
		if err != nil {
			log.Fatalln("Failed to read RFC index file: %v", err)
			return "", err
		}
		return string(data), nil
	} else if os.IsNotExist(err) {
		err = DownloadRFC(canonicalName)
		if err != nil {
			log.Fatalln("Failed to dowload RFC %s: %v", canonicalName, err)
			return "", err
		}
	} else {
		// TODO: Proper error handling
		// https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return "", err
	}
	// TODO: Cleanup returns
	return "", nil
}

func DownloadRFC(canonicalName string) error {
	log.Println("Downloading RFC %s", canonicalName)
	res, err := http.Get(RFCBaseURL + canonicalName + ".txt")
	if err != nil {
		return err
	}

	rfc, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	res.Body.Close()

	log.Println("INFO: Saving RFC %s", canonicalName)
	f, err := os.Create(path.Join(GetAppCacheDir(), "rfc", canonicalName+".txt"))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(rfc)

	return err
}
