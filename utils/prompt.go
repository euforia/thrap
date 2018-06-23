package utils

import (
	"bufio"
	"io"
)

func PromptUntilNoError(prompt string, out io.Writer, in io.Reader, f func([]byte) error) {
	var (
		lb []byte
	)
	err := io.ErrUnexpectedEOF
	for err != nil {
		out.Write([]byte(prompt))

		rd := bufio.NewReader(in)
		lb, _ = rd.ReadBytes('\n')
		lb = lb[:len(lb)-1]
		err = f(lb)
	}
}
