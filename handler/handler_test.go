package handler

import (
	"aDi/log"
	"aDi/model"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestJson(t *testing.T) {
	// 创建注释提取器
	ce := NewCommentExtractor()

	// 从目录中提取所有.go文件的注释
	err := ce.ExtractCommentsFromDir("../model") // 替换为你的模型目录路径
	if err != nil {
		panic(err)
	}

	// 创建Schema生成器
	generator := NewSchemaGenerator(ce)

	// 为每个需要生成schema的结构体生成schema
	schemas := make(map[string]interface{})

	// 示例：为Person结构体生成schema
	personSchema := generator.GenerateSchema(&model.FPDetailInfo{})
	schemas["Person"] = personSchema
	// 格式化输出
	schemaJSON, err := json.MarshalIndent(personSchema, "", "    ")
	if err != nil {
		panic(err)
	}
	log.Debug(schemaJSON)
}

// CommentExtractor 用于提取结构体字段的注释
type CommentExtractor struct {
	comments map[string]string
	fset     *token.FileSet
}

// NewCommentExtractor 创建一个新的注释提取器
func NewCommentExtractor() *CommentExtractor {
	return &CommentExtractor{
		comments: make(map[string]string),
		fset:     token.NewFileSet(),
	}
}

// ExtractCommentsFromDir 从目录中提取所有.go文件的注释
func (ce *CommentExtractor) ExtractCommentsFromDir(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 只处理.go文件
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			if err := ce.ExtractComments(path); err != nil {
				return fmt.Errorf("error processing file %s: %v", path, err)
			}
		}
		return nil
	})
}

// ExtractComments 从单个文件中提取注释
func (ce *CommentExtractor) ExtractComments(filename string) error {
	node, err := parser.ParseFile(ce.fset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						structName := typeSpec.Name.Name
						ce.extractStructComments(structName, structType)
					}
				}
			}
		}
	}
	return nil
}

// extractStructComments 提取结构体字段的注释
func (ce *CommentExtractor) extractStructComments(structName string, structType *ast.StructType) {
	for _, field := range structType.Fields.List {
		if field.Comment != nil {
			comment := strings.TrimSpace(strings.TrimPrefix(field.Comment.Text(), "//"))
			for _, name := range field.Names {
				key := fmt.Sprintf("%s.%s", structName, name.Name)
				ce.comments[key] = comment
			}
		}
	}
}

// GetComment 获取指定字段的注释
func (ce *CommentExtractor) GetComment(structName, fieldName string) string {
	return ce.comments[fmt.Sprintf("%s.%s", structName, fieldName)]
}

// SchemaGenerator JSON Schema生成器
type SchemaGenerator struct {
	ce *CommentExtractor
}

// NewSchemaGenerator 创建新的Schema生成器
func NewSchemaGenerator(ce *CommentExtractor) *SchemaGenerator {
	return &SchemaGenerator{ce: ce}
}

// GenerateSchema 生成JSON Schema
func (sg *SchemaGenerator) GenerateSchema(v interface{}) map[string]interface{} {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	schema := map[string]interface{}{
		"type":       "object",
		"properties": make(map[string]interface{}),
	}

	properties := schema["properties"].(map[string]interface{})
	required := []string{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = field.Name
		}
		jsonName := strings.Split(jsonTag, ",")[0]

		if strings.Contains(jsonTag, "required") {
			required = append(required, jsonName)
		}

		fieldSchema := sg.generateFieldSchema(t.Name(), field)
		properties[jsonName] = fieldSchema
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}

func (sg *SchemaGenerator) generateFieldSchema(prefix string, field reflect.StructField) map[string]interface{} {
	schema := make(map[string]interface{})

	comment := sg.ce.GetComment(prefix, field.Name)
	if comment != "" {
		schema["description"] = comment
	}

	switch field.Type.Kind() {
	case reflect.String:
		schema["type"] = "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		schema["type"] = "integer"
	case reflect.Float32, reflect.Float64:
		schema["type"] = "number"
	case reflect.Bool:
		schema["type"] = "boolean"
	case reflect.Slice:
		schema["type"] = "array"
		schema["items"] = sg.generateFieldSchema(field.Name, reflect.StructField{Type: field.Type.Elem()})
	case reflect.Struct:
		nestedSchema := sg.GenerateSchema(reflect.New(field.Type).Interface())
		for k, v := range nestedSchema {
			schema[k] = v
		}
	case reflect.Ptr:
		return sg.generateFieldSchema(field.Name, reflect.StructField{Type: field.Type.Elem()})
	}

	return schema
}
