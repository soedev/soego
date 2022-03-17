// MIT License
//
// Copyright (c) 2020 go-kratos
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package main

import (
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	errorsPackage = protogen.GoImportPath("github.com/soedev/soego/core/eerrors")
	codesPackage  = protogen.GoImportPath("google.golang.org/grpc/codes")
)

// generateFile generates a _errors.pb.go file containing ego errors definitions.
func generateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Enums) == 0 {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + "_errors.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-go-errors. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	g.QualifiedGoIdent(codesPackage.Ident(""))
	generateFileContent(gen, file, g)
	return g
}

// generateFileContent generates the ego errors definitions, excluding the package statement.
func generateFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) {
	if len(file.Enums) == 0 {
		return
	}

	g.P("// This is a compile-time assertion to ensure that this generated file")
	g.P("// is compatible with the ego package it is being compiled against.")
	g.P("const _ = ", errorsPackage.Ident("SupportPackageIsVersion1"))
	g.P()
	index := 0
	for _, enum := range file.Enums {
		if !generationErrorsSection(gen, file, g, enum) {
			index++
		}
	}
	// If all enums do not contain 'errors.code', the current file is skipped
	if index == 0 {
		g.Skip()
	}
}

const (
	fileLevelCommentAnnotation  = "plugins"
	fieldLevelCommentAnnotation = "code"
)

func generationErrorsSection(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, enum *protogen.Enum) bool {
	var ew errorWrapper
	for _, v := range enum.Values {
		annos := getAnnotations(string(v.Comments.Leading))
		eCode := annos[fieldLevelCommentAnnotation]
		desc := string(v.Desc.Name())
		err := &errorInfo{
			Name:            string(enum.Desc.Name()),
			Value:           desc,
			UpperCamelValue: strcase.ToCamel(strings.ToLower(desc)),
			LowerCamelValue: strcase.ToLowerCamel(strings.ToLower(desc)),
			Code:            strcase.ToCamel(strings.ToLower(eCode.val)),
			Key:             string(v.Desc.FullName()),
		}
		ew.Errors = append(ew.Errors, err)
	}
	if len(ew.Errors) == 0 {
		return true
	}
	g.P(ew.execute())
	return false
}

var filedLevelCommentRgx, _ = regexp.Compile(`@(\w+)=([_a-zA-Z0-9-,]+)`)
var fileLevelCommentRgx, _ = regexp.Compile(`@(\w+)=([_a-zA-Z0-9-,]+)`)

type annotation struct {
	name string
	val  string
}

func getAnnotations(comment string) map[string]annotation {
	matches := filedLevelCommentRgx.FindAllStringSubmatch(comment, -1)
	return findMatchesFromComments(matches)
}

func findMatchesFromComments(matches [][]string) map[string]annotation {
	annotations := make(map[string]annotation)
	for _, v := range matches {
		annotations[v[1]] = annotation{
			name: v[1],
			val:  v[2],
		}
	}
	return annotations
}

func getFileLevelAnnotations(locs []*descriptorpb.SourceCodeInfo_Location) map[string]annotation {
	comments := ""
	for _, loc := range locs {
		comments += loc.String()
	}
	matches := fileLevelCommentRgx.FindAllStringSubmatch(comments, -1)
	return findMatchesFromComments(matches)
}

func needGenerate(locs []*descriptorpb.SourceCodeInfo_Location) bool {
	annos := getFileLevelAnnotations(locs)
	anno, ok := annos[fileLevelCommentAnnotation]
	if !ok {
		return false
	}
	plugins := strings.Split(anno.val, ",")
	for _, p := range plugins {
		// if protobuf contains "@plugins=protoc-gen-go-errors" annotation, then we should generate errors stub code
		if p == "protoc-gen-go-errors" {
			return true
		}
	}
	return false
}
