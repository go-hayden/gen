package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"path"
	"regexp"

	cl "github.com/fatih/color"
	tb "github.com/go-hayden/toolbox"
)

var importReg *regexp.Regexp
var implReg *regexp.Regexp
var returnUrlReg *regexp.Regexp
var replaceHash *regexp.Regexp

func init() {
	importReg = regexp.MustCompile(`^#import\s+"\S+"`)
	implReg = regexp.MustCompile(`^@implementation\s+mparam_\S+`)
	returnUrlReg = regexp.MustCompile(`^\s*return\s+@"https{0,1}:\S+"`)
	replaceHash = regexp.MustCompile(`return|\s|"|@|;`)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-h" {
		printHelp()
		return
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		cl.Red("错误：由于无法获得当前路径，未能生成NVNetwork+mapi.h和NVModels.h文件！")
		return
	}

	config := readConfig(dir)
	if config == nil {
		cl.Red("读取当前路径下gen_config.json文件失败, Gen程序退出")
		return
	}
	config.CurrentPath = dir
	config.print()
	if !config.AddMparamInit && !config.GenMapiHeader && !config.GenModelHeader && config.URLMaps == nil {
		cl.Green("根据配置，Gen程序没有做任何事！")
		return
	}

	c := make(chan error, 3)
	go generateMapiHeader(config, c)
	go generateModelHeader(config, c)
	go generateMapiMFile(config, c)
	for index := 0; index < 3; index++ {
		err := <-c
		if err != nil {
			cl.Red("错误：gen执行失败，原因->" + err.Error())
			return
		}
	}
	if config.GenMapiHeader {
		cl.Green("生成文件：" + path.Join(dir, "NVNetworkClient+mapi.h"))
	}
	if config.GenModelHeader {
		cl.Green("生成文件：" + path.Join(dir, "NVModels.h"))
	}
	if config.AddMparamInit || config.URLMaps != nil {
		cl.Green("重新生成文件：" + path.Join(dir, "NVNetwork", "NVNetworkClient+mapi.m"))
	}
}

const helpString = `
Gen由golang编写，具体使用文档参考：http://wiki.sankuai.com/pages/viewpage.action?pageId=729312051

`

func printHelp() {
	print(helpString)
}

func readConfig(currentPath string) *GenConfig {
	configPath := path.Join(currentPath, "gen_config.json")
	if !tb.FileExists(configPath) {
		cl.Red("警告：映射文件不存在[" + configPath + "]")
		return nil
	}

	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		cl.Red("警告：无法读取配置文件[" + configPath + "]")
		return nil
	}

	var config *GenConfig
	if err := json.Unmarshal(bytes, &config); err != nil {
		cl.Red("警告：无法解析配置文件，请确认配置文件格式正确[" + configPath + "]")
		return nil
	}
	return config
}
