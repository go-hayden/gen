package main

import (
	"strconv"
	"strings"
)

type URLMap struct {
	From string `json:"from,omitempty" bson:"from,omitempty"`
	To   string `json:"to,omitempty" bson:"to,omitempty"`
}

type GenConfig struct {
	CurrentPath       string
	AddMparamInit     bool      `json:"add_mparam_init,omitempty" bson:"add_mparam_init,omitempty"`
	GenMapiHeader     bool      `json:"gen_mapi_header,omitempty" bson:"gen_mapi_header,omitempty"`
	GenModelHeader    bool      `json:"gen_model_header,omitempty" bson:"gen_model_header,omitempty"`
	MapiExtraImports  []string  `json:"mapi_extra_imports,omitempty" bson:"mapi_extra_imports,omitempty"`
	ModelExtraImports []string  `json:"model_extra_imports,omitempty" bson:"model_extra_imports,omitempty"`
	URLMaps           []*URLMap `json:"url_maps,omitempty" bson:"url_maps,omitempty"`
}

func (s *GenConfig) print() {
	println("\n****** Gen配置信息 ******")
	println("当前路径：" + s.CurrentPath)
	println("生成 All in one 的 NVNetworkClient+mapi.h 文件：" + strconv.FormatBool(s.GenMapiHeader))
	println("生成 All in one 的 NVModels.h 文件：" + strconv.FormatBool(s.GenModelHeader))
	println("为Babel生成的mparam参数类添加init方法：" + strconv.FormatBool(s.AddMparamInit))
	if s.GenMapiHeader && s.MapiExtraImports != nil && len(s.MapiExtraImports) > 0 {
		println("NVNetworkClient+mapi.h 文件额外添加的头文件：" + strings.Join(s.MapiExtraImports, ", "))
	}
	if s.GenModelHeader && s.ModelExtraImports != nil && len(s.ModelExtraImports) > 0 {
		println("NVModels.h 文件额外添加的头文件：" + strings.Join(s.ModelExtraImports, ", "))
	}
	if s.URLMaps != nil && len(s.URLMaps) > 0 {
		println("URL映射：{")
		for _, urlmap := range s.URLMaps {
			println("  " + urlmap.From + " -> " + urlmap.To)
		}
		println("}")
	}
	println("************************\n")
}
