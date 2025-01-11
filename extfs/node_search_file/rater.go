package nodesearchfile

type FileRater interface {
	Rate(filePath string, tokens []string) (int8, error)
	Tokenize(text string) []string
}
