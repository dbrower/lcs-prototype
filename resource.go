package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Resource struct {
	ID       string      `xml:"identifier"`
	Level    AccessLevel `xml:-`
	XMLName  xml.Name
	Metadata []Pair `xml:",any"`
	Files    []File
}

type Pair struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type File struct {
	ID     string
	MD5    []byte
	SHA256 []byte
	Size   int64
}

type AccessLevel int

const (
	AccessPrivate AccessLevel = iota
	AccessLogin
	AccessPublic
)

func LoadXML(resourcePath string) (*Resource, error) {
	var result Resource

	mdContents, err := ioutil.ReadFile(path.Join(resourcePath, "meta.xml"))
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(mdContents, &result)
	if err != nil {
		return nil, err
	}

	switch len(result.Get("accesslevel")) {
	case 0:
		result.Level = AccessPrivate
	case 1:
		result.Level = parseAccessLevel(result.Get("accesslevel")[0])
	default:
		return nil, errors.New("Too many access levels")
	}

	err = (&result).scanFiles(resourcePath)

	return &result, err
}

func parseAccessLevel(s string) AccessLevel {
	switch strings.ToLower(s) {
	case "public":
		return AccessPublic
	case "login":
		return AccessLogin
	default:
		return AccessPrivate
	}
}

func (res *Resource) scanFiles(resourcePath string) error {

	filepath.Walk(resourcePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// we don't scan directory entries themselves
			return nil
		}
		ID := strings.TrimPrefix(path, resourcePath)
		if ID == "FILES.xml" || ID == "LOG.txt" {
			return nil
		}

		found := false
		for _, f := range res.Files {
			if f.ID == ID {
				// file was already scanned
				found = true
			}
		}
		if !found {
			// new file
			fmt.Println("File added:", ID)
			res.Files = append(res.Files, File{ID: ID, Size: info.Size()})
		}
		return nil
	})

	return nil
}

func (res Resource) Get(fieldname string) []string {
	var result []string
	for _, pair := range res.Metadata {
		if pair.XMLName.Local == fieldname {
			result = append(result, pair.Value)
		}
	}
	return result
}
