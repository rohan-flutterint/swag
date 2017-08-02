package gen

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/swaggo/swag"
)

type Gen struct {
}

func New() *Gen {
	return &Gen{}
}

func (g *Gen) Build(searchDir, mainApiFile string) error {
	log.Println("Generate swagger docs....")
	p := swag.New()
	p.ParseApi(searchDir, mainApiFile)
	swagger := p.GetSwagger()

	b, _ := json.MarshalIndent(swagger, "", "    ")

	os.MkdirAll(path.Join(searchDir, "docs"), os.ModePerm)
	docs, _ := os.Create(path.Join(searchDir, "docs", "docs.go"))
	defer docs.Close()

	packageTemplate.Execute(docs, struct {
		Timestamp time.Time
		Doc       string
	}{
		Timestamp: time.Now(),
		Doc:       "`" + string(b) + "`",
	})

	log.Printf("create docs.go at  %+v", docs.Name())
	return nil
}

var packageTemplate = template.Must(template.New("").Parse(`// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// {{ .Timestamp }}

package docs

import (
	"github.com/swaggo/swag/swagger"
)

var doc = {{.Doc}}

type s struct{}

func (s *s) ReadDoc() string {
	return doc
}
func init() {
	swagger.Register(swagger.Name, &s{})
}
`))
