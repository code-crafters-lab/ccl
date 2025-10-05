package extension

import (
	"sort"
	"strings"

	"entgo.io/ent/entc/gen"
	"go.uber.org/zap"
)

type sortField struct {
	Name string
	Sort int
}

var sortFields = []sortField{
	{"id", 1},
	{"tenant_id", 2},
	{"name", 3},
	{"sort", -170},
	{"status", -160},
	{"description", -150},
	{"assigned_by", -140},
	{"expired_at", -130},
	{"last_login_at", -120},
	{"created_at", -104},
	{"created_by", -103},
	{"updated_at", -102},
	{"updated_by", -101},
}

// FiledSortHook 字段确保所有的字段都有标签
func FiledSortHook(logger *zap.SugaredLogger) gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			// 创建一个权重映射，方便快速查找
			weightMap := make(map[string]int)
			for _, sf := range sortFields {
				weightMap[sf.Name] = sf.Sort
			}
			for _, node := range g.Nodes {
				logger.Debugf("正在为表 %s 按权重重新排序字段...", node.Name)

				a := make([]string, len(node.Fields))
				for i, field := range node.Fields {
					a[i] = field.Name
				}
				logger.Debugf("排序前：%s", strings.Join(a, ", "))

				// 获取节点的所有字段
				fields := node.Fields

				// 3. 使用 sort.Slice 进行自定义排序
				sort.Slice(fields, func(i, j int) bool {
					fieldI := fields[i]
					fieldJ := fields[j]

					// 从权重映射中获取两个字段的权重
					weightI, hasI := weightMap[fieldI.Name]
					weightJ, hasJ := weightMap[fieldJ.Name]

					// 情况1: 两个字段都有指定权重 (正数或负数)
					if hasI && hasJ {
						// a. 如果一个是正数，一个是负数，正数在前
						if weightI > 0 && weightJ < 0 {
							return true // I是正数，J是负数，I排在J前面
						}
						if weightI < 0 && weightJ > 0 {
							return false // I是负数，J是正数，I排在J后面
						}
						// b. 如果两者同号（都为正或都为负），则按数值从小到大排序
						return weightI < weightJ
					}

					// 情况2: 只有一个字段有指定权重
					if hasI {
						// I 有权重，J 没有。
						// 如果 I 的权重是正数，它应该排在 J 前面。
						// 如果 I 的权重是负数，它应该排在 J 后面。
						return weightI > 0
					}
					if hasJ {
						// J 有权重，I 没有。
						// 如果 J 的权重是正数，它应该排在 I 前面，所以 I 应该在 J 后面，返回 false。
						// 如果 J 的权重是负数，它应该排在 I 后面，所以 I 应该在 J 前面，返回 true。
						return weightJ < 0
					}

					// 情况3: 两个字段都不在排序规则中
					// 保持它们在 Schema 中定义的原始顺序
					return fieldI.Position.Index < fieldJ.Position.Index
				})

				// 4. 将排序后的字段切片替换回节点
				node.Fields = fields

				b := make([]string, len(fields))
				for i, field := range fields {
					b[i] = field.Name
				}
				logger.Debugf("排序后：%s", strings.Join(b, ", "))

			}
			return next.Generate(g)
		})
	}
}
