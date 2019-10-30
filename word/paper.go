package word

type CT_PureWord struct {
	SourceUri string `json:"source_uri"`
	WordText  string `json:"word_text"`
}

func NewCT_PureWord() *CT_PureWord {
	return &CT_PureWord{}
}

func ParsePaper(w *Word) *CT_PureWord {
	pWord := NewCT_PureWord()
	pWord.WordText = w.getPureText()

	return pWord
}
