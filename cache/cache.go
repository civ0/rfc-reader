package cache

import (
	"fmt"
	"github.com/civ0/rfc-reader/util"
	"github.com/pkg/errors"
	"os"
	"path"
)

const XMLIndexURL = "https://www.rfc-editor.org/in-notes/rfc-index.xml"
const RFCBaseURL = "https://www.rfc-editor.org/rfc/"

type RFCLocalCopyStatus int

const (
	RFCLocalCopyStatusAbsent RFCLocalCopyStatus = iota
	RFCLocalCopyStatusPresent
	RFCLocalCopyStatusUnknown
)

func GetAppCacheDir() string {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		err = errors.Wrap(err, "Failed to get UserCacheDir")
		util.FatalExit(err)
	}
	appCacheDir := path.Join(userCacheDir, "rfc-reader")
	return appCacheDir
}

func UpdateCache() error {
	appCacheDir := GetAppCacheDir()
	// also create "rfc" sub dir for caching RFCs
	err := os.MkdirAll(path.Join(appCacheDir, "rfc"), 0755)
	if err != nil {
		err = errors.Wrapf(err, "Failed to create app cache dir %s", appCacheDir)
		util.FatalExit(err)
	}

	indexFilePath := path.Join(appCacheDir, "rfc-index.xml")
	exists, err := util.FileExists(indexFilePath)
	if err != nil {
		err = errors.Wrapf(err, "Failed to check if index file %s is present", indexFilePath)
		util.FatalExit(err)
	}

	if exists == false {
		err = DownloadIndex(indexFilePath)
		if err != nil {
			err = errors.Wrap(err, "Failed to download RFC index file")
		}
	}

	return err
}

func DownloadIndex(path string) error {
	err := util.DownloadFile(XMLIndexURL, path)
	if err != nil {
		err = errors.Wrapf(err, "Failed to Download RFC index file from %s to %s", XMLIndexURL, path)
	}
	return err
}

func GetRFCFilePath(canonicalName string) string {
	rfcFilePath := path.Join(GetAppCacheDir(), "rfc", canonicalName+".txt")
	return rfcFilePath
}

func RFCPresent(canonicalName string) (RFCLocalCopyStatus, error) {
	res := RFCLocalCopyStatusUnknown

	rfcFilePath := GetRFCFilePath(canonicalName)
	exists, err := util.FileExists(rfcFilePath)
	if err != nil {
		err = errors.Wrapf(err, "Failed to check if RFC %s is present", canonicalName)
	}

	if exists == true {
		res = RFCLocalCopyStatusPresent
	} else {
		res = RFCLocalCopyStatusAbsent
	}

	return res, err
}

func ReadRFC(canonicalName string) (rfc string, err error) {
	defer func() {
		if r := recover(); r != nil {
			rfc = fmt.Sprintf("%+v\n", err) // Return error message as RFC content to display to user
			err = r.(error)
		}
	}()

	localCopyStatus, err := RFCPresent(canonicalName)
	if err != nil {
		panic(errors.WithStack(err))
	}

	if localCopyStatus != RFCLocalCopyStatusPresent {
		err = DownloadRFC(canonicalName)
		if err != nil {
			panic(errors.WithStack(err))
		}
	}

	rfcFilePath := GetRFCFilePath(canonicalName)
	rfcData, err := util.ReadFile(rfcFilePath)
	if err != nil {
		panic(errors.WithStack(err))
	}
	rfc = string(rfcData)

	return
}

func DownloadRFC(canonicalName string) error {
	url := RFCBaseURL + canonicalName + ".txt"
	rfcFilePath := GetRFCFilePath(canonicalName)
	err := util.DownloadFile(url, rfcFilePath)
	if err != nil {
		err = errors.Wrapf(err, "Failed to download RFC %s from %s to %s", canonicalName, url, rfcFilePath)
	}
	return err
}
