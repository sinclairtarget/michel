/*
* Package build is responsible for coordinating the build.
*
* A build of the site proceeds as follows:
* 	1. Load the Michel config.
* 	2. Load site page and asset metadata. If there is none, quit here.
* 	3. Clean the target dir.
* 	4. Load content metadata.
* 	5. Load layouts.
* 	6. Load partials.
* 	7. For each site page:
* 	       a. Load page template
* 	       b. Parse it
* 	       c. ExecuteTemplate() with layouts defined in the page frontmatter
* 	8. For each site asset:
* 	     Copy it to the target dir
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
	"github.com/sinclairtarget/michel/internal/site"
)

// Input directories
const (
	ContentDir  string = "content"
	SiteDir            = "site"
	LayoutsDir         = "layouts"
	PartialsDir        = "partials"
)

// Output directory
const TargetDir string = "public"

// Scope for a build.
//
// This is the relevant universe of inputs to a build.
type scope struct {
	config   config.Config
	site     site.Site
	corpus   content.Corpus
	layouts  []Layout
	partials []Partial
	start    time.Time
}

func Build(logger *slog.Logger) error {
	var (
		scope scope
		err   error
	)

	logger.Debug("beginning build")
	scope.start = time.Now()

	logger.Debug("loading config")
	scope.config, err = config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	logger.Debug("loading site metadata")
	scope.site, err = site.LoadSite(SiteDir)
	if err != nil {
		return fmt.Errorf("failed to load site metadata: %v", err)
	}

	if scope.site.NumPages()+scope.site.NumAssets() == 0 {
		logger.Debug("build done because site is empty")
		return nil
	}

	logger.Debug("cleaning target directory")
	err = clean(TargetDir)
	if err != nil {
		return fmt.Errorf("failed to clean target directory: %v", err)
	}

	logger.Debug("loading content metadata")
	scope.corpus, err = content.LoadCorpus(ContentDir)
	if err != nil {
		return fmt.Errorf("failed to load content metadata: %v", err)
	}

	logger.Debug("loading layouts")
	scope.layouts, err = loadLayouts(LayoutsDir)
	if err != nil {
		return fmt.Errorf("failed to load layouts: %w", err)
	}

	logger.Debug("loading partials")
	scope.partials, err = loadPartials(PartialsDir)
	if err != nil {
		return fmt.Errorf("failed to load partials: %w", err)
	}

	rootTmpl := template.New("root")
	rootTmpl, err = addPartials(rootTmpl, scope.partials)
	if err != nil {
		return fmt.Errorf("failed to parse partials: %w", err)
	}

	logger.Debug("processing pages")
	for page := range scope.site.Pages() {
		targetPath := mapPagePath(page.Path, SiteDir, TargetDir)
		logger.Debug(
			"processing page",
			"path",
			page.Path,
			"targetPath",
			targetPath,
		)
		err = processPage(
			page,
			targetPath,
			scope,
			template.Must(rootTmpl.Clone()),
		)
		if err != nil {
			return fmt.Errorf(
				"failed to process page \"%s\": %w",
				page.Path,
				err,
			)
		}
	}

	logger.Debug("processing assets")
	for asset := range scope.site.Assets() {
		targetPath := mapAssetPath(asset.Path, SiteDir, TargetDir)
		logger.Debug(
			"processing asset",
			"path",
			asset.Path,
			"targetPath",
			targetPath,
		)
		err = processAsset(asset, targetPath)
		if err != nil {
			return fmt.Errorf(
				"failed to process asset \"%s\": %w",
				asset.Path,
				err,
			)
		}
	}

	elapsed := time.Now().Sub(scope.start)
	logger.Debug(
		"build complete",
		"durationMs",
		elapsed.Milliseconds(),
		"pages",
		scope.site.NumPages(),
		"assets",
		scope.site.NumAssets(),
	)
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
	metadata site.PageMetadata,
	targetPath string,
	scope scope,
	rootTmpl *template.Template,
) error {
	// Set up output file
	err := os.MkdirAll(filepath.Dir(targetPath), 0o755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fout, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf(
			"failed to create file at \"%s\": %w",
			targetPath,
			err,
		)
	}
	defer fout.Close()

	// Add layouts
	layoutKeys := metadata.Layouts
	tmpl, err := addLayouts(rootTmpl, scope.layouts, layoutKeys)
	if err != nil {
		return err // TODO: Handle layout not found
	}

	// Parse page template
	tmplName := filepath.Base(metadata.Path)
	tmpl = tmpl.New(tmplName)

	dot := Dot{
		Config:  &scope.config,
		Content: &scope.corpus,
		Page:    &metadata,
		Now:     scope.start,
	}
	tmpl.Funcs(dot.FuncMap(tmpl, fout))

	page, err := metadata.LoadPage()
	if err != nil {
		return err
	}

	tmpl, err = tmpl.Parse(page.TemplateText)
	if err != nil {
		return fmt.Errorf(
			"failed to parse template \"%s\": %w",
			page.Path,
			err,
		)
	}

	// Execute template and write output
	var execName string
	if len(layoutKeys) > 0 {
		// If we have layouts, we should start executing with the first one
		execName = templateName("layouts", layoutKeys[0])
	} else {
		// No layouts? Just execute the page template
		execName = tmplName
	}

	tmpl = template.Must(tmpl.Clone())
	err = tmpl.ExecuteTemplate(fout, execName, dot)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// Just copies file to output directory unmodified.
func processAsset(asset site.AssetMetadata, targetPath string) error {
	source, err := os.Open(asset.Path)
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
