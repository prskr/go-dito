package routing

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/prskr/go-dito/core/ports"
)

const sneakLength = 512

func StatusCode(status int) ports.ResponseProvider {
	return ports.ResponseProviderFunc(func(writer http.ResponseWriter) {
		writer.WriteHeader(status)
	})
}

func Json(status int, inlineJson string) ports.ResponseProvider {
	return ports.ResponseProviderFunc(func(writer http.ResponseWriter) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(status)
		if _, err := writer.Write([]byte(inlineJson)); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	})
}

func File(status int, filePath, contentType string) ports.ResponseProvider {
	return &FileProvider{
		Status:      status,
		FilePath:    filePath,
		ContentType: contentType,
	}
}

var (
	_ ports.ResponseProvider = (*FileProvider)(nil)
)

type FileProvider struct {
	Status      int
	FilePath    string
	ContentType string
}

func (f *FileProvider) Apply(writer http.ResponseWriter) {
	file, err := os.Open(f.FilePath)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	defer func() {
		_ = file.Close()
	}()

	var reader io.Reader = file

	if f.ContentType == "" {
		bufReader := bufio.NewReader(reader)
		sneakData, err := bufReader.Peek(sneakLength)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		f.ContentType = http.DetectContentType(sneakData)
		reader = bufReader
	}

	writer.Header().Set("Content-Type", f.ContentType)
	writer.WriteHeader(f.Status)

	if _, err := io.Copy(writer, reader); err != nil && !errors.Is(err, io.EOF) {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func JSFile(status int, filePath string) ports.ResponseProvider {
	return ports.ResponseProviderFunc(func(writer http.ResponseWriter) {
		writer.WriteHeader(http.StatusInternalServerError)
	})
}
