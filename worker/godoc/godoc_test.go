package godoc

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type (
	// 定义数据类型
	Function struct {
		Name   string
		Params string
		Return string
		Body   string
	}

	Method struct {
		Recv   string
		Name   string
		Params string
		Return string
		Body   string
	}

	Interface struct {
		Name    string
		Methods []Method
	}

	Struct struct {
		Name   string
		Fields string
	}

	FileData struct {
		Path       string
		Content    string // 存储文件内容
		Functions  []Function
		Methods    []Method
		Interfaces []Interface
		Structs    []Struct
	}

	FileResult map[string]FileData
)

func TestDoc(t *testing.T) {
	Walk("golibrary")
}
func Walk(root string) { //todo gen std lib provider for cpp library
	os.RemoveAll("tmp")
	result := make(FileResult)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Base(path) == "provider.go" {
			return nil
		}
		if strings.HasSuffix(path, "_gen.go") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			if data, err := parseGoFile(path); err == nil {
				result[path] = data
			} else {
				fmt.Printf("Error parsing %s: %v\n", path, err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking path: %v\n", err)
	}

	/*
		i := 0
		for path, data := range result {
			i++
			//gen provider
			const provider = "Provider"
			getNamer := func(name string) (objectName, exportInitVarName string) {
				name = stream.TrimExtension(name)
				switch name {
				case "datasize":
					name = "DataSize"
				case "httpclient":
					name = "HTTPClient"
				case "logger":
					name = "Logger"
				case "redis":
					name = "Redis"
				case "byteslice":
					name = "ByteSlice"

				}
				return stream.ToCamelToLower(name) + provider, stream.ToCamelUpper(name)
			}
			base := filepath.Base(path)
			base = stream.TrimExtension(base)
			objectName, exportInitVarName := getNamer(stream.TrimExtension(base))

			g := stream.NewGeneratedFile()
			if len(data.Functions) > 0 {
				g.P("type ", objectName, " struct {")
				for _, function := range data.Functions {
					switch function.Name {
					case "String": //fmt.Stringer
						continue
					case "object":
						function.Name = stream.ToCamelUpper(base)
					}
					g.P(function.Name, " func ", function.Params, function.Return)
				}
				g.P("}")

				g.P("var ", exportInitVarName, "=", objectName, "{")
				for _, function := range data.Functions {
					switch function.Name {
					case "String": //fmt.Stringer
						continue
					case "*object":
						function.Name = stream.ToCamelUpper(base)
					}
					g.P(function.Name, ": func", function.Params, function.Return, "{")
					g.P(strings.TrimPrefix(function.Body, "{"), ",")
				}
				g.P("}")
				g.P()
			}

			if len(data.Methods) > 0 {
				for j, method := range data.Methods {
					if j == 0 {
						fieldName := ""
						switch {
						case strings.Contains(method.Recv, "["):
							//(s *M[K, V])
							split := strings.Split(method.Recv, "[")
							split = strings.Split(split[0], "*")
							fieldName = split[1]
						case strings.Contains(method.Recv, "*"):
							split := strings.Split(method.Recv, "*")
							fieldName = strings.TrimSuffix(split[1], ")")
						default:
							//(b BitField)
							split := strings.Split(method.Recv, " ")
							fieldName = strings.TrimSuffix(split[1], ")")
						}
						objectName, exportInitVarName = getNamer(fieldName)
						g.P("type ", objectName, " struct {")
					}
					switch method.Name {
					case "String": //fmt.Stringer
						continue
					case "*object":
						method.Name = stream.ToCamelUpper(base)
					}
					g.P(method.Name, " func ", method.Params, method.Return)
				}
				g.P("}")

				g.P("var ", exportInitVarName, "=", objectName, "{")
				for _, method := range data.Methods {
					switch method.Name {
					case "String": //fmt.Stringer
						continue
					case "*object":
						method.Name = stream.ToCamelUpper(base)
					}
					g.P(method.Name, ": func", method.Params, method.Return, "{")
					g.P(strings.TrimPrefix(method.Body, "{"), ",")
				}
				g.P("}")
				g.P()
			}
			if stream.FileExists("tmp/generated/" + base) {
				base = stream.TrimExtension(base) + strconv.Itoa(i) + ".go"
			}
			g.InsertPackageWithImports("golibrary")
			stream.WriteGoFile("tmp/lib/"+base+".go", g.String())
		}
	*/

	// 打印结果
	for path, data := range result {
		fmt.Printf("\nFile: %s\n", path)
		fmt.Println("Functions:")
		for _, f := range data.Functions {
			fmt.Printf("  %s(%s) %s\n", f.Name, f.Params, f.Return)
		}

		fmt.Println("\nMethods:")
		for _, m := range data.Methods {
			fmt.Printf("  (%s) %s(%s) %s\n", m.Recv, m.Name, m.Params, m.Return)
		}

		fmt.Println("\nInterfaces:")
		for _, i := range data.Interfaces {
			fmt.Printf("  interface %s\n", i.Name)
			for _, m := range i.Methods {
				fmt.Printf("    %s(%s) %s\n", m.Name, m.Params, m.Return)
			}
		}

		fmt.Println("\nStructs:")
		for _, s := range data.Structs {
			fmt.Printf("  struct %s\n", s.Name)
			if s.Fields != "" {
				// 格式化字段输出
				fieldLines := strings.SplitSeq(s.Fields, "\n")
				for line := range fieldLines {
					fmt.Printf("    %s\n", strings.TrimSpace(line))
				}
			}
		}
		fmt.Println()
	}
}

// 从AST节点获取对应的源代码
func getNodeCode(fset *token.FileSet, node ast.Node, content string) string {
	if node == nil {
		return ""
	}
	start := fset.Position(node.Pos()).Offset
	end := fset.Position(node.End()).Offset
	return strings.TrimSpace(content[start:end])
}

// 获取字段列表的代码
func getFieldListCode(fset *token.FileSet, fl *ast.FieldList, content string) string {
	if fl == nil || len(fl.List) == 0 {
		return ""
	}
	return getNodeCode(fset, fl, content)
}

// 主要解析函数
func parseGoFile(path string) (FileData, error) {
	data := FileData{
		Path: path,
	}

	// 读取文件内容
	content, err := os.ReadFile(path)
	if err != nil {
		return data, err
	}
	data.Content = string(content)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return data, err
	}

	// 遍历所有顶级声明
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					switch t := ts.Type.(type) {
					case *ast.StructType:
						fields := getNodeCode(fset, t, data.Content)
						data.Structs = append(data.Structs, Struct{
							Name:   ts.Name.Name,
							Fields: fields,
						})
					case *ast.InterfaceType:
						iface := Interface{
							Name: ts.Name.Name,
						}
						for _, m := range t.Methods.List {
							if len(m.Names) > 0 {
								method := getMethodFromField(fset, m, data.Content)
								iface.Methods = append(iface.Methods, method)
							}
						}
						data.Interfaces = append(data.Interfaces, iface)
					}
				}
			}

		case *ast.FuncDecl:
			if d.Recv == nil {
				// 普通函数
				name := d.Name.Name
				params := getFieldListCode(fset, d.Type.Params, data.Content)
				returns := getFieldListCode(fset, d.Type.Results, data.Content)

				if params == "" {
					params = "()"
				}
				data.Functions = append(data.Functions, Function{
					Name:   name,
					Params: params,
					Return: returns,
					Body:   getNodeCode(fset, d.Body, data.Content),
				})
			} else {
				// 方法
				recv := getFieldListCode(fset, d.Recv, data.Content)
				name := d.Name.Name
				params := getFieldListCode(fset, d.Type.Params, data.Content)
				returns := getFieldListCode(fset, d.Type.Results, data.Content)

				if params == "" {
					params = "()"
				}
				data.Methods = append(data.Methods, Method{
					Recv:   recv,
					Name:   name,
					Params: params,
					Return: returns,
					Body:   getNodeCode(fset, d.Body, data.Content),
				})
			}
		}
	}

	return data, nil
}

// 从字段中获取方法（用于接口方法）
func getMethodFromField(fset *token.FileSet, f *ast.Field, content string) Method {
	if len(f.Names) == 0 {
		return Method{}
	}

	name := f.Names[0].Name
	var params, returns string

	if ft, ok := f.Type.(*ast.FuncType); ok {
		params = getFieldListCode(fset, ft.Params, content)
		returns = getFieldListCode(fset, ft.Results, content)
	}

	return Method{
		Name:   name,
		Params: params,
		Return: returns,
	}
}
