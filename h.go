package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path"

	"strings"

	tb "github.com/go-hayden/toolbox"
)

func generateMapiHeader(cfg *GenConfig, c chan error) {
	if !cfg.GenMapiHeader {
		c <- nil
		return
	}
	mapiMPath := path.Join(cfg.CurrentPath, "NVNetwork", "NVNetworkClient+mapi.m")
	if !tb.FileExists(mapiMPath) {
		err := errors.New("无法找到NVNetworkClient+mapi.m文件")
		c <- err
		return
	}
	var buffer bytes.Buffer
	buffer.WriteString("\n//This file is generated, never modify this mannualy\n\n")
	for _, ei := range cfg.MapiExtraImports {
		tmp := strings.TrimSpace(ei)
		if len(tmp) > 0 {
			buffer.WriteString("#import \"" + strings.TrimSpace(ei) + "\"\n")
		}
	}
	tb.ReadLine(mapiMPath, func(line string, finished bool, err error, stop *bool) {
		if isImport(line) {
			buffer.WriteString(line + "\n")
		}
	})

	mapiHPath := path.Join(cfg.CurrentPath, "NVNetworkClient+mapi.h")
	err := ioutil.WriteFile(mapiHPath, buffer.Bytes(), 0644)
	if err != nil {
		c <- err
		return
	}

	c <- nil
}

func generateModelHeader(cfg *GenConfig, c chan error) {
	if !cfg.GenModelHeader {
		c <- nil
		return
	}
	modelMPath := path.Join(cfg.CurrentPath, "NVModels", "NVModels.m")
	if !tb.FileExists(modelMPath) {
		err := errors.New("无法找到NVModels.m文件")
		c <- err
		return
	}
	var buffer bytes.Buffer
	buffer.WriteString("\n//This file is generated, never modify this mannualy\n\n")
	for _, ei := range cfg.ModelExtraImports {
		tmp := strings.TrimSpace(ei)
		if len(tmp) > 0 {
			buffer.WriteString("#import \"" + strings.TrimSpace(ei) + "\"\n")
		}
	}
	buffer.WriteString("#import \"NVObject.h\"\n")
	buffer.WriteString("#import \"NVBaseModel.h\"\n\n")

	tb.ReadLine(modelMPath, func(line string, finished bool, err error, stop *bool) {
		if isImport(line) {
			buffer.WriteString(line + "\n")
		}
	})

	modelHPath := path.Join(cfg.CurrentPath, "NVModels.h")
	err := ioutil.WriteFile(modelHPath, buffer.Bytes(), 0644)
	if err != nil {
		c <- err
		return
	}

	c <- nil
}

func isImport(line string) bool {
	return importReg.MatchString(line)
}
