/*
* Package responsible for coordinating the build.
*
* To build the site, we:
* 	1. Clean the target dir.
* 	2. Load site metadata.
* 	3. Load content.
* 	4. Load partials, prefixed with "partials/"
* 	5. For each site path:
* 		 If it is a page (*.html, *.html.tmpl):
* 	       a. Read YAML frontmatter
* 	       b. Load layouts defined in frontmatter, prefixed with layouts/
* 	       c. Load template
* 	       d. ExecuteTemplate() with first layout
* 	     Otherwise, it is an asset:
* 	       Copy it to the target dir
 */
package build

import (
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/sinclairtarget/michel/internal/content"
	"github.com/sinclairtarget/michel/internal/site"
)

type Options struct {
	SiteDir   string
	TargetDir string
}

func Build(logger *slog.Logger, options Options) error {
	start := time.Now()
	logger.Debug("beginning build")

	logger.Debug("cleaning target directory")
	err := clean(options.TargetDir)
	if err != nil {
		return fmt.Errorf("failed to clean target directory: %v", err)
	}

	logger.Debug("loading site")
	siteMetadata := site.Load(options.SiteDir)

	logger.Debug("loading partials templates")
	tmpl, err := loadPartials("partials")
	if err != nil {
		return fmt.Errorf("failed to load partials templates: %w", err)
	}

	logger.Debug("loading content")
	data := struct {
		SiteName string
		Content  content.Content
	}{
		SiteName: siteMetadata.Config.Name,
		Content: content.Content{
			Path: "content/two-houses-in-cambridgeport.txt",
			Frontmatter: content.Frontmatter{
				Title: "Two Houses in Cambridgeport",
			},
			Html: "<p>Foo bar</p>",
		},
	}

	logger.Debug("processing site pages and assets")
	seq, finish := siteMetadata.Paths()
	for sitePath := range seq {
		targetPath, err := target(options.SiteDir, options.TargetDir, sitePath)
		if err != nil {
			return fmt.Errorf("could not map path: %v", err)
		}

		if site.IsPage(sitePath) {
			logger.Debug("processing page", "path", sitePath)
			err = processPage(sitePath, targetPath, tmpl, data)
			if err != nil {
				return fmt.Errorf(
					"failed to process page \"%s\": %v",
					sitePath,
					err,
				)
			}
		} else {
			logger.Debug("processing asset", "path", sitePath)
			err = processAsset(sitePath, targetPath)
			if err != nil {
				return fmt.Errorf(
					"failed to process asset \"%s\": %v",
					sitePath,
					err,
				)
			}
		}
	}

	err = finish()
	if err != nil {
		return err
	}

	elapsed := time.Now().Sub(start)
	logger.Debug("build complete", "durationMs", elapsed.Milliseconds())
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

// TODO: Layout dir should be configurable
// TODO: Don't hardcode extension
func mapLayoutNameToPath(name string) string {
	return fmt.Sprintf("layouts/%s.html.tmpl", name)
}

func processPage(
	sourcePath string,
	targetPath string,
	partialsTmpl *template.Template,
	data any,
) error {
	page, err := site.LoadPage(sourcePath)
	if err != nil {
		return fmt.Errorf(
			"failed to load site page \"%s\": %w",
			sourcePath,
			err,
		)
	}

	tmpl := partialsTmpl
	tmplName := filepath.Base(sourcePath)
	execName := tmplName

	if len(page.Frontmatter.Layouts) > 0 {
		var layoutPaths []string
		for _, layoutName := range page.Frontmatter.Layouts {
			path := mapLayoutNameToPath(layoutName)
			layoutPaths = append(layoutPaths, path)
		}

		tmpl, err = loadLayouts(partialsTmpl, layoutPaths)
		if err != nil {
			return fmt.Errorf("failed to load layouts: %w", err)
		}
		execName = "layouts/" + page.Frontmatter.Layouts[0]
	}

	if tmpl != nil {
		tmpl, err = tmpl.New(tmplName).Parse(page.TemplateText)
	} else {
		tmpl, err = template.New(tmplName).Parse(page.TemplateText)
	}
	if err != nil {
		return fmt.Errorf(
			"failed to parse site template \"%s\": %w",
			sourcePath,
			err,
		)
	}

	f, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf(
			"failed to create file at \"%s\": %w",
			targetPath,
			err,
		)
	}
	defer f.Close()

	err = tmpl.ExecuteTemplate(f, execName, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func processAsset(sourcePath string, targetPath string) error {
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
