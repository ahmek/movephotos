package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var picTypes = map[string]bool{
	"jpg":  true,
	"jpeg": true,
	"png":  true,
	"gif":  true,
	"bmp":  true,
}

var videoTypes = map[string]bool{
	"mp4": true,
	"mov": true,
	"avi": true,
	"wmv": true,
	"mkv": true,
	"rm":  true,
	"f4v": true,
	"flv": true,
	"swf": true,
}

// isVideoFileFix 检测文件后缀是否为视频格式
func isVideoFileFix(fix string) string {
	if _, ok := videoTypes[fix]; ok {
		return "videos"
	}
	return ""
}

// isPicFileFix 检测文件后缀是否为图片格式
func isPicFileFix(fix string) string {
	if _, ok := picTypes[fix]; ok {
		return "pictures"
	}
	return ""
}

func main() {
	var (
		err   error
		total int32
		fi    os.FileInfo
		file  *os.File
	)

	files, _ := getAllFile(".")
	for _, v := range files {
		file, err = os.OpenFile(v, os.O_CREATE|os.O_APPEND, 6)
		if err != nil {
			fmt.Printf("读取文件出错: %s \n", err.Error())
			continue
		}
		if fi, err = file.Stat(); err != nil {
			fmt.Printf("未知错误: %s \n", err.Error())
			continue
		}
		file.Close()

		arr := strings.Split(fi.Name(), ".")
		if len(arr) < 1 {
			continue
		}

		fileFix := strings.ToLower(arr[len(arr)-1])
		ft := isPicFileFix(fileFix)
		if len(ft) == 0 {
			ft = isVideoFileFix(fileFix)
		}
		if len(ft) == 0 {
			fmt.Printf("发现不支持移动的文件类型，程序将不操作它: %s\n", fi.Name())
			continue
		}

		modTime := fi.ModTime()
		dir := fmt.Sprintf("./%s/%d/%02d", ft, modTime.Year(), modTime.Month())
		if _, err = os.Stat(dir); os.IsNotExist(err) {
			if err = os.MkdirAll(dir, 0666); err != nil {
				fmt.Printf("创建目录失败: %s, 原因: %s\n", dir, err.Error())
				continue
			}
			fmt.Println("创建目录: ", strings.TrimLeft(dir, "./"))
		}

		newFilePath := dir + strings.TrimLeft(v, ".")
		if err = os.Rename(v, newFilePath); err != nil {
			fmt.Printf("移动文件失败: %s, 原因: %s\n", v, err.Error())
			continue
		}
		total++
		fmt.Printf("移动文件 %s 到 %s\n", fi.Name(), strings.TrimLeft(newFilePath, "./"))
	}

	exitMsg := "整理完成，即将退出。"
	if total > 0 {
		exitMsg += fmt.Sprintf("总共移动了%d个文件。", total)
	} else {
		exitMsg += "没有移动任何文件。"
	}
	fmt.Println(exitMsg)
	time.Sleep(time.Second * 20)
}

// getAllFile 获取指定目录所有文件
func getAllFile(pathname string) ([]string, error) {
	var s []string
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}

	for _, fi := range rd {
		if !fi.IsDir() {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}
