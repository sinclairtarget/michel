/*
* Package responsible for coordinating the build.
*
* To build the site, we:
* 	1. Load the Michel config.
* 	2. Clean the target dir.
* 	3. Load content.
* 	4. Load partials, prefixed with "partials/"
* 	5. For each page path:
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

	"github.com/sinclairtarget/michel/internal/config"
	"github.com/sinclairtarget/michel/internal/content"
	"github.com/sinclairtarget/michel/internal/content/myst"
	"github.com/sinclairtarget/michel/internal/page"
	"github.com/sinclairtarget/michel/internal/util"
)

// Michel config
const ConfigFilename string = "michel.yaml"

// Input directories
const (
	ContentDir string = "content"
	PagesDir = "site"
	LayoutsDir = "layouts"
	PartialsDir = "partials"
)

// Output directory
const TargetDir string = "public"

func Build(logger *slog.Logger) error {
	start := time.Now()
	logger.Debug("beginning build")

	logger.Debug("loading config")
	cfg := config.Load(ConfigFilename)

	logger.Debug("cleaning target directory")
	err := clean(TargetDir)
	if err != nil {
		return fmt.Errorf("failed to clean target directory: %v", err)
	}

	logger.Debug("loading content")
	contentCollection, err := content.LoadAllContent(ContentDir)
	data := struct {
		Config  config.Config
		Content content.Collection
	}{
		Config:  cfg,
		Content: contentCollection,
	}

	logger.Debug("loading partials")
	tmpl, err := loadPartials(PartialsDir)
	if err != nil {
		return fmt.Errorf("failed to load partials templates: %w", err)
	}

	logger.Debug("processing pages and assets")
	seq, finish := util.WalkPaths(PagesDir)
	for path := range seq {
		if page.IsPage(path) {
			logger.Debug("processing page", "path", path)
			targetPath, err := mapPagePath(
				path,
				PagesDir,
				TargetDir,
			)
			if err != nil {
				return fmt.Errorf("could not map path: %w", err)
			}

			tmpl = template.Must(tmpl.Clone())
			err = processPage(path, targetPath, tmpl, data)
			if err != nil {
				return fmt.Errorf(
					"failed to process page \"%s\": %w",
					path,
					err,
				)
			}
		} else {
			logger.Debug("processing asset", "path", path)
			targetPath, err := mapAssetPath(
				path,
				PagesDir,
				TargetDir,
			)
			if err != nil {
				return fmt.Errorf("could not map path: %w", err)
			}

			err = processAsset(path, targetPath)
			if err != nil {
				return fmt.Errorf(
					"failed to process asset \"%s\": %w",
					path,
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
	page, err := page.LoadPage(sourcePath)
	if err != nil {
		return fmt.Errorf(
			"failed to load page \"%s\": %w",
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
		for _, layoutKey := range layouts {
			path, err := layoutPathFromKey(layoutKey, LayoutsDir)
			if err != nil {
				return err
			}

			layoutPaths = append(layoutPaths, path)
		}

		tmpl, err = loadLayouts(LayoutsDir, layoutPaths, partialsTmpl)
		if err != nil {
			return fmt.Errorf("failed to load layouts: %w", err)
		}
		execName = layouts[0]
	}

	tmpl = tmpl.New(tmplName)
	funcMap := template.FuncMap{
		"html": myst.RenderHTML,
	}
	tmpl.Funcs(funcMap)

	tmpl, err = tmpl.Parse(page.TemplateText)
	if err != nil {
		return fmt.Errorf(
			"failed to parse template \"%s\": %w",
			sourcePath,
			err,
		)
	}

	err = os.MkdirAll(filepath.Dir(targetPath), 0o755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
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

	tmpl = template.Must(tmpl.Clone())
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

	err = os.MkdirAll(filepath.Dir(targetPath), 0o755)
	if err != nil {
		return err
	}

	target, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = io.Copy(target, source)
	return err
}
