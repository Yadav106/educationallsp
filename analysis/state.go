package analysis

type State struct {
	// map of filenames to content
	Documents map[string]string
}

func NewState() State {
	return State{
		Documents: map[string]string{},
	}
}

func (s *State) OpenDocument(document, text string) {
	s.Documents[document] = text
}
