//go:build ignore

package main

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/valyala/fasthttp"
)

// -----------------------------------------------------------------------------

type TemplateConfig struct {
	Methods []Method
}

type Method struct {
	Name     string
	InTypes  []string
	OutTypes []string
}

// -----------------------------------------------------------------------------

func main() {
	funcs := template.FuncMap{
		"processParams": processParams,
	}

	// Generator template
	tmpl := template.Must(template.New("").Funcs(funcs).Parse(`package go_webserver

// Code generated by go generate; DO NOT EDIT.

import (
	"crypto/tls"
	"io"
	"mime/multipart"
	"net"
	"time"

	"github.com/valyala/fasthttp"
)

// -----------------------------------------------------------------------------

{{range $method := .Methods}}
{{-  $outTypesLen := len $method.OutTypes }}
func (req *RequestContext) {{$method.Name}}({{processParams $method.InTypes "in" true}}) ` +
		`({{processParams $method.OutTypes "out" true}}) {
	{{- if gt $outTypesLen 0 }}
		{{processParams $method.OutTypes "out" false}} = req.ctx.{{$method.Name}}({{processParams $method.InTypes "in" false}})
	{{- else }}
		req.ctx.{{$method.Name}}({{processParams $method.InTypes "in" false}})
	{{- end }}
	return
}
{{end}}
`))

	tmplConfig := TemplateConfig{
		Methods: make([]Method, 0),
	}

	// Process 'fasthttp.RequestCtx' public methods
	reqCtx := &fasthttp.RequestCtx{}
	v := reflect.ValueOf(reqCtx)
	t := reflect.TypeOf(reqCtx)

	// Iterate through all public methods
	for methodIdx := 0; methodIdx < v.NumMethod(); methodIdx++ {
		method := t.Method(methodIdx)
		methodType := method.Type

		// Skip some public methods
		switch method.Name {
		case "NotFound":
		case "NotModified":
		case "Success":
		case "Error":

			// These are replaced with versions that takes into account proxies
		case "Host":
		case "RemoteIP":

		default:
			// Add this method
			tmplMethod := Method{
				Name:     method.Name,
				InTypes:  make([]string, 0),
				OutTypes: make([]string, 0),
			}

			// We are assuming the first parameter is the pointer to 'fasthttp.RequestCtx'
			if methodType.NumIn() >= 2 {
				for idx := 1; idx < methodType.NumIn(); idx++ {
					tmplMethod.InTypes = append(tmplMethod.InTypes, methodType.In(idx).String())
				}
			}

			for idx := 0; idx < methodType.NumOut(); idx++ {
				tmplMethod.OutTypes = append(tmplMethod.OutTypes, methodType.Out(idx).String())
			}

			tmplConfig.Methods = append(tmplConfig.Methods, tmplMethod)
		}
	}

	// Generate output
	output := &bytes.Buffer{}
	err := tmpl.Execute(output, tmplConfig)
	if err != nil {
		log.Fatalf("Error executing template [err=%v]", err)
	}

	var data []byte

	data, err = format.Source(output.Bytes())
	if err != nil {
		log.Fatalf("Error formatting generated code [err=%v]", err)
	}

	err = ioutil.WriteFile("request_generated.go", data, os.ModePerm)
	if err != nil {
		log.Fatalf("Error writing generated file [err=%v]", err)
	}
}

// -----------------------------------------------------------------------------
// Helpers

func processParams(items []string, namePrefix string, addTypes bool) string {
	hasNamePrefix := len(namePrefix) > 0

	result := make([]string, 0)
	for idx := range items {
		s := ""

		if hasNamePrefix {
			s = namePrefix + strconv.Itoa(idx+1)
			if addTypes {
				s = s + " "
			}
		}

		if addTypes {
			s = s + items[idx]
		}

		result = append(result, s)
	}

	return strings.Join(result, ", ")
}
