package util

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
)

func FatalExit(err error) {
	fmt.Printf("%+v\n", err)
	os.Exit(1)
}

func FileExists(path string) (bool, error) {
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil // Don't return err, as it is identified as IsNotExist
	} else {
		// TODO: Possible to properly check this?
		// File may or may not exist, but difficult to check
		return false, errors.Wrapf(err, "Failed to check if file %s exists", path)
	}
}

func ReadFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read file %s", path)
	}
	return data, nil
}

func DownloadFile(url, path string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		} else {
			err = nil
		}
	}()

	response, err := http.Get(url)
	if err != nil {
		panic(errors.Wrapf(err, "Failed to download file from %s", url))
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(errors.Wrap(err, "Failed to reade response body"))
	}
	err = response.Body.Close()
	if err != nil {
		panic(errors.Wrap(err, "Failed to close response body"))
	}

	f, err := os.Create(path)
	if err != nil {
		panic(errors.Wrapf(err, "Failed to create file %s", path))
	}

	defer func() {
		ferr := f.Close()
		if ferr != nil {
			panic(errors.Wrap(ferr, "Failed to close file"))
		}
	}()

	_, err = f.Write(data)
	if err != nil {
		panic(errors.Wrapf(err, "Failed to write data to file %s", err))
	}

	return
}
