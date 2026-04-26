package rsa

import (
	"encoding/pem"
	"os"
	"strings"
)

func ReadKey(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content = []byte(strings.TrimSpace(string(content)))

	block, _ := pem.Decode(content)
	if block == nil {
		return nil, os.ErrInvalid // 明确返回错误：不是合法 PEM 格式
	}

	return block.Bytes, nil
}
