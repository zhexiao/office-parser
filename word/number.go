package word

type NUM_Decimal int

func (n NUM_Decimal) String() string {
	switch n {
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	case 5:
		return "5"
	case 6:
		return "6"
	case 7:
		return "7"
	case 8:
		return "8"
	case 9:
		return "9"
	case 10:
		return "10"
	default:
		return ""
	}
}

type NUM_DecimalEnclosedCircle int

func (n NUM_DecimalEnclosedCircle) String() string {
	switch n {
	case 1:
		return "①"
	case 2:
		return "②"
	case 3:
		return "③"
	case 4:
		return "④"
	case 5:
		return "⑤"
	case 6:
		return "⑥"
	case 7:
		return "⑦"
	case 8:
		return "⑧"
	case 9:
		return "⑨"
	case 10:
		return "⑩"
	default:
		return ""
	}
}

type NUM_Counting int

func (n NUM_Counting) String() string {
	switch n {
	case 1:
		return "一"
	case 2:
		return "二"
	case 3:
		return "三"
	case 4:
		return "四"
	case 5:
		return "五"
	case 6:
		return "六"
	case 7:
		return "七"
	case 8:
		return "八"
	case 9:
		return "九"
	case 10:
		return "十"
	default:
		return ""
	}
}

type NUM_UpperLetter int

func (n NUM_UpperLetter) String() string {
	switch n {
	case 1:
		return "A"
	case 2:
		return "B"
	case 3:
		return "C"
	case 4:
		return "D"
	case 5:
		return "E"
	case 6:
		return "F"
	case 7:
		return "G"
	case 8:
		return "H"
	default:
		return ""
	}
}

type NUM_UpperRoman int

func (n NUM_UpperRoman) String() string {
	switch n {
	case 1:
		return "Ⅰ"
	case 2:
		return "Ⅱ"
	case 3:
		return "Ⅲ"
	case 4:
		return "Ⅳ"
	case 5:
		return "Ⅴ"
	case 6:
		return "Ⅵ"
	case 7:
		return "Ⅶ"
	case 8:
		return "Ⅷ"
	case 9:
		return "Ⅸ"
	case 10:
		return "Ⅹ"
	default:
		return ""
	}
}
