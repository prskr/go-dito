package routing

import (
	"bufio"
	"errors"
	"io"
	"net/http"

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
	_ ports.CwdInjectable    = (*FileProvider)(nil)
	_ ports.ResponseProvider = (*FileProvider)(nil)
)

type FileProvider struct {
	Status      int
	FilePath    string
	ContentType string
	CWD         ports.CWD
}

func (f *FileProvider) Apply(writer http.ResponseWriter) {
	file, err := f.CWD.Open(f.FilePath)
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

func (f *FileProvider) InjectCwd(cwd ports.CWD) {
	f.CWD = cwd
}

func JSFile(cwd ports.CWD, status int, filePath string) ports.ResponseProvider {
	return ports.ResponseProviderFunc(func(writer http.ResponseWriter) {
		writer.WriteHeader(http.StatusInternalServerError)
	})
}
