package gamedata

type QualityType int64

const (
	QualityTypeLowest QualityType = iota
	QualityTypeWhite
	QualityTypeGreen
	QualityTypeBlue
	QualityTypePurple
	QualityTypeOrange
	QualityTypeRed
	QualityTypeHighest
)
