package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Info struct {
	Title    string `json:"Title"`
	PartName string `json:"PartName"`
	PartNo   string `json:"PartNo"`
}

func main() {
	downloaddir, err := filepath.Abs(".")
	if err != nil {
		return
	}
	filepath.Walk(downloaddir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".info" {
			videos, ext := []string{}, ""
			dir := filepath.Dir(path)
			flvfiles, err := filepath.Glob(filepath.Join(dir, "*.flv"))
			if len(flvfiles) != 0 && err == nil {
				videos, ext = flvfiles, ".flv"
			}
			mp4files, err := filepath.Glob(filepath.Join(dir, "*.mp4"))
			if len(mp4files) != 0 && err == nil {
				videos, ext = mp4files, ".mp4"
			}
			if len(videos) == 0 || ext == "" {
				return nil
			}
			infofile, err := ioutil.ReadFile(path)
			if err != nil {
				return nil
			}
			info := Info{}
			err = json.Unmarshal(infofile, &info)
			if err != nil {
				return nil
			}
			if info.PartName == "" {
				info.PartName = info.Title
			}
			targetdir := filepath.Join(downloaddir, info.Title)
			_, err = os.Stat(targetdir)
			if os.IsNotExist(err) {
				fmt.Println("mkdir", targetdir)
				os.Mkdir(targetdir, os.ModeDir)
			}
			oldpath := videos[0]
			newpath := filepath.Join(targetdir, fmt.Sprintf("%s-%s.%s", info.PartNo, info.PartName, ext))
			fmt.Println(oldpath, "--->", newpath)
			os.Rename(oldpath, newpath)
		}
		return nil
	})
}
