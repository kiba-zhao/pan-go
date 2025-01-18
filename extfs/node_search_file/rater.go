package nodesearchfile

import (
	"bufio"
	"mime"
	"os"
	nodeitem "pan/extfs/node_item"
	"path"
	"regexp"
	"slices"
	"strings"
)

type FileRater interface {
	Rate(filePath string, tokens []string) (uint, error)
	Tokenize(text string) []string
}

type fileRaterImpl struct {
}

func (fr *fileRaterImpl) Rate(filePath string, tokens []string) (uint, error) {

	score := uint(0)
	stat, err := os.Stat(filePath)
	if err != nil {
		return score, err
	}

	if !stat.IsDir() {
		mimeType, err := nodeitem.GenerateMimeType(filePath)
		if err != nil {
			return score, err
		}

		mimeTypeScore, err := rateWithMimeType(filePath, tokens, mimeType)
		if err != nil {
			return score, err
		}
		score += mimeTypeScore
	}

	score += rateWithFilePath(filePath, tokens)
	return score, nil
}

func (fr *fileRaterImpl) Tokenize(text string) []string {
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
