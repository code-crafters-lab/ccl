package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func tplOut(outputDir string, pkg string, filename string) string {
	return filepath.FromSlash(filepath.Join(outputDir, strings.ReplaceAll(pkg, ".", "/"), filename))
}

func GenerateJavaCode(meta *Meta, outputDir string) error {
	// 定义模板映射（模板文件名 → 生成的 Java 文件路径）
	templateMap := map[string]string{
		"templates/controller.tpl":   tplOut(outputDir, fmt.Sprintf("%s.controller", meta.JavaPackage), fmt.Sprintf("%sController.java", meta.ServiceName)),
		"templates/service.tpl":      tplOut(outputDir, fmt.Sprintf("%s.service", meta.JavaPackage), fmt.Sprintf("%sService.java", meta.ServiceName)),
		"templates/service_impl.tpl": tplOut(outputDir, fmt.Sprintf("%s.service.impl", meta.JavaPackage), fmt.Sprintf("%sServiceImpl.java", meta.ServiceName)),
	}

	// 1. 定义你的自定义函数
	// 这里我们直接用 strings.ToLower，也可以包装一下
	lower := func(s string) string {
		return strings.ToLower(s)
	}
	// indent 函数将给定的文本块的每一行都缩进指定的级别。
	// level: 缩进级别，每个级别对应 4 个空格。
	// text: 需要被缩进的原始文本。
	indent := func(level int, text string) string {
		if level <= 0 || text == "" {
			return text
		}

		// 根据级别生成缩进字符串
		indentStr := strings.Repeat(" ", level*4)

		// 按行分割文本
		lines := strings.Split(text, "\n")

		// 对每一行进行缩进处理
		for i, line := range lines {
			// 只有非空行才进行缩进，避免在空行上添加多余的空格
			if line != "" {
				lines[i] = indentStr + line
			}
		}

		// 将处理后的行重新拼接成一个字符串
		return strings.Join(lines, "\n")
	}

	// 2. 创建一个 FuncMap 并将你的函数注册进去
	funcMap := template.FuncMap{
		"lower":  lower,
		"indent": indent,
	}

	// 加载并解析所有模板
	tpl, err := template.New("example").Funcs(funcMap).ParseGlob("templates/*.tpl")

	if err != nil {
		return err
	}

	// 渲染每个模板，生成文件
	for tplPath, javaFilePath := range templateMap {
		// 创建输出目录（递归创建多级目录）
		if err := os.MkdirAll(filepath.Dir(javaFilePath), 0755); err != nil {
			return err
		}

		// 创建 Java 文件
		file, err := os.Create(javaFilePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// 渲染模板并写入文件
		if err := tpl.ExecuteTemplate(file, filepath.Base(tplPath), meta); err != nil {
			return err
		}
	}

	return nil
}

// 将 Java 包名转为文件路径（com.example → com/example）
func replacePackageToPath(pkg string) string {
	return filepath.FromSlash(pkg)
}
