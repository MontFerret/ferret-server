package dal

type (
	Comparator int

	LogicOperator int

	FilteringField struct {
		Name       string      `json:"name"`
		Value      interface{} `json:"value"`
		Comparator Comparator  `json:"comparator"`
	}

	Filtering struct {
		Fields   []FilteringField
		Operator LogicOperator
	}
)

const (
	ComparatorEq  Comparator = 0
	ComparatorGt  Comparator = 1
	ComparatorGte Comparator = 2
	ComparatorLt  Comparator = 3
	ComparatorLte Comparator = 4
	ComparatorIn  Comparator = 5

	LogicOperatorAnd LogicOperator = 0
	LogicOperatorOr  LogicOperator = 1
)

var (
	comparatorValues = map[Comparator]string{
		ComparatorEq:  "==",
		ComparatorGt:  ">",
		ComparatorGte: ">=",
		ComparatorLt:  "<",
		ComparatorLte: "<=",
		ComparatorIn:  "IN",
	}

	logicOperatorValues = map[LogicOperator]string{
		LogicOperatorAnd: "AND",
		LogicOperatorOr:  "OR",
	}
)

func (c Comparator) String() string {
	v, found := comparatorValues[c]

	if found {
		return v
	}

	return comparatorValues[ComparatorEq]
}

func (l LogicOperator) String() string {
	v, found := logicOperatorValues[l]

	if found {
		return v
	}

	return logicOperatorValues[LogicOperatorAnd]
}
