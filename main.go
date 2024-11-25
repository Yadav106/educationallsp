package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/Yadav106/educationallsp/analysis"
	"github.com/Yadav106/educationallsp/lsp"
	"github.com/Yadav106/educationallsp/rpc"
)

func main() {
	logger := getLogger("/Users/macbook/Desktop/Programming/educationallsp/log.txt")
	logger.Println("Mic Check! 1! 2! 3!")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analysis.NewState()
	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s", err)
			continue
		}

		handleMessage(logger, writer, state, method, content)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state analysis.State, method string, content []byte) {
	logger.Printf("Received message with method: %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Hey, we couldn't parse this: %s", err)
		}

		logger.Printf("Connected to: %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version,
		)

		msg := lsp.NewInitializeResponse(request.ID)
		writeResponse(writer, msg)

		logger.Println("Sent Reply!")

	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("ERROR! textDocument/didOpen: %s", err)
			return
		}

		logger.Printf("Opened: %s", request.Params.TextDocument.URI)
		diagnostics := state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
		writeResponse(writer, lsp.PublishDiagnosticsNotification{
			Notification: lsp.Notification{
				RPC:    "2.0",
				Method: "textDocument/publishDiagnostics",
			},
			Params: lsp.PublishDiagnosticsParams{
				URI:         request.Params.TextDocument.URI,
				Diagnostics: diagnostics,
			},
		})

	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("ERROR! textDocument/didChange: %s", err)
			return
		}

		logger.Printf("Changed: %s", request.Params.TextDocument.URI)
		for _, change := range request.Params.ContentChanges {
			diagnostics := state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
			writeResponse(writer, lsp.PublishDiagnosticsNotification{
				Notification: lsp.Notification{
					RPC:    "2.0",
					Method: "textDocument/publishDiagnostics",
				},
				Params: lsp.PublishDiagnosticsParams{
					URI:         request.Params.TextDocument.URI,
					Diagnostics: diagnostics,
				},
			})
		}

	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("ERROR! textDocument/hover: %s", err)
			return
		}

		// create a response
		response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		// write it back
		writeResponse(writer, response)

	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("ERROR! textDocument/definition: %s", err)
			return
		}

		// create a response
		response := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		// write it back
		writeResponse(writer, response)

	case "textDocument/codeAction":
		var request lsp.CodeActionRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("ERROR! textDocument/codeAction: %s", err)
			return
		}

		// create a response
		response := state.TextDocumentCodeAction(request.ID, request.Params.TextDocumentIdentifer.URI)
		// write it back
		writeResponse(writer, response)

	case "textDocument/completion":
		var request lsp.CompletionRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("ERROR! textDocument/hover: %s", err)
			return
		}

		// create a response
		response := state.TextDocumentCompletion(request.ID, request.Params.TextDocument.URI)
		// write it back
		writeResponse(writer, response)

	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}

func getLogger(fileName string) *log.Logger {
	logfile, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("give a better file 🗿")
	}

	return log.New(logfile, "[educationallsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
