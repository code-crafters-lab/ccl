package main

import (
	"ccl/db/extension"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"go.uber.org/zap"
)

var logger *zap.Logger
var err error

func main() {
	log := logger.Sugar()
	opts := []entc.Option{
		entc.Extensions(extension.Extension(log)),
	}
	if err := entc.Generate("./ent/schema", &gen.Config{
		Package: "ccl/db/ent/generated",
		Target:  "./ent/generated",
		Features: []gen.Feature{
			//gen.FeaturePrivacy,
			gen.FeatureEntQL,
			gen.FeatureGlobalID,
			gen.FeatureModifier,
		},
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
