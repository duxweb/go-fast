package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var (
	dir = flag.String("dir", "", "directory to scan (default: current directory)")
)

type FileData struct {
	Name        string
	Annotations []*AnnotationData
}

type AnnotationData struct {
	Name   string
	Params map[string]interface{}
	Func   string
}

type Import struct {
	Name string
	As   string
}

func main() {
	flag.Parse()

	// 自动检测项目根目录
	scanDir := *dir
	if scanDir == "" {
		scanDir = "."
	}
	
	// 确保 runtime 目录存在
	runtimeDir := filepath.Join(scanDir, "runtime")
	if err := os.MkdirAll(runtimeDir, 0755); err != nil {
		log.Fatalf("Failed to create runtime directory: %v", err)
	}
	
	// 输出文件固定为 runtime/annotations.go
	outputFile := filepath.Join(runtimeDir, "annotations.go")
	
	// 自动检测包名
	packageName := "runtime"

	files := []*FileData{}
	imports := []*Import{}
	importMap := make(map[string]bool)

	// 扫描 app 目录
	appDir := filepath.Join(scanDir, "app")
	if _, err := os.Stat(appDir); err == nil {
		err := filepath.WalkDir(appDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
				return err
			}

			// 跳过特殊目录
			if strings.Contains(path, "vendor/") || 
			   strings.Contains(path, ".git/") {
				return nil
			}

			annotations, importInfo := parseGoFile(scanDir, path)
			if len(annotations) == 0 {
				return nil
			}

			relPath, _ := filepath.Rel(scanDir, path)
			files = append(files, &FileData{
				Name:        filepath.ToSlash(relPath),
				Annotations: annotations,
			})

			if importInfo != nil && !importMap[importInfo.Name] {
				imports = append(imports, importInfo)
				importMap[importInfo.Name] = true
			}

			return nil
		})

		if err != nil {
			log.Fatal(err)
		}
	}

	if err := generateFile(outputFile, packageName, files, imports); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated %s with %d files\n", outputFile, len(files))
}

var annotationRe = regexp.MustCompile(`@(\w+)\(([^()]*)\)`)

func parseGoFile(rootDir, path string) ([]*AnnotationData, *Import) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, nil
	}

	// 计算包导入信息
	relPath, _ := filepath.Rel(rootDir, filepath.Dir(path))
	relPath = filepath.ToSlash(relPath)
	
	var importInfo *Import
	var pkgAlias string
	
	// 只处理 app 目录下的文件
	if strings.HasPrefix(relPath, "app") {
		// 生成包别名：app/system/web -> AppSystemWeb
		parts := strings.Split(relPath, "/")
		for i, part := range parts {
			if part != "" {
				parts[i] = strings.ToUpper(part[:1]) + part[1:]
			}
		}
		pkgAlias = strings.Join(parts, "")
		
		// 获取模块名（从 go.mod 中读取）
		moduleName := getModuleName(rootDir)
		if moduleName == "" {
			// 如果找不到 go.mod，使用默认项目名
			moduleName = "dux-project"
		}
		
		importInfo = &Import{
			Name: moduleName + "/" + relPath,
			As:   pkgAlias,
		}
	}

	// 收集所有注解
	annotations := []*AnnotationData{}
	
	// 遍历所有注释组
	for _, comment := range node.Comments {
		text := comment.Text()
		if !strings.Contains(text, "@") {
			continue
		}

		matches := annotationRe.FindStringSubmatch(text)
		if len(matches) < 3 {
			continue
		}

		ann := &AnnotationData{
			Name:   matches[1],
			Params: parseParams(matches[2]),
		}

		// 查找关联的声明
		ann.Func = findDeclForComment(node, comment, pkgAlias)
		annotations = append(annotations, ann)
	}

	return annotations, importInfo
}

func findDeclForComment(file *ast.File, comment *ast.CommentGroup, pkgAlias string) string {
	commentText := strings.TrimSpace(comment.Text())
	
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Doc != nil && strings.TrimSpace(d.Doc.Text()) == commentText {
				return formatName(d.Name.Name, pkgAlias)
			}
		case *ast.GenDecl:
			if d.Doc != nil && strings.TrimSpace(d.Doc.Text()) == commentText {
				return extractGenDeclName(d, pkgAlias)
			}
		}
	}
	return ""
}

func extractGenDeclName(d *ast.GenDecl, pkgAlias string) string {
	for _, spec := range d.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			if _, ok := s.Type.(*ast.StructType); ok {
				return formatName(s.Name.Name+"{}", pkgAlias)
			}
		case *ast.ValueSpec:
			if len(s.Names) > 0 {
				return formatName(s.Names[0].Name, pkgAlias)
			}
		}
	}
	return ""
}

func formatName(name, pkgAlias string) string {
	if pkgAlias != "" {
		return pkgAlias + "." + name
	}
	return name
}

func parseParams(paramStr string) map[string]interface{} {
	params := make(map[string]interface{})
	if paramStr == "" {
		return params
	}

	for _, part := range strings.Split(paramStr, ",") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.Trim(strings.TrimSpace(kv[1]), `"`)
			params[key] = value
		}
	}
	return params
}

func getModuleName(dir string) string {
	for current := dir; ; {
		goModPath := filepath.Join(current, "go.mod")
		if content, err := os.ReadFile(goModPath); err == nil {
			for _, line := range strings.Split(string(content), "\n") {
				if strings.HasPrefix(strings.TrimSpace(line), "module ") {
					parts := strings.Fields(line)
					if len(parts) >= 2 {
						return parts[1]
					}
				}
			}
		}
		
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return ""
}

func generateFile(outputFile, packageName string, files []*FileData, imports []*Import) error {
	const tmpl = `// Code generated by annotation-gen. DO NOT EDIT.
package {{.Package}}

import (
	"github.com/duxweb/go-fast/v2/annotation"
	{{- range .Imports}}
	{{.As}} "{{.Name}}"
	{{- end}}
)

// GetAnnotations 返回所有注解的索引
func GetAnnotations() []*annotation.File {
	return []*annotation.File{
		{{- range .Files}}
		{
			Name: "{{.Name}}",
			Annotations: []*annotation.Annotation{
				{{- range .Annotations}}
				{
					Name: "{{.Name}}",
					Params: map[string]any{
						{{- range $key, $value := .Params}}
						"{{$key}}": "{{$value}}",
						{{- end}}
					},
					{{- if .Func}}
					Func: {{.Func}},
					{{- end}}
				},
				{{- end}}
			},
		},
		{{- end}}
	}
}
`

	t, err := template.New("index").Parse(tmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, map[string]any{
		"Package": packageName,
		"Files":   files,
		"Imports": imports,
	})
	if err != nil {
		return err
	}

	return os.WriteFile(outputFile, buf.Bytes(), 0644)
}