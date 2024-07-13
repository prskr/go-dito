package routing

import (
	"bufio"
	"errors"
	"io"
	"net/http"

	"github.com/prskr/go-dito/core/ports"
)

const sneakLength = 512

func Json(status int, inlineJson string) ports.ResponseProvider {
	return ports.ResponseProviderFunc(func(writer http.ResponseWriter) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(status)
		if _, err := writer.Write([]byte(inlineJson)); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	})
}

func File(cwd ports.CWD, status int, filePath, contentType string) ports.ResponseProvider {
	return ports.ResponseProviderFunc(func(writer http.ResponseWriter) {
		f, err := cwd.Open(filePath)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		defer func() {
			_ = f.Close()
		}()

		var reader io.Reader = f

		if contentType == "" {
			bufReader := bufio.NewReader(reader)
			sneakData, err := bufReader.Peek(sneakLength)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}

			contentType = http.DetectContentType(sneakData)
			reader = bufReader
		}

		writer.Header().Set("Content-Type", contentType)
		writer.WriteHeader(status)

		if _, err := io.Copy(writer, reader); err != nil && !errors.Is(err, io.EOF) {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	})
}

func JSFile(cwd ports.CWD, status int, filePath string) ports.ResponseProvider {
	return ports.ResponseProviderFunc(func(writer http.ResponseWriter) {
		writer.WriteHeader(http.StatusInternalServerError)
	})
}
