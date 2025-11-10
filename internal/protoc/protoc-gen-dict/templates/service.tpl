package {{.JavaPackage}}.service;

{{if .ProtoPackage}}import {{.ProtoPackage}}.{{.RequestType}};{{end}}
{{if .ProtoPackage}}import {{.ProtoPackage}}.{{.ResponseType}};{{end}}

public interface {{.ServiceName}}Service {
    {{range .Methods}}
    {{.ResponseType}} {{.Name}}({{.RequestType}} request);
    {{end}}
}