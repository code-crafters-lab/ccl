package {{.JavaPackage}}.controller;

import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import javax.annotation.Resource;
{{ if .ProtoPackage -}}import {{.ProtoPackage}}.{{.RequestType}};{{ end -}}
{{ if .ProtoPackage -}}import {{.ProtoPackage}}.{{.ResponseType}};{{ end -}}
import {{.JavaPackage}}.service.{{.ServiceName}}Service;

{{ define "resource" -}}
    @Resource
    private {{.ServiceName}}Service {{.ServiceName | lower}}Service;
{{ end -}}

{{ define "method" -}}
    @PostMapping("{{.Url}}")
    public {{.ResponseType}} {{.Name}}(@RequestBody {{.RequestType}} request) {

    }
{{ end -}}

@RestController
public class {{.ServiceName}}Controller {

    {{template "resource" .}}

    {{range .Methods -}}
    {{template "method" . }}
    {{end}}
}