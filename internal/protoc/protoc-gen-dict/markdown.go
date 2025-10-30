package main

import (
	"bytes"
	"fmt"
	"path"

	"google.golang.org/protobuf/compiler/protogen"
)

// 生成 markdown
func generateMarkdown(dictionaries []*dict, plugin *protogen.Plugin, file *protogen.File, opts option) {
	if !opts.markdown {
		return
	}
	filename := fmt.Sprintf("%s.dict.md", file.GeneratedFilenamePrefix)
	if opts.java {
		filename = fmt.Sprintf("%s.dict.md", path.Base(file.GeneratedFilenamePrefix))
	}
	g := plugin.NewGeneratedFile(filename, file.GoImportPath)

	buf := bytes.Buffer{}
	for i, d := range dictionaries {
		buf.WriteString(fmt.Sprintf("## %d. %s\n", i+1, d.Name))
		buf.WriteByte('\n')
		buf.WriteString(fmt.Sprintf("- 字典编码：%s\n", d.FullCode))
		if d.Description != "" {
			buf.WriteString(fmt.Sprintf("- 字典编码：%s\n", d.Description))
		}
		buf.WriteByte('\n')

		buf.WriteString("| 名称 | 编码 | 值 | 排序 | 描述 |\n")
		buf.WriteString("| ---- | ---- | :----: | :----: | ---- |\n")
		for _, item := range d.Items {
			buf.WriteString(fmt.Sprintf("| %s | %v | %s | %d | %s |\n", item.Name, item.Code, item.Value, item.Sort, item.Description))
		}
	}

	_, _ = g.Write(buf.Bytes())
}
