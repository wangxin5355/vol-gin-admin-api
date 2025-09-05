package model

type CPrefix int

const (
	Role     CPrefix = iota // 0
	UID                     // 1
	HDImg                   // 2
	Token                   // 3
	CityList                // 4
)

func (s CPrefix) String() string {
	switch s {
	case Role:
		return "0"
	case UID:
		return "1"
	case HDImg:
		return "2"
	case Token:
		return "3"
	case CityList:
		return "4"
	default:
		return "UNKNOWN"
	}
}
