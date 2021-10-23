package httpd

import (
	_ "embed"
	"fmt"
	"html/template"
	"strings"
)

//go:embed template.html
var templateRaw string

type fileContentsFragment struct {
	Name       string `json:"name"`
	Link       string `json:"link"`
	Icon       string `json:"icon"`
	Annotation string `json:"annotation"`
}

// Used for displaying a directory index.
func renderDirectoryIndexPage(dirName string, files []fileContentsFragment, ReadOnly bool) string {
	tpl, err := template.New("directoryIndexPage").Parse(templateRaw)

	if err != nil {
		return fmt.Sprintf("%v\n", err)
	}

	type TemplateInputDirComponent struct {
		Link string
		Name string
	}
	type TemplateInput struct {
		DirName string
		DirComponents []TemplateInputDirComponent
		Files   []fileContentsFragment
		ReadOnly bool
	}
	var components []TemplateInputDirComponent
	link := ""
	for i, part := range strings.Split(dirName, "/") {
		link = link + part + "/"
		if i == 0 {
			part = "root"
		}
		components = append(components, TemplateInputDirComponent{
			Link: link[:],
			Name: part,
		})
	}

	input := TemplateInput{
		DirName: dirName,
		DirComponents: components,
		Files:   files,
		ReadOnly: ReadOnly,
	}

	builder := strings.Builder{}
	_ = tpl.Execute(&builder, input)

	return builder.String()
}
