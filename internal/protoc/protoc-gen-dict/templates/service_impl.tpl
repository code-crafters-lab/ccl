package {{.JavaPackage}}.service.impl;

import org.springframework.stereotype.Service;
{{if .ProtoPackage}}import {{.ProtoPackage}}.{{.RequestType}};{{end}}
{{if .ProtoPackage}}import {{.ProtoPackage}}.{{.ResponseType}};{{end}}
    import {{.JavaPackage}}.service.{{.ServiceName}}Service;

    @Service
    public class {{.ServiceName}}ServiceImpl implements {{.ServiceName}}Service {

    {{range .Methods}}
    @Override
    public {{.ResponseType}} {{.Name}}({{.RequestType}} request) {
    // TODO: 实现业务逻辑
    return {{.ResponseType}}.newBuilder().build();
}
    {{end}}
}