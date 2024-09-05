package gofilter

const (
	IN  = iota
	NIN // 和in相反
	GT
	LT
	EQ
	NE
	LE
	GE

	//in  （contains）                包含
	//nin （not contains）            不包含
	//lt  （less than）               小于
	//le  （less than or equal to）   小于等于
	//eq  （equal to）                等于
	//ne  （not equal to）            不等于
	//ge  （greater than or equal to）大于等于
	//gt  （greater than）            大于
	// in 和 nin、eq、ne 用于string
	// 整形和浮点型可以使用全部
)
