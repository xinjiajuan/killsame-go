package main

import (
	"crypto/sha256"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type config []struct {
	Dir Dir_One `yaml:"dir"`
}
type Dir_One struct {
	Enable        bool   `yaml:"enable"`
	DirPath       string `yaml:"dirPath"`
	CheckInterval string `yaml:"checkInterval"`
}

func main() {
	arg := os.Args
	if len(arg) != 2 {
		println("请指定配置文件\n格式:killsame config.yaml")
		return
	}
	config := config{}
	yamlfile, err := ioutil.ReadFile(arg[1])
	if err != nil {
		println(err.Error())
		return
	}
	if err = yaml.Unmarshal(yamlfile, &config); err != nil {
		println(err.Error())
		return
	}
	for _, Adir := range config {
		if Adir.Dir.Enable {
			go RunTicker(Adir.Dir)
		}
	}
	select {}
}
func RunTicker(one Dir_One) {
	duration, er := time.ParseDuration(one.CheckInterval)
	if er != nil {
		fmt.Printf("配置文件检查时间间隔解析错误!%s", er.Error())
		return
	}
	ticker := time.NewTicker(duration)
	DelSameFile(one)
	for {
		<-ticker.C
		DelSameFile(one)
	}
}
func DelSameFile(dirConf Dir_One) {
	fmt.Printf("<<<<<<<%s目录文件检查开始,开始时间%s>>>>>>>\n", dirConf.DirPath, time.Now().Format(time.RFC1123))
	filelist := make(map[string]string)
	files, err := ioutil.ReadDir(dirConf.DirPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			filesha256, er := GetFileSHA256(dirConf.DirPath + file.Name())
			if er != nil {
				println(er.Error())
				return
			}
			_, ok := filelist[filesha256]
			if ok {
				err := os.Remove(dirConf.DirPath + file.Name())
				if err != nil {
					println(dirConf.DirPath+file.Name(), "delete Err")
				} else {
					println(dirConf.DirPath+file.Name(), "is deleted")
				}
			} else {
				filelist[filesha256] = file.Name()
			}
		}
	}
	fmt.Printf(">>>>>>>%s目录文件检查结束,结束时间%s<<<<<<<\n\n", dirConf.DirPath, time.Now().Format(time.RFC1123))
}
func GetFileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum, nil
}
