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
		DirName       string
		RelativeRootDir  string
		DirComponents []TemplateInputDirComponent
		Files         []fileContentsFragment
		ReadOnly      bool
	}
	var components []TemplateInputDirComponent
	link := ""
	splitDirName := strings.Split(dirName, "/")
	relativeRootDir := "."
	if len(splitDirName) > 2 {
		relativeRootDir = relativeRootDir + strings.Repeat("/..", len(splitDirName) - 2)
	}
	for i, part := range splitDirName {
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
		DirName:       dirName,
		RelativeRootDir: relativeRootDir,
		DirComponents: components,
		Files:         files,
		ReadOnly:      ReadOnly,
	}

	builder := strings.Builder{}
	_ = tpl.Execute(&builder, input)

	return builder.String()
}
