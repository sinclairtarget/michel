package build

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/sinclairtarget/michel/internal/content"
	"github.com/sinclairtarget/michel/internal/site"
)

type Options struct {
	SiteDir     string
	TargetDir   string
	ShouldClean bool
}

// Build the static website.
//
//  1. Clean target dir.
//  2. Load site configuration.
//  3. Load content.
//  4. Load partials, prefixed with "partials/"
//  5. For each site path:
//     a. Read YAML frontmatter
//     b. Load layouts (if any), prefixed with partials/
//     c. Load template
//     d. ExecuteTemplate() with first layout
func Build(logger *slog.Logger, options Options) error {
	if options.ShouldClean {
		logger.Debug("cleaning target directory")
		err := clean(options.TargetDir)
		if err != nil {
			return fmt.Errorf("failed to clean target directory: %v", err)
		}
	}

	logger.Debug("loading site")
	siteMetadata := site.Load(options.SiteDir)

	logger.Debug("loading partials templates")
	tmpl, err := loadPartials("partials")
	if err != nil {
		return fmt.Errorf("failed to load partials templates: %w", err)
	}

	page, err := site.LoadPage("site/two-houses-in-cambridgeport.html.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load site page: %w", err)
	}

	executedTmplName := "two-houses-in-cambridgeport.html.tmpl"
	if len(page.Frontmatter.Layouts) > 0 {
		tmpl, err = tmpl.ParseFiles(page.Frontmatter.Layouts...)
		if err != nil {
			return fmt.Errorf("failed to parse layout template: %w", err)
		}

		executedTmplName = filepath.Base(page.Frontmatter.Layouts[0])
	}

	tmpl, err = tmpl.New("two-houses-in-cambridgeport.html.tmpl").Parse(page.TemplateText)
	if err != nil {
		return fmt.Errorf("failed to parse article template: %w", err)
	}

	logger.Debug("defined templates", "templates", tmpl.DefinedTemplates())

	f, err := os.Create("public/two-houses-in-cambridgeport.html")
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	data := struct {
		SiteName string
		Content  content.Content
	}{
		SiteName: siteMetadata.Config.Name,
		Content: content.Content{
			Title:    "Two Houses I Like In Cambridgeport",
			BodyText: "This is my article about houses in Cambridgeport",
		},
	}
	err = tmpl.ExecuteTemplate(f, executedTmplName, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// 	logger.Debug("processing site pages")
	// 	seq, finish := site.Paths()
	// 	for sitePath := range seq {
	// 		targetPath, err := target(options.SiteDir, options.TargetDir, sitePath)
	// 		if err != nil {
	// 			return fmt.Errorf("could not map path: %v", err)
	// 		}
	//
	// 		err = process(sitePath, targetPath)
	// 		if err != nil {
	// 			return fmt.Errorf("failed to process \"%s\": %v", sitePath, err)
	// 		}
	// 	}
	//
	// 	err := finish()
	// 	if err != nil {
	// 		return err
	// 	}

	return nil
}

func clean(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}

	err = os.Mkdir(dir, 0o755)
	if err != nil {
		return err
	}

	return nil
}

// Returns output path under target dir given path in source dir.
func target(siteDir string, targetDir string, path string) (string, error) {
	relative, err := filepath.Rel(siteDir, path)
	if err != nil {
		return "", err
	}

	return filepath.Join(targetDir, relative), nil
}

func process(sourcePath string, targetPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	target, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = io.Copy(target, source)
	return err
}
