package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// Info 是每个子视频的元数据
type Info struct {
	Title    string `json:"Title"`
	PartName string `json:"PartName"`
	PartNo   string `json:"PartNo"` // fix video name order problem
}

func main() {
	// 解析哔哩哔哩下载目录绝对路径
	downloaddir, err := filepath.Abs(".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// 如果 info 文件同目录有 flv 文件则整理
	filepath.Walk(downloaddir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".info" {
			dir := filepath.Dir(path)
			videosfile, ext := []string{}, ""
			flvfile, err := filepath.Glob(filepath.Join(dir, "*.flv"))
			if len(flvfile) != 0 && err == nil {
				videosfile = flvfile
				ext = ".flv"
			}
			mp4file, err := filepath.Glob(filepath.Join(dir, "*.mp4"))
			if len(mp4file) != 0 && err == nil {
				videosfile = mp4file
				ext = ".mp4"
			}
			if len(videosfile) == 0 || err != nil || ext == "" {
				return nil
			}
			infofile, err := ioutil.ReadFile(path)
			if err != nil {
				return nil
			}
			infostruct := Info{}
			err = json.Unmarshal(infofile, &infostruct)
			if err != nil {
				return nil
			}
			if infostruct.PartName == "" {
				infostruct.PartName = infostruct.Title
			}
			targetdir := filepath.Join(downloaddir, infostruct.Title)
			_, err = os.Stat(targetdir)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Println("mkdir", targetdir)
					os.Mkdir(targetdir, os.ModeDir)
				}
			}
			oldpath := videosfile[0]
			newpath := filepath.Join(targetdir, fmt.Sprintf("%s-%s.%s", infostruct.PartNo, infostruct.PartName, ext))
			fmt.Println(oldpath, "--->", newpath)
			os.Rename(oldpath, newpath)
		}
		return nil
	})
	fmt.Println("done!")
	time.Sleep(5 * time.Second)
}
