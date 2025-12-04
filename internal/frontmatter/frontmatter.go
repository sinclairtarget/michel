/*
* Handles reading of files that might have YAML frontmatter.
*
* Frontmatter appears in a separate block that must be at the beginning of the
* file. The YAML frontmatter is demarcated by a "---" line that appears at the
* beginning and end of the block. Any number of "-" characters can be used as
* long as there are at least three.
 */
package frontmatter

import (
	"bufio"
	"io"
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

type Result[T any] struct {
	Frontmatter T      // Loaded frontmatter if there was any
	Text        string // Main text from file
}

func ReadFile[T any](r io.Reader) (Result[T], error) {
	var (
		result                  Result[T]
		yamlBuilder             strings.Builder
		textBuilder             strings.Builder
		numDemarcationLinesSeen int
	)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if numDemarcationLinesSeen < 2 && isDemarcationLine(line) {
			numDemarcationLinesSeen += 1
			continue
		}

		if numDemarcationLinesSeen == 1 {
			_, err := yamlBuilder.WriteString(line)
			if err != nil {
				return result, err
			}

			err = yamlBuilder.WriteByte('\n')
			if err != nil {
				return result, err
			}
		} else {
			_, err := textBuilder.WriteString(line)
			if err != nil {
				return result, err
			}

			err = textBuilder.WriteByte('\n')
			if err != nil {
				return result, err
			}
		}
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
