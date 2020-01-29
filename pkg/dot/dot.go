package dot

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

func ToSVG(dot string) (string, error) {
	cmd := exec.Command("dot", "-Tsvg")
	cmd.Stdin = bytes.NewReader([]byte(dot))

	r, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	err = cmd.Start()
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(r)
	out := ""
	for {
		chunk, err := reader.ReadString('\n')
		if chunk != "" {
			out = fmt.Sprintf("%s%s", out, chunk)
		}

		if err != nil {
			if err == io.EOF {
				return out, nil
			}
			return "", err
		}
	}
}
