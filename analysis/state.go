package analysis

import (
	"fmt"
	"strings"

	"github.com/Yadav106/educationallsp/lsp"
)

type State struct {
	// map of filenames to content
	Documents map[string]string
}

func NewState() State {
	return State{
		Documents: map[string]string{},
	}
}

func getDiagnosticsForFile(text string) []lsp.Diagnostic {
  diagnostics := []lsp.Diagnostic{}
  for row, line := range strings.Split(text, "\n") {
    // diagnostics = append(diagnostics, )
    if strings.Contains(line, "VS Code") {
      idx := strings.Index(text, "VS Code")
      diagnostics = append(diagnostics, lsp.Diagnostic{
      	Range:    LineRange(row, idx, idx + len("VS Code")),
      	Severity: 1,
      	Source:   "Custom Built Language Server",
      	Message:  "Diagnostic provided by custom language server",
      })
    }
  }

  return diagnostics
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text

  return getDiagnosticsForFile(text)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text

  return getDiagnosticsForFile(text)
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	document := s.Documents[uri]
	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("File: %s, Characters %d", uri, len(document)),
		},
	}
}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	// document := s.Documents[uri]
	return lsp.DefinitionResponse{
		Response: lsp.Response{RPC: "2.0", ID: &id},
		Result: lsp.Location{
			URI: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
				End: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
			},
		},
	}
}

func (s *State) TextDocumentCodeAction(id int, uri string) lsp.CodeActionResponse {
	text := s.Documents[uri]
	actions := []lsp.CodeAction{}

	uri_split := strings.Split(uri, "/")
	title := uri_split[len(uri_split)-2]

	text_split := strings.Split(text, "\n")

	if (len(text_split) > 0) && (strings.Index(text_split[0], title) < 0) {
		appendChange := map[string][]lsp.TextEdit{}
		appendChange[uri] = []lsp.TextEdit{
			{
				Range:   LineRange(0, 0, 0),
				NewText: fmt.Sprintf("# %s\n", title),
			},
		}

		actions = append(actions, lsp.CodeAction{
			Title: "Add title to the markdown file",
			Edit: &lsp.WorkspaceEdit{
				Changes: appendChange,
			},
		})
	}

	for row, line := range text_split {
		idx := strings.Index(line, "VS Code")
		if idx >= 0 {
			replaceChange := map[string][]lsp.TextEdit{}
			replaceChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "Neovim",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: "Replace VS Code with a superior editor",
				Edit: &lsp.WorkspaceEdit{
					Changes: replaceChange,
				},
			})

			censorChange := map[string][]lsp.TextEdit{}
			censorChange[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "VS C*de",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: "Censor to VS C*de",
				Edit:  &lsp.WorkspaceEdit{Changes: censorChange},
			})
		}
	}

	return lsp.CodeActionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: actions,
	}
}

func (s *State) TextDocumentCompletion(id int, uri string) lsp.CompletionResponse {
	items := []lsp.CompletionItem{
    {
    	Label:         "LSP Autocomplete",
    	Detail:        "Testing Autocompletion",
    	Documentation: "Very cool isn't it",
    },
  }

	return lsp.CompletionResponse{
		Response: lsp.Response{RPC: "2.0", ID: &id},
		Result:   items,
	}
}

func LineRange(line, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      line,
			Character: start,
		},
		End: lsp.Position{
			Line:      line,
			Character: end,
		},
	}
}
