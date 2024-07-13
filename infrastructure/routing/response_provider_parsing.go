package routing

import (
	"fmt"
	"net/http"

	"github.com/prskr/go-dito/core/ports"
	"github.com/prskr/go-dito/infrastructure/grammar"
)

func ParseResponseProvider(call *grammar.Call, cwd ports.CWD) (ports.ResponseProvider, error) {
	switch call.Signature() {
	case "json(string)":
		rawJson, _ := call.Params[0].AsString()
		return Json(http.StatusOK, rawJson), nil
	case "json(int,string)":
		status, _ := call.Params[0].AsInt()
		rawJson, _ := call.Params[1].AsString()

		return Json(status, rawJson), nil
	case "file(string)":
		filePath, _ := call.Params[0].AsString()
		return File(cwd, http.StatusOK, filePath, ""), nil
	case "file(string,string)":
		filePath, _ := call.Params[0].AsString()
		contentType, _ := call.Params[1].AsString()

		return File(cwd, http.StatusOK, filePath, contentType), nil
	case "file(int,string)":
		status, _ := call.Params[0].AsInt()
		filePath, _ := call.Params[1].AsString()

		return File(cwd, status, filePath, ""), nil
	case "file(int,string,string)":
		status, _ := call.Params[0].AsInt()
		filePath, _ := call.Params[1].AsString()
		contentType, _ := call.Params[2].AsString()

		return File(cwd, status, filePath, contentType), nil
	default:
		return nil, fmt.Errorf("unknown response provider '%s'", call.String())
	}
}
