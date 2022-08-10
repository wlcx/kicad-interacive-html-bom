package main 

import (
	"embed"
	"time"
	"flag"
	"fmt"
	"path/filepath"
	"io"
	"net/url"
	"strings"
	"html/template"
	log "github.com/sirupsen/logrus"
	"os"
	"bytes"
	"github.com/yuin/goldmark"
)

//go:embed main.html.tmpl
var fs embed.FS

var version = "unknown"

func bail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type page struct {
	Content *template.HTML // A map from page name to rendered partial HTML
	Title string
	UrlPath string
	IsFullWidth bool
}

type site struct {
	projectName string
	projectVersion string
	pages map[string]page
}

func newSite(projectName string, projectVersion string) site {
	return site {
		projectName: projectName,
		projectVersion: projectVersion,
		pages: make(map[string]page),
	}
}

func renderMd(raw io.Reader) (*bytes.Buffer, error) {
	md, err := io.ReadAll(raw)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := goldmark.Convert(md, &buf); err !=nil {
		return nil, err
	}

	return &buf, nil
}

// Take a file and output the html we want to appear in the content area of the site
// - Markdown files (*.md) are rendered into HTML
// - HTML (*.html) and PDF (*.pdf) files are linked to in an iframe
func processFile(path string) (template.HTML, bool, error) {
	if _, err := os.Stat(path); err != nil {
		return "", false, fmt.Errorf("stat %s: %w", path, err)
	}
	var rawhtml string
	var shouldCopy bool
	classes := []string{"content"}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md":
		shouldCopy = false
		classes = append(classes, "container")
		f, err := os.Open(path)
		if err != nil {
			return "", false, fmt.Errorf("open md file: %w", err)
		}
		defer f.Close()
		rendered, err := renderMd(f)
		if err != nil {
			return "", false, fmt.Errorf("render %s: %w", path, err)
		}
		rawhtml = rendered.String()
	case ".pdf", ".html":
		shouldCopy = true
		classes = append(classes, "fullsize")
		rawhtml = fmt.Sprintf(`<iframe src="%s" frameborder="0"></iframe>`, filepath.Base(path))
	default:
		return "", false, fmt.Errorf("unsupported extension: %s", path)
	}

	return template.HTML(fmt.Sprintf(`<div class="%s">%s</div>`, strings.Join(classes, " "), rawhtml)), shouldCopy, nil

}

// Add a page with the given html content
func (s *site) AddPage(title string, content *template.HTML) error {
	var urlPath string
	// The first page added is our homepage
	if len(s.pages) == 0 {
		urlPath = "index.html"
	} else {
		urlPath = url.PathEscape(strings.ToLower(title)) + ".html"
	}

	if _, exists := s.pages[urlPath]; exists {
		return fmt.Errorf("Duplicate page path: %s. Use a different title", urlPath)
	}
	s.pages[urlPath] = page{Content: content, Title: title, UrlPath: urlPath}
	return nil
}

type TemplateData struct{
	ProjectName string
	ProjectVersion string
	PageTitle string
	Links []page
	Content *template.HTML
	GeneratedAt string
	SelfVersion string
}

func (s *site) Render(outpath string) error {
	t, err := template.ParseFS(fs, "main.html.tmpl")
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	// A map from page title to url path
	var links ([]page)
	for _, page := range s.pages {
		//TODO sorting
		links = append(links, page)
	}
	for _, page := range s.pages {
		data := TemplateData{
			ProjectName: s.projectName,
			PageTitle: page.Title,
			Links: links,
			Content: page.Content,
			ProjectVersion: s.projectVersion,
			GeneratedAt: time.Now().Format(time.RFC3339),
			SelfVersion: version,
		}

		f, err := os.Create(filepath.Join(outpath, page.UrlPath))
		if err != nil {
			return fmt.Errorf("create file %s: %w", page.UrlPath, err)
		}
		defer f.Close()

		if err := t.Execute(f, &data); err != nil {
			return fmt.Errorf("execute template in page %s: %w", page.Title, err)
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	srcf, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open %s: %w", src, err)
	}
	defer srcf.Close()
	dstf, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("open %s: %w", dst, err)
	}
	defer dstf.Close()
	_, err = io.Copy(dstf, srcf); if err != nil {
		return fmt.Errorf("copy %s to %s: %w", src, dst, err)
	}
	return nil
}

func main() {
	showVersion := flag.Bool("version", false, "Show the version and exit")
	projectName := flag.String("projectName", "KiCAD Project", "The name of the project we are generating a site for. This is shown in the generated HTML.")
	projectVersion := flag.String("projectVersion", "", "The version of the project we are generating a site for. This is shown in the generated HTML.")
	out := flag.String("out", "", "Directory to generate output in")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	s := newSite(*projectName, *projectVersion)
	if len(flag.Args()) == 0 {
		log.Fatal("No files included. Specify file to include at the end of the command invocation")
	}

	var includes []string
	for _, include := range flag.Args() {
		expanded, err := filepath.Glob(include)
		bail(err)
		includes = append(includes, expanded...)
	}
	for _, include := range includes {
		log.Infof("Including file %s", include)
		html, shouldCopy, err := processFile(include)
		if shouldCopy {
			bail(copyFile(include, filepath.Join(*out, filepath.Base(include))))
		}
		bail(err)
		base := filepath.Base(include)
		bail(s.AddPage(base, &html))
	}
	bail(s.Render(*out))
}
