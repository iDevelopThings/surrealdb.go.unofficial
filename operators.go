package surrealdb

type Operator struct {
	value string
}

func (o *Operator) String() string {
	return o.value
}

type OperatorTypes struct {
	// Symbol Operators
	Exact           Operator
	NotEqual        Operator
	AllEqual        Operator
	AnyEqual        Operator
	Equal           Operator
	NotLike         Operator
	AllLike         Operator
	AnyLike         Operator
	Like            Operator
	LessThanOrEqual Operator
	LessThan        Operator
	MoreThanOrEqual Operator
	MoreThan        Operator
	Add             Operator
	Sub             Operator
	Mul             Operator
	Div             Operator

	// Phrase Operators
	And         Operator
	Or          Operator
	ContainAll  Operator
	ContainAny  Operator
	ContainNone Operator
	NotContain  Operator
	Contain     Operator
	AllInside   Operator
	AnyInside   Operator
	NoneInside  Operator
	NotInside   Operator
	Inside      Operator
	Outside     Operator
	Intersects  Operator
}

var Operators = OperatorTypes{
	// Symbol Operators
	Exact:           Operator{value: "=="},
	NotEqual:        Operator{value: "!="},
	AllEqual:        Operator{value: "*="},
	AnyEqual:        Operator{value: "?="},
	Equal:           Operator{value: "="},
	NotLike:         Operator{value: "!~"},
	AllLike:         Operator{value: "*~"},
	AnyLike:         Operator{value: "?~"},
	Like:            Operator{value: "~"},
	LessThanOrEqual: Operator{value: "<="},
	LessThan:        Operator{value: "<"},
	MoreThanOrEqual: Operator{value: ">="},
	MoreThan:        Operator{value: ">"},
	Add:             Operator{value: "+"},
	Sub:             Operator{value: "-"},
	Mul:             Operator{value: "*"},
	Div:             Operator{value: "/"},

	// Phrase Operators
	And:         Operator{value: "AND"},
	Or:          Operator{value: "OR"},
	ContainAll:  Operator{value: "CONTAINSALL"},
	ContainAny:  Operator{value: "CONTAINSANY"},
	ContainNone: Operator{value: "CONTAINSNONE"},
	NotContain:  Operator{value: "CONTAINSNOT"},
	Contain:     Operator{value: "CONTAINS"},
	AllInside:   Operator{value: "ALLINSIDE"},
	AnyInside:   Operator{value: "ANYINSIDE"},
	NoneInside:  Operator{value: "NONEINSIDE"},
	NotInside:   Operator{value: "NOTINSIDE"},
	Inside:      Operator{value: "INSIDE"},
	Outside:     Operator{value: "OUTSIDE"},
	Intersects:  Operator{value: "INTERSECTS"},
}
