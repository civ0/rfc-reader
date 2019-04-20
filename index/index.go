package index

import (
	"encoding/xml"
	"github.com/civ0/rfc-reader/cache"
	"github.com/civ0/rfc-reader/util"
	"github.com/pkg/errors"
	"io/ioutil"
	"path"
	"strings"
)

type RFCIndex struct {
	XMLName    xml.Name   `xml:"rfc-index"`
	RFCEntries []RFCEntry `xml:"rfc-entry"`
}

type RFCEntry struct {
	XMLName           xml.Name `xml:"rfc-entry"`
	DocID             string   `xml:"doc-id"`
	Title             string   `xml:"title"`
	Authors           []Author
	Date              Date
	Format            Format
	Keywords          []string `xml:"keywords>kw"`
	Abstract          string   `xml:"abstract>p"`
	Draft             string   `xml:"draft"`
	CurrentStatus     string   `xml:"current-status"`
	PublicationStatus string   `xml:"publication-status"`
	Stream            string   `xml:"stream"`
	Area              string   `xml:"area"`
	WGAcronym         string   `xml:"wg_acronym"`
	DOI               string   `xml:"doi"`
}

func (e RFCEntry) CanonicalName() string {
	var str strings.Builder
	str.WriteString("rfc")
	num := e.DocID[3:]
	for num[0] == '0' {
		num = num[1:]
	}
	str.WriteString(num)
	return str.String()
}

type Author struct {
	Name  string `xml:"name"`
	Title string `xml:"title"`
}

type Date struct {
	Month string `xml:"month"`
	Year  string `xml:"year"`
}

type Format struct {
	FileFormat string `xml:"file-format"`
	CharCount  string `xml:"char-count"`
	PageCount  string `xml:"page-count"`
}

func ReadIndex() (RFCIndex, error) {
	var index RFCIndex

	indexFilePath := path.Join(cache.GetAppCacheDir(), "rfc-index.xml")
	data, err := ioutil.ReadFile(indexFilePath)
	if err != nil {
		err = errors.Wrap(err, "Failed to read RFC index file")
		util.FatalExit(err)
	}

	err = xml.Unmarshal(data, &index)
	if err != nil {
		err = errors.Wrap(err, "Failed to unmarshal XML index")
		util.FatalExit(err)
	}

	return index, nil
}
