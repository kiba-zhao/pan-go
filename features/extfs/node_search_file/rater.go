package nodesearchfile

import (
	"bufio"
	"mime"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"
)

type FileRater interface {
	// Rate returns the score of the file
	Rate(filePath string, tokens []string) (uint, error)
	// Tokenize returns the tokens of the text
	Tokenize(text string) []string
}

type stdFileRater struct {
}

var _ = (FileRater)((*stdFileRater)(nil))

// Rate computes the number of tokens that are found within the filename
// specified by filePath. It returns the count of matched tokens as a uint.
// An error is returned if there is any issue during the computation.

func (fr *stdFileRater) Rate(filePath string, tokens []string) (uint, error) {
	matchedCount := 0

	_, filename := path.Split(filePath)
	for _, token := range tokens {
		if strings.Contains(filename, token) {
			matchedCount++
		}
	}

	return uint(matchedCount), nil
}

// Tokenize splits the given text into tokens. The function first trims
// leading and trailing whitespace, and then splits on one or more
// whitespace characters. The result is a slice of the tokens.
func (fr *stdFileRater) Tokenize(text string) []string {
	text = strings.Trim(text, " ")
	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	return strings.Split(text, " ")
}

func rateWithFilePath(filePath string, tokens []string) uint {
	matchedCount := 0

	_, filename := path.Split(filePath)
	for _, token := range tokens {
		if strings.Contains(filename, token) {
			matchedCount++
		}
	}

	return uint(matchedCount)
}

func rateWithMimeType(filePath string, tokens []string, mimeType string) (uint, error) {
	exts, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return 0, err
	}

	if len(exts) == 0 {
		return 0, err
	}

	if slices.Contains(exts, ".txt") {
		return rateWithTextFile(filePath, tokens)
	}

	return 0, err
}

func rateWithTextFile(filePath string, tokens []string) (uint, error) {

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	matchedIdxes := make([]int, 0)
	for scanner.Scan() {
		text := scanner.Text()
		for idx, token := range tokens {
			offset, ok := slices.BinarySearch(matchedIdxes, idx)
			if ok {
				continue
			}
			if strings.Contains(text, token) {
				matchedIdxes = slices.Insert(matchedIdxes, offset, idx)
			}
		}
		if len(matchedIdxes) >= len(tokens) {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return uint(len(matchedIdxes)), nil
}
