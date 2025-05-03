package godoc

import (
	"gioui.org/layout"
	"github.com/ddkwork/ux"
	"go/ast"
	"go/parser"
	"go/token"
	"iter"
	"os"
	"path/filepath"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
)

// Godoc todo 单个返回值没有小括号的强制加上了，泛型的函数和方法 not 提取  C:\Program Files\Go\src\go\doc
type Godoc struct {
	Path    string
	Func    string
	Method  string
	Comment string
}

func Layout(libDir string) ux.Widget {
	jsonName := filepath.Base(mylog.Check2(filepath.Abs(libDir)))
	t := ux.NewTreeTable(Godoc{})
	t.TableContext = ux.TableContext[Godoc]{
		CustomContextMenuItems: func(gtx layout.Context, n *ux.Node[Godoc]) iter.Seq[ux.ContextMenuItem] {
			return func(yield func(ux.ContextMenuItem) bool) {

			}
		},
		MarshalRowCells: func(n *ux.Node[Godoc]) (cells []ux.CellData) {

			return ux.MarshalRow(n.Data, func(key string, field any) (value string) {

				return ""
			})
		},
		UnmarshalRowCells: func(n *ux.Node[Godoc], rows []ux.CellData) {
			if n.Container() {
				n.Data.Path = n.SumChildren()
			}
			n.Data = ux.UnmarshalRow[Godoc](rows, func(key, value string) (field any) {
				return nil
			})
		},
		RowSelectedCallback: func() {

		},
		RowDoubleClickCallback: func() {

		},
		SetRootRowsCallBack: func() {
			mylog.Check(os.Chdir(libDir))
			mylog.Check(filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") && !strings.HasSuffix(info.Name(), "_test.go") {
					if countFunctionsAndMethods(path) {
						container := ux.NewContainerNode(path, Godoc{
							Path:    "",
							Func:    "",
							Method:  "",
							Comment: "",
						})
						t.Root.AddChild(container)
						processFile(path, container)
						return err
					}
					processFile(path, t.Root)
				}
				return nil
			}))
		},
		JsonName:   jsonName,
		IsDocument: false,
	}
	return t
}

func countFunctionsAndMethods(filePath string) bool {
	totalFunctions := 0
	totalMethods := 0

	fset := token.NewFileSet()
	node := mylog.Check2(parser.ParseFile(fset, filePath, nil, parser.ParseComments))

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Recv == nil && ast.IsExported(x.Name.Name) { // 排除非导出的函数
				totalFunctions++
			}
		case *ast.GenDecl:
			if x.Tok == token.TYPE {
				for _, spec := range x.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					if st, ok := ts.Type.(*ast.StructType); ok {
						totalMethods += len(st.Fields.List) // 假设结构体的字段即为方法数量
					}
				}
			}
		}
		if x, ok := n.(*ast.FuncDecl); ok {
			if x.Recv != nil {
				totalMethods++
			}
		}
		return true
	})

	return totalFunctions+totalMethods > 1
}

func processFile(filePath string, parent *ux.Node[Godoc]) {
	fset := token.NewFileSet()
	node := mylog.Check2(parser.ParseFile(fset, filePath, nil, parser.ParseComments))
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Recv == nil && ast.IsExported(x.Name.Name) { // 排除非导出的函数
				parent.AddChildByData(Godoc{
					Path:    filePath,
					Func:    formatFuncSignature(x),
					Method:  "",
					Comment: "",
				})
			}
		case *ast.GenDecl: // 处理结构体
			//if x.Tok == token.TYPE {
			//	for _, spec := range x.Specs {
			//		ts, ok := spec.(*ast.TypeSpec)
			//		if !ok {
			//			continue
			//		}
			//		if st, ok := ts.Type.(*ast.StructType); ok {
			//			parent.Append([]string{filePath, fmt.Sprintf("Struct: %s", ts.Name.Name)})
			//			for _, field := range st.Fields.List {
			//				// 处理结构体的字段
			//			}
			//		}
			//	}
			//}
		}
		if x, ok := n.(*ast.FuncDecl); ok {
			if x.Recv != nil {
				recv := x.Recv.List[0].Type
				if ident, ok := recv.(*ast.Ident); ok {
					structName := ident.Name
					if ast.IsExported(x.Name.Name) {
						parent.AddChildByData(Godoc{
							Path:    filePath,
							Func:    "",
							Method:  "func (" + structName + ") " + formatFuncSignature(x),
							Comment: "",
						})
					}
				}
			}
		}
		return true
	})
}

func formatFuncSignature(decl *ast.FuncDecl) string {
	var buf strings.Builder
	buf.WriteString(decl.Name.Name)

	buf.WriteByte('(')
	writeFieldList(&buf, decl.Type.Params)
	buf.WriteByte(')')

	if decl.Type.Results != nil {
		buf.WriteByte(' ')
		buf.WriteByte('(')
		writeFieldList(&buf, decl.Type.Results)
		buf.WriteByte(')')
	}

	return buf.String()
}

func writeFieldList(buf *strings.Builder, list *ast.FieldList) {
	if list != nil {
		for i, p := range list.List {
			if i > 0 {
				buf.WriteString(", ")
			}
			writeField(buf, p)
		}
	}
}

func writeField(buf *strings.Builder, field *ast.Field) {
	for i, name := range field.Names {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(name.Name)
	}
	if field.Type != nil {
		buf.WriteByte(' ')
		writeType(buf, field.Type)
	}
}

func writeType(buf *strings.Builder, typ ast.Expr) {
	switch t := typ.(type) {
	case *ast.Ident:
		buf.WriteString(t.Name)
	case *ast.StarExpr:
		buf.WriteByte('*')
		writeType(buf, t.X)
	case *ast.Ellipsis:
		buf.WriteString("...")
		writeType(buf, t.Elt)
	case *ast.ArrayType:
		if t.Len != nil {
			buf.WriteByte('[')
			writeType(buf, t.Len)
			buf.WriteByte(']')
		}
		writeType(buf, t.Elt)
	case *ast.SelectorExpr:
		writeType(buf, t.X)
		buf.WriteByte('.')
		buf.WriteString(t.Sel.Name)
	case *ast.FuncType:
		buf.WriteString("func(")
		writeFieldList(buf, t.Params)
		buf.WriteByte(')')
		if t.Results != nil {
			if len(t.Results.List) == 1 {
				buf.WriteByte(' ')
			} else {
				buf.WriteString(" (")
			}
			writeFieldList(buf, t.Results)
			if len(t.Results.List) > 1 {
				buf.WriteByte(')')
			}
		}
	default:
		// handle other cases as needed
	}
}
