package main

import (
	"ccl/db/ent/extension"
	"embed"
	"os"
	"path/filepath"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	err    error
	//go:embed templates/**
	templateDir embed.FS
)

func main() {
	if !clean() {
		return
	}
	log := logger.Sugar()
	opts := []entc.Option{
		//entc.Dependency(
		//	entc.DependencyType(&zap.Logger{}),
		//),
		entc.Extensions(extension.Extension(log)),
	}
	// 2. 递归遍历源目录
	template := gen.MustParse(gen.NewTemplate("ccl").
		ParseFS(templateDir, "templates/*.tmpl", "templates/*/*.tmpl"))
	log.Info(template)

	if err := entc.Generate("./ent/schema", &gen.Config{
		Features: []gen.Feature{
			//gen.FeaturePrivacy,
			gen.FeatureEntQL,
			gen.FeatureGlobalID,
			gen.FeatureModifier,
		},
		Templates: []*gen.Template{template},
	}, opts...); err != nil {
		log.Error("running ent codegen:", zap.Error(err))
	}
	log.Infof("generated code for schema in %s", "./ent/schema")
}

func init() {
	logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger.Info("ent 代码生成", zap.String("version", "v1.0.0"))
}

// 需要保留的目录和文件
var keep = map[string]bool{
	"schema":      true, // 目录
	"extension":   true, // 目录
	"generate.go": true, // 文件
}

func clean() bool {
	logger.Debug("正在清理已生成文件")
	// 定义要清理的目标目录
	entDir := "./ent"
	// 读取 ent 目录下的所有条目
	entries, _ := os.ReadDir(entDir)
	logger.Debug("--- 开始清理 ent 目录 ---")
	// 遍历每个条目
	for _, entry := range entries {
		name := entry.Name()

		// 如果条目是需要保留的，则跳过
		if keep[name] {
			logger.Info("已保留", zap.String("name", name))
			continue
		}

		// 构建完整的路径
		fullPath := filepath.Join(entDir, name)

		// 判断是文件还是目录，并执行删除
		if entry.IsDir() {
			// 删除整个目录及其内容
			if err := os.RemoveAll(fullPath); err != nil {
				logger.Error("删除目录失败", zap.Error(err))
				return false
			}
		} else {
			// 删除单个文件
			if err := os.Remove(fullPath); err != nil {
				logger.Error("删除文件失败", zap.Error(err))
				return false
			}
		}
	}
	return true
}
