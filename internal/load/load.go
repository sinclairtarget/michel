/*
* Handles reading of files (possibly with YAML frontmatter) from disk.
*
* This package just handles the reading of arbitrary YAML frontmatter without
* specifying a required shape.
*
* Frontmatter appears in a separate block that must be at the beginning of the
* file. The YAML frontmatter is demarcated by a "---" line that appears at the
* beginning and end of the block. Any number of "-" characters can be used as
* long as there are at least three.
 */
package load

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const demarcationChar = '-'

func isDemarcationLine(line string) bool {
	count := 0
	for _, c := range line {
		if c == demarcationChar {
			count += 1
		} else {
			return false
		}
	}

	return count >= 3
}

// Result of loading a file from disk.
type Result[TFrontmatter any] struct {
	Frontmatter TFrontmatter // Loaded frontmatter if there was any
	Text        string       // Main text from file
}

type Opts struct {
	FrontmatterOnly bool
}

func ReadFile[TFrontmatter any](
	path string,
	opts Opts,
) (result Result[TFrontmatter], err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error reading file \"%s\": %w", path, err)
		}
	}()

	f, err := os.Open(path)
	if err != nil {
		return result, err
	}
	defer f.Close()

	var (
		yamlBuilder             strings.Builder
		textBuilder             strings.Builder
		lineIndex               int
		numDemarcationLinesSeen int
	)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		demarcationAllowed := lineIndex == 0 || numDemarcationLinesSeen == 1
		if demarcationAllowed && isDemarcationLine(line) {
			numDemarcationLinesSeen += 1

			if numDemarcationLinesSeen == 2 && opts.FrontmatterOnly {
				break
			}

			lineIndex += 1
			continue
		}

		if numDemarcationLinesSeen == 1 { // In YAML block
			_, err := yamlBuilder.WriteString(line)
			if err != nil {
				return result, err
			}

			err = yamlBuilder.WriteByte('\n')
			if err != nil {
				return result, err
			}
		} else { // Passed YAML block
			_, err := textBuilder.WriteString(line)
			if err != nil {
				return result, err
			}

			err = textBuilder.WriteByte('\n')
			if err != nil {
				return result, err
			}
		}

		lineIndex += 1
	}

	if yamlBuilder.Len() > 0 {
		err := yaml.Unmarshal([]byte(yamlBuilder.String()), &result.Frontmatter)
		if err != nil {
			return result, err
		}
	}

	result.Text = textBuilder.String()
	return result, nil
}
