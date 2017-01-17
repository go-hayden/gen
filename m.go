package main

import (
	"bytes"
	"io/ioutil"
	"path"

	"errors"

	"strings"

	tb "github.com/go-hayden/toolbox"
)

const initMethodString = `
- (instancetype)init
{
    self = [super init];
    if (self) {
        self.mapi_compress = NO;
        self.mapi_cacheType = NVCacheTypeDisabled;
    }
    return self;
}

`

func generateMapiMFile(cfg *GenConfig, c chan error) {
	if !cfg.AddMparamInit && cfg.URLMaps == nil {
		c <- nil
		return
	}
	mapiMPath := path.Join(cfg.CurrentPath, "NVNetwork", "NVNetworkClient+mapi.m")
	if !tb.FileExists(mapiMPath) {
		c <- errors.New("无法找到NVNetworkClient+mapi.m文件")
		return
	}
	var buffer bytes.Buffer
	tb.ReadLine(mapiMPath, func(line string, finished bool, err error, stop *bool) {
		if cfg.URLMaps != nil && returnUrlReg.MatchString(line) {
			newLine := replaceUrl(line, cfg.URLMaps)
			buffer.WriteString(newLine + "\n")
		} else {
			buffer.WriteString(line + "\n")
			if cfg.AddMparamInit && isImpl(line) {
				buffer.WriteString(initMethodString)
			}
		}
	})

	err := ioutil.WriteFile(mapiMPath, buffer.Bytes(), 0644)
	if err != nil {
		c <- err
		return
	}

	c <- nil
}

func replaceUrl(line string, mapUrl []*URLMap) string {
	returnString := line
	for _, mp := range mapUrl {
		from := strings.TrimSpace(mp.From)
		to := strings.TrimSpace(mp.To)
		if len(from) == 0 || len(to) == 0 {
			continue
		}
		index := strings.Index(returnString, from)
		if index < 0 {
			continue
		}
		returnString = strings.Replace(returnString, from, to, -1)
	}
	if returnString != line {
		println("URL替换：" + replaceHash.ReplaceAllString(line, "") + " -> " + replaceHash.ReplaceAllString(returnString, ""))
	}
	return returnString
}

func isImpl(line string) bool {
	return implReg.MatchString(line)
}
