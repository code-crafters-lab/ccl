package service

import (
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Meta struct {
	ProtoPackage string       // proto 文件中的 package
	JavaPackage  string       // 生成 Java 代码的包名（如 com.example.service）
	ServiceName  string       // 服务名（如 UserService）
	Methods      []MethodMeta // RPC 方法列表
}

type MethodMeta struct {
	Name         string // 方法名（如 getUser）
	Deprecated   bool   // 是否已废弃
	RequestType  string // 请求消息类型（如 GetUserRequest）
	ResponseType string // 响应消息类型（如 UserResponse）
	HttpMethod   string // HTTP 方法（如 GET）
	Url          string // HTTP URL（如 /user/{id}）
}

func NewMeta(service *protogen.Service, file *protogen.File) *Meta {
	fileOptions := file.Proto.GetOptions()
	meta := &Meta{
		ProtoPackage: "",
		JavaPackage:  fileOptions.GetJavaPackage(),
		ServiceName:  strings.TrimSuffix(service.GoName, "Service"),
		Methods:      make([]MethodMeta, len(service.Methods)),
	}
	for i, method := range service.Methods {
		meta.Methods[i] = newMethodMeta(*method)
	}
	return meta
}

func newMethodMeta(method protogen.Method) MethodMeta {
	methodMeta := MethodMeta{
		Name:         method.GoName,
		RequestType:  method.Input.GoIdent.GoName,
		ResponseType: method.Output.GoIdent.GoName,
	}
	methodMeta.extractHttpRule(method.Desc)
	return methodMeta
}

func (mm *MethodMeta) extractHttpRule(methodDesc protoreflect.MethodDescriptor) {
	methodOptions := methodDesc.Options()
	options := methodOptions.(*descriptorpb.MethodOptions)
	if options == nil {
		return
	}

	if options.Deprecated != nil {
		mm.Deprecated = *options.Deprecated
	}

	if proto.HasExtension(methodOptions, annotations.E_Http) {
		// 自定义选项的 "扩展编号"：来自 google/api/annotations.proto 中定义的 http 选项
		// 格式：<protobuf包名>.<选项名> → 对应编号需通过 descriptor 查找
		// 这里直接使用编译后的 Go 常量（annotations.E_Http 等价于 "google.api.http" 的扩展编号）
		ext := proto.GetExtension(methodOptions, annotations.E_Http)
		if ext == nil {
			return // 无 HttpRule 注解
		}

		// 将扩展值转换为 HttpRule 类型
		httpRule, ok := ext.(*annotations.HttpRule)
		if ok {
			mm.setHttpInfo(httpRule)
		}
	}

}

func (mm *MethodMeta) setHttpInfo(httpRule *annotations.HttpRule) {
	// 解析 HTTP 方法（get/post/put/delete/patch 中仅一个非空）
	var (
		method string
		path   string
	)

	switch {
	case httpRule.GetGet() != "":
		method = "GET"
		path = httpRule.GetGet()
	case httpRule.GetPost() != "":
		method = "POST"
		path = httpRule.GetPost()
	case httpRule.GetPut() != "":
		method = "PUT"
		path = httpRule.GetPut()
	case httpRule.GetDelete() != "":
		method = "DELETE"
		path = httpRule.GetDelete()
	case httpRule.GetPatch() != "":
		method = "PATCH"
		path = httpRule.GetPatch()
	}

	mm.HttpMethod = method
	mm.Url = path
}
