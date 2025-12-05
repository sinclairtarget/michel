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

	"github.com/sinclairtarget/michel/internal/site"
)

var siteDir string = "site"
var targetDir string = "public"
var contentDir string = "content"
var layoutsDir string = "layouts"
var partialsDir string = "partials"

func Build(logger *slog.Logger) error {
	start := time.Now()
	logger.Debug("beginning build")

	logger.Debug("cleaning target directory")
	err := clean(targetDir)
	if err != nil {
		return fmt.Errorf("failed to clean target directory: %v", err)
	}

	logger.Debug("loading site")
	siteMetadata := site.Load(siteDir)

	logger.Debug("loading partials templates")
	tmpl, err := loadPartials(partialsDir)
	if err != nil {
		return fmt.Errorf("failed to load partials templates: %w", err)
	}

	logger.Debug("loading content")
	contentLibrary, err := loadContent(contentDir)
	data := struct {
		Site    site.Site
		Content ContentLibrary
	}{
		Site:    siteMetadata,
		Content: contentLibrary,
	}

	logger.Debug("processing site pages and assets")
	seq, finish := siteMetadata.Paths()
	for sitePath := range seq {
		if site.IsPage(sitePath) {
			logger.Debug("processing page", "path", sitePath)
			targetPath, err := mapPagePath(
				sitePath,
				siteDir,
				targetDir,
			)
			if err != nil {
				return fmt.Errorf("could not map path: %v", err)
			}

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
			targetPath, err := mapAssetPath(
				sitePath,
				siteDir,
				targetDir,
			)
			if err != nil {
				return fmt.Errorf("could not map path: %v", err)
			}

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

	layouts := page.Frontmatter.LayoutsFullName()
	if len(layouts) > 0 {
		var layoutPaths []string
		for _, layoutName := range layouts {
			path, err := layoutPathFromName(layoutName, layoutsDir)
			if err != nil {
				return err
			}

			layoutPaths = append(layoutPaths, path)
		}

		tmpl, err = loadLayouts(layoutsDir, layoutPaths, partialsTmpl)
		if err != nil {
			return fmt.Errorf("failed to load layouts: %w", err)
		}
		execName = layouts[0]
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
