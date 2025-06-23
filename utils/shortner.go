package utils

type CodeGenerator interface {
	Generate(input string) string
}
