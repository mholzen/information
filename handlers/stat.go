package handlers

import (
	"bufio"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
	Error   error
}

func (f FileInfo) Reader() (*bufio.Reader, error) {
	file, err := os.Open(f.Name)
	if err != nil {
		return nil, err
	}
	// defer file.Close()	// TODO: associate with context
	return bufio.NewReader(file), nil
}

func (f FileInfo) header() ([]byte, error) {
	reader, err := f.Reader()
	if err != nil {
		return nil, err
	}

	header, err := reader.Peek(512)
	if err != io.EOF && err != nil {
		return nil, err
	}
	return header, nil
}

func (f FileInfo) ContentType() (string, error) {
	if mimeType := mime.TypeByExtension(Extension(f.Name)); mimeType != "" {
		return mimeType, nil
	}
	header, err := f.header()
	if err != nil {
		return "", err
	}
	return http.DetectContentType(header), nil
}

func (f FileInfo) String() (string, error) {
	reader, err := f.Reader()
	if err != nil {
		return "", err
	}
	res, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func NewFileInfo(stat os.FileInfo) FileInfo {
	return FileInfo{
		Name:    stat.Name(),
		Size:    stat.Size(),
		Mode:    stat.Mode(),
		ModTime: stat.ModTime(),
		IsDir:   stat.IsDir(),
	}
}

func StatsHandler(c echo.Context) error {
	path := Filepath(c.Param("file"))
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err != nil {
		return err
	}
	if stat.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		list := []FileInfo{}

		for _, entry := range entries {

			stat, err := os.Stat(filepath.Join(path, entry.Name()))

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

		return c.JSON(http.StatusOK, list)
	} else {
		return c.JSON(http.StatusOK, NewFileInfo(stat))
	}
}
