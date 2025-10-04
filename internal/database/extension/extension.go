package extension

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"go.uber.org/zap"
)

type cclExtension struct {
	logger *zap.Logger
	entc.Extension
}

func (ccl cclExtension) Hooks() []gen.Hook {
	return []gen.Hook{
		FiledSortHook(ccl.logger.Sugar()),
	}
}

// Annotations of the extensions.
func (cclExtension) Annotations() []entc.Annotation { return []entc.Annotation{} }

// Templates of the extensions.
func (cclExtension) Templates() []*gen.Template { return []*gen.Template{} }

// Options of the extensions.
func (cclExtension) Options() []entc.Option { return []entc.Option{} }

func Extension(logger *zap.SugaredLogger) entc.Extension {
	return &cclExtension{
		logger: logger.Desugar(),
	}
}
