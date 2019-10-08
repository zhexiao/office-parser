package utils

//基本类型
type BasicType string

const (
	BT_DANXT BasicType = "单选题"
	BT_DUOXT BasicType = "多选题"
	BT_TKT   BasicType = "填空"
	BT_PDT   BasicType = "判断题"
	BT_JDT   BasicType = "解答题"
	BT_ZWT   BasicType = "作文题"
)

func (b BasicType) Val() string {
	switch b {
	case BT_DANXT:
		return "DANXT"
	case BT_DUOXT:
		return "DUOXT"
	case BT_TKT:
		return "TKT"
	case BT_PDT:
		return "PDT"
	case BT_JDT:
		return "JDT"
	case BT_ZWT:
		return "ZWT"
	}
	return ""
}

//应用场景
type ResUsage int

const (
	RU_ILA ResUsage = 1
	RU_LAS ResUsage = 2
	RU_CA  ResUsage = 3
	RU_PKG ResUsage = 4
	RU_HW  ResUsage = 5
)

func (r ResUsage) Val() string {
	switch r {
	case RU_ILA:
		return "ILA"
	case RU_LAS:
		return "LAS"
	case RU_CA:
		return "CA"
	case RU_PKG:
		return "PKG"
	case RU_HW:
		return "HW"
	}
	return ""
}
