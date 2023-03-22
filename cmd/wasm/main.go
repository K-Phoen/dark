//go:build js && wasm

package main

import (
	"strings"
	"syscall/js"

	"github.com/K-Phoen/dark/internal/pkg/converter"
	"go.uber.org/zap"
)

func dashboardToDarkFunc() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "invalid nnumber of arguments passed"
		}

		logger := zap.NewNop()
		conv := converter.NewJSON(logger)

		inputJSON := args[0].String()
		output := &strings.Builder{}

		if err := conv.ToK8SManifest(strings.NewReader(inputJSON), output, converter.K8SManifestOptions{
			Folder:    "folder_name",
			Name:      "name",
			Namespace: "default",
		}); err != nil {
			logger.Fatal("Could not convert dashboard", zap.Error(err))
		}

		return output.String()
	})
}

func main() {
	js.Global().Set("dashboardToDark", dashboardToDarkFunc())

	c := make(chan struct{}, 0)
	<-c
}
