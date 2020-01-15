package word

import "github.com/zhexiao/office-parser/bases"

type CT_PureWord struct {
	SourceUri string `json:"source_uri"`
	WordText  string `json:"word_text"`
}

func NewCT_PureWord() *CT_PureWord {
	return &CT_PureWord{}
}

//解析试卷
func ParsePaper(fileBytes []byte) (*CT_PureWord, error) {
	ctWord := NewCT_Word()

	err := ctWord.read(fileBytes)
	if err != nil {
		return nil, bases.NewOpError(bases.NormalError, err.Error())
	}

	data, err := ctWord.getWordData()
	if err != nil {
		return nil, bases.NewOpError(bases.NormalError, err.Error())
	}

	pWord := NewCT_PureWord()
	pWord.WordText = data

	return pWord, nil
}
