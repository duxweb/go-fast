package annotation

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

type File struct {
	Name        string
	Annotations []*Annotation
}

type Annotation struct {
	Name string
	Data map[string]any
	Func any
}

type Import struct {
	Name string
	As   string
}

func Run() {
	files := make([]*File, 0)
	imports := make([]*Import, 0)
	dir := "./app"
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			// 解析Go文件
			fileAnnotations, fileImport := parseGoFile(dir, path)
			if err != nil {
				log.Printf("Error parsing file %s: %s\n", path, err)
				return nil
			}
			if fileAnnotations != nil {
				files = append(files, &File{
					Name:        fileImport.Name,
					Annotations: fileAnnotations,
				})
				imports = append(imports, fileImport)
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	generateIndexFile(files, imports)

}

func parseGoFile(dir string, path string) ([]*Annotation, *Import) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		return nil, nil
	}

	relPackagePath, err := filepath.Rel(dir, filepath.Dir(path))

	fullPackageName := "dux-project/app/"
	fullPackageName += filepath.ToSlash(relPackagePath)

	parts := strings.Split(relPackagePath, "/")
	for i, part := range parts {
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	asPackage := "app"
	asPackage += strings.Join(parts, "")
	Imports := &Import{
		Name: fullPackageName,
		As:   asPackage,
	}

	var fileAnnotations []*Annotation

	for _, comment := range file.Comments {
		var annotation *Annotation
		for _, decl := range file.Decls {
			declFun, ok := decl.(*ast.FuncDecl)
			if !ok {
				// 不为函数跳过
				continue
			}
			// 判断函数所属注释
			if strings.TrimSpace(comment.Text()) == strings.TrimSpace(declFun.Doc.Text()) {
				funName := declFun.Name.Name
				annotation = parseDocs(comment.Text())
				annotation.Func = asPackage + "." + funName
				break
			}
		}

		// 如果不是函数或方法的注释，则直接解析注释内容
		if annotation == nil {
			annotation = parseDocs(comment.Text())
		}

		if annotation != nil {
			fileAnnotations = append(fileAnnotations, annotation)
		}
	}

	return fileAnnotations, Imports
}

func parseDocs(text string) *Annotation {
	if !strings.Contains(text, "@") {
		return nil
	}
	reg := regexp.MustCompile(`@(\w+)\(([^()]*)\)`)
	match := reg.FindStringSubmatch(text)
	if len(match) < 3 {
		return nil
	}
	name := match[1]
	data := parseParams(match[2])
	return &Annotation{
		Name: name,
		Data: data,
	}
}

func generateIndexFile(files []*File, imports []*Import) {
	tmpl := `
package main

import (
	{{- range .imports}}
	{{.As}}	"{{.Name}}"
	{{- end}}
)


type Annotation struct {
	Name string
	Params map[string]any
	Func any
}

type File struct {
	Name        string
	Annotations []*Annotation
}

var Annotations = []*File{
	{{- range $file := .files}}
	{
		Name: "{{$file.Name}}",
		Annotations: []*Annotation{
			{{- range $item := $file.Annotations }}
			{
				Name: "{{$item.Name}}",
				Params: map[string]any{
					{{- range $key, $value := $item.Data }}
					"{{ $key }}": {{ $value }},
					{{- end }}
				},
				{{- if $item.Func }}
				Func: {{$item.Func}},
				{{- end }}
			},
			{{- end}}
		},
	},
	{{- end}}
	
}
`

	t := template.Must(template.New("index").Parse(tmpl))

	file, err := os.Create("./runtime/annotations.go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = t.Execute(file, map[string]any{
		"files":   files,
		"imports": imports,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Index file generated successfully.")
}

func parseParams(paramString string) map[string]interface{} {
	params := make(map[string]interface{})
	parts := strings.Split(paramString, ",")
	for _, part := range parts {
		keyValue := strings.SplitN(part, "=", 2)
		if len(keyValue) == 2 {
			key := strings.TrimSpace(keyValue[0])
			value := strings.TrimSpace(keyValue[1])
			params[key] = value
		}
	}
	return params
}
