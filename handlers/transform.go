package handlers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/russross/blackfriday/v2"
)

type Payload struct {
	Content string      `json:"content"`
	Data    interface{} `json:"data"`
}

type Transform (func(Payload) (Payload, error))

func ToHtml(input Payload) (Payload, error) {
	text, ok := input.Data.(string)
	if !ok {
		return input, fmt.Errorf("cannot convert '%T' to string", input.Data)
	}
	htmlContent := blackfriday.Run([]byte(text))
	response := Payload{
		Content: "text/html",
	}
	response.Data = string(htmlContent)
	return response, nil
}

func DirEntries(fileInfo FileInfo) ([]FileInfo, error) {
	entries, err := os.ReadDir(fileInfo.Name)
	if err != nil {
		return nil, err
	}
	list := []FileInfo{}

	for _, entry := range entries {

		stat, err := os.Stat(filepath.Join(fileInfo.Name, entry.Name()))

		f := FileInfo{
			Name:    entry.Name(),
			Size:    stat.Size(),
			Mode:    stat.Mode(),
			ModTime: stat.ModTime(),
			IsDir:   entry.IsDir(),
			Error:   err,
		}
		list = append(list, f)
	}
	return list, nil
}

func TextString(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	return string(bytes), err
}

func ToText(input Payload) (Payload, error) {
	fileInfo, ok := input.Data.(FileInfo)
	if !ok {
		return input, fmt.Errorf("cannot convert '%T' to FileInfo", input.Data)
	}

	if fileInfo.IsDir {
		entries, err := DirEntries(fileInfo)
		if err != nil {
			return input, err
		}
		res := Payload{
			Content: "application/json+[]fileinfo",
			Data:    entries,
		}
		return res, nil
	}

	string, err := TextString(fileInfo.Name)
	if err != nil {
		return input, err
	}

	res := Payload{
		Content: "text/markdown",
		Data:    string,
	}
	return res, nil
}
