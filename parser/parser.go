package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

const (
	structComment     = "easyjson:json"
	structSkipComment = "easyjson:skip"
	structDataComment = "easyjson:model"
)

type Parser struct {
	PkgPath     string
	PkgName     string
	StructNames []string
	Models      []string
	AllStructs  bool
}

type visitor struct {
	*Parser

	name string
}

func (p *Parser) needType(comments *ast.CommentGroup) (skip, explicit, model bool) {
	if comments == nil {
		return
	}

	for _, v := range comments.List {
		comment := v.Text

		if len(comment) > 2 {
			switch comment[1] {
			case '/':
				// -style comment (no newline at the end)
				comment = comment[2:]
			case '*':
				/*-style comment */
				comment = comment[2 : len(comment)-2]
			}
		}

		for _, comment := range strings.Split(comment, "\n") {
			comment = strings.TrimSpace(comment)

			if strings.HasPrefix(comment, structSkipComment) {
				return true, false, false
			}
			if strings.HasPrefix(comment, structComment) {
				return false, true, false
			}
			if strings.HasPrefix(comment, structDataComment) {
				return false, false, true
			}
		}
	}

	return
}

func (v *visitor) Visit(n ast.Node) (w ast.Visitor) {
	switch n := n.(type) {
	case *ast.Package:
		return v
	case *ast.File:
		v.PkgName = n.Name.String()
		return v

	case *ast.GenDecl:
		skip, explicit, model := v.needType(n.Doc)

		if skip || explicit || model {
			for _, nc := range n.Specs {
				switch nct := nc.(type) {
				case *ast.TypeSpec:
					nct.Doc = n.Doc
				}
			}
		}

		return v
	case *ast.TypeSpec:
		skip, explicit, model := v.needType(n.Doc)
		if skip {
			return nil
		}
		if !explicit && !v.AllStructs {
			return nil
		}

		v.name = n.Name.String()

		// Allow to specify non-structs explicitly independent of '-all' flag.
		if explicit && !model {
			v.StructNames = append(v.StructNames, v.name)
			return nil
		}
		if model {
			v.Models = append(v.Models, v.name)
		}

		return v
	case *ast.StructType:
		v.StructNames = append(v.StructNames, v.name)
		return nil
	}
	return nil
}

func (p *Parser) Parse(fname string, isDir bool) error {
	var err error
	if p.PkgPath, err = getPkgPath(fname, isDir); err != nil {
		return err
	}

	fset := token.NewFileSet()
	if isDir {
		packages, err := parser.ParseDir(fset, fname, excludeTestFiles, parser.ParseComments)
		if err != nil {
			return err
		}

		for _, pckg := range packages {
			ast.Walk(&visitor{Parser: p}, pckg)
		}
	} else {
		f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		ast.Walk(&visitor{Parser: p}, f)
	}
	return nil
}

func excludeTestFiles(fi os.FileInfo) bool {
	return !strings.HasSuffix(fi.Name(), "_test.go")
}
