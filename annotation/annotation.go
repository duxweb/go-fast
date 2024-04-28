package annotation

import (
	"bytes"
	"github.com/gookit/color"
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

var Annotations = make([]*File, 0)

type File struct {
	Name        string
	Annotations []*Annotation
}

type Annotation struct {
	Name   string
	Params map[string]any
	Func   any
}

type Import struct {
	Name string
	As   string
}

func Run() {
	files := make([]*File, 0)
	imports := make([]*Import, 0)
	dir := "./app"
	existingImports := make(map[string]bool)
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

				if !existingImports[fileImport.Name] {
					imports = append(imports, fileImport)
					existingImports[fileImport.Name] = true
				}
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
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if strings.TrimSpace(comment.Text()) != strings.TrimSpace(d.Doc.Text()) {
					continue
				}
				funName := d.Name.Name
				item := parseDocs(comment.Text())
				if item == nil {
					break
				}
				annotation = item
				annotation.Func = asPackage + "." + funName
				break
			case *ast.GenDecl:
				if strings.TrimSpace(comment.Text()) != strings.TrimSpace(d.Doc.Text()) {
					continue
				}

				switch d.Tok {
				case token.CONST, token.VAR:
					for _, spec := range d.Specs {
						if vs, ok := spec.(*ast.ValueSpec); ok {
							for _, ident := range vs.Names {
								varName := ident.Name
								item := parseDocs(comment.Text())
								if item == nil {
									break
								}
								annotation = item
								annotation.Func = asPackage + "." + varName
								break
							}
						}
					}
				case token.TYPE:
					for _, spec := range d.Specs {
						if ts, ok := spec.(*ast.TypeSpec); ok {
							if _, ok = ts.Type.(*ast.StructType); ok {
								structName := ts.Name.Name
								item := parseDocs(comment.Text())
								if item == nil {
									break
								}
								annotation = item
								annotation.Func = asPackage + "." + structName + "{}"
								break
							}
						}
					}
				}
			}

		}

		// 如果不是函数或结构体的注释，则直接解析注释内容
		if annotation == nil {
			common := parseDocs(comment.Text())
			if common != nil {
				annotation = common
			}
		}

		if annotation == nil {
			continue
		}
		fileAnnotations = append(fileAnnotations, annotation)
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
	if name == "" {
		return nil
	}
	return &Annotation{
		Name:   name,
		Params: data,
	}
}

func generateIndexFile(files []*File, imports []*Import) {

	tmpl := `
package runtime

import (
	{{- range .imports}}
	{{.As}}	"{{.Name}}"
	{{- end}}
	"github.com/duxweb/go-fast/annotation"
)



var Annotations = []*annotation.File{
	{{- range $file := .files}}
	{
		Name: "{{$file.Name}}",
		Annotations: []*annotation.Annotation{
			{{- range $item := $file.Annotations }}
			{
				Name: "{{$item.Name}}",
				Params: map[string]any{
					{{- range $key, $value := $item.Params }}
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

	var buf bytes.Buffer
	var err error
	err = t.Execute(&buf, map[string]any{
		"files":   files,
		"imports": imports,
	})
	if err != nil {
		log.Fatal(err)
	}

	outputSize := int64(buf.Len())

	fileInfo, _ := os.Stat("./runtime/annotations.go")
	fileSize := fileInfo.Size()

	if outputSize != fileSize {
		color.Redln("⇨ runtime found an update, please restart")
	}

	file, err := os.Create("./runtime/annotations.go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(buf.Bytes())
	if err != nil {
		return
	}
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
