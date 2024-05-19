package base

const ParamBaseKey = "p"

type OrderType string

const (
	OrderTypeASC  OrderType = "ASC"
	OrderTypeDESC OrderType = "DESC"
)

type ConditionOperator string

const (
	ConditionOperatorEq ConditionOperator = "="
	ConditionOperatorIn ConditionOperator = "IN"
)

type SearchResultType int

const (
	SearchResultTypeNotSearched SearchResultType = 0
	SearchResultTypeFound       SearchResultType = 1
	SearchResultTypeNotFound    SearchResultType = 2
)

type OperationType int

const (
	OperationTypeInsert OperationType = iota + 1
	OperationTypeUpdate
	OperationTypeDelete
)

type OrderPair struct {
	Column    string
	OrderType OrderType
}

type OrderPairs []*OrderPair
