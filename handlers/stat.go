package handlers

import (
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
