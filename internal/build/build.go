/*
* Package build is responsible for coordinating the build.
*
* To build the site, we:
* 	1. Load the Michel config.
* 	2. Clean the target dir.
* 	3. Load content.
* 	4. Load layouts.
* 	5. Load partials.
* 	6. For each page path:
* 		 If it is a page (*.html, *.html.tmpl):
* 	       a. Read YAML frontmatter
* 	       b. Use layouts defined in frontmatter
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
	"github.com/sinclairtarget/michel/internal/page"
	"github.com/sinclairtarget/michel/internal/util"
)

// Input directories
const (
	ContentDir  string = "content"
	PagesDir           = "site"
	LayoutsDir         = "layouts"
	PartialsDir        = "partials"
)

// Output directory
const TargetDir string = "public"

func Build(logger *slog.Logger) error {
	start := time.Now()
	logger.Debug("beginning build")

	logger.Debug("loading config")
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	logger.Debug("cleaning target directory")
	err = clean(TargetDir)
	if err != nil {
		return fmt.Errorf("failed to clean target directory: %v", err)
	}

	logger.Debug("loading content")
	collection, err := content.LoadCollection(ContentDir)

	logger.Debug("loading layouts")
	layouts, err := page.LoadLayouts(LayoutsDir)
	if err != nil {
		return fmt.Errorf("failed to load layouts: %w", err)
	}

	logger.Debug("loading partials")
	partials, err := page.LoadPartials(PartialsDir)
	if err != nil {
		return fmt.Errorf("failed to load partials: %w", err)
	}

	partialsTmpl := template.New("root")
	partialsTmpl, err = page.AddPartials(partialsTmpl, partials)
	if err != nil {
		return fmt.Errorf("failed to parse partials: %w", err)
	}

	logger.Debug("processing pages and assets")
	seq, finish := util.WalkFiles(PagesDir)
	for path := range seq {
		if page.IsPage(path) {
			targetPath := mapPagePath(path, PagesDir, TargetDir)
			logger.Debug(
				"processing page",
				"path",
				path,
				"targetPath",
				targetPath,
			)
			err = processPage(
				path,
				targetPath,
				cfg,
				collection,
				layouts,
				template.Must(partialsTmpl.Clone()),
				start,
			)
			if err != nil {
				return fmt.Errorf(
					"failed to process page \"%s\": %w",
					path,
					err,
				)
			}
		} else {
			targetPath := mapAssetPath(path, PagesDir, TargetDir)
			logger.Debug(
				"processing asset",
				"path",
				path,
				"targetPath",
				targetPath,
			)
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

	return nil
}

func processPage(
	sourcePath string,
	targetPath string,
	cfg config.Config,
	collection content.Collection,
	layouts []page.Layout,
	partialsTmpl *template.Template,
	now time.Time,
) error {
	p, err := page.LoadPage(PagesDir, sourcePath)
	if err != nil {
		return fmt.Errorf(
			"failed to load page \"%s\": %w",
			sourcePath,
			err,
		)
	}

	// Set up output file
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

	// Add layouts
	layoutKeys := p.Frontmatter.Layouts
	tmpl, err := page.AddLayouts(partialsTmpl, layouts, layoutKeys)
	if err != nil {
		return err // TODO: Handle layout not found
	}

	// Parse page template
	tmplName := filepath.Base(sourcePath)
	tmpl = tmpl.New(tmplName)

	dot := page.Dot{
		Config:  &cfg,
		Content: &collection,
		Now:     now,
	}
	tmpl.Funcs(dot.FuncMap(tmpl, f))

	tmpl, err = tmpl.Parse(p.TemplateText)
	if err != nil {
		return fmt.Errorf(
			"failed to parse template \"%s\": %w",
			sourcePath,
			err,
		)
	}

	// Execute template and write output
	var execName string
	if len(layoutKeys) > 0 {
		// If we have layouts, we should start executing with the first one
		execName = page.TemplateName("layouts", layoutKeys[0])
	} else {
		// No layouts? Just execute the page template
		execName = tmplName
	}

	tmpl = template.Must(tmpl.Clone())
	err = tmpl.ExecuteTemplate(f, execName, dot)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// Just copies file to output directory unmodified.
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
