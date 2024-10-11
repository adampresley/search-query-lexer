package searchquerylexer

var DefaultComparatorConfig = ComparatorConfig{
	Equal:              "=",
	NotEqual:           "!=",
	LessThan:           "<",
	GreaterThan:        ">",
	LessThanEqualTo:    "<=",
	GreaterThanEqualTo: ">=",
	Like:               "=~",
	NotLike:            "!~",
}

var DefaultConnectiveConfig = ConnectiveConfig{
	And: "and",
	Or:  "or",
}
