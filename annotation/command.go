package annotation

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
)

func Command() []*cli.Command {
	gen := &cli.Command{
		Name:  "annotation:gen",
		Usage: "Generating Annotated Data",
		Action: func(ctx *cli.Context) error {
			fset := token.NewFileSet()
			path, _ := filepath.Abs("./app/home/web/index.go")
			f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				log.Println(err)
			}

			for _, c := range f.Comments {
				fmt.Println("comment: ", c.Text())
			}

			ast.Inspect(f, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.FuncDecl:
					if x.Doc != nil {
						fmt.Println("func comment: ", x.Doc.Text())
					}
				case *ast.Field:
					if x.Doc != nil {
						fmt.Println("field comment: ", x.Doc.Text())
					}
				}

				return true
			})

			return nil
		},
	}

	return []*cli.Command{
		gen,
	}
}
