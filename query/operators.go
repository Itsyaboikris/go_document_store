package query

type Operator string

const (
    // Comparison Operators
    OpEquals       Operator = "$eq"
    OpNotEquals    Operator = "$ne"
    OpGreater      Operator = "$gt"
    OpGreaterEqual Operator = "$gte"
    OpLess         Operator = "$lt"
    OpLessEqual    Operator = "$lte"
    OpIn           Operator = "$in"
    OpNotIn        Operator = "$nin"
    
    // Logical Operators
    OpAnd          Operator = "$and"
    OpOr           Operator = "$or"
    OpNot          Operator = "$not"
    OpNor          Operator = "$nor"
    
    // Element Operators
    OpExists       Operator = "$exists"
    OpType         Operator = "$type"

    // Evaluation Operators
    OpRegex        Operator = "$regex"
    OpMod          Operator = "$mod"
    
    // Array Operators
    OpAll          Operator = "$all"
    OpSize         Operator = "$size"
    OpElemMatch    Operator = "$elemMatch"
)

func IsComparisonOperator(op Operator) bool {
    switch op {
    case OpEquals, OpNotEquals, OpGreater, OpGreaterEqual, OpLess, OpLessEqual, OpIn, OpNotIn:
        return true
    default:
        return false
    }
}

func IsLogicalOperator (op Operator) bool {
    switch op {
    case OpAnd, OpOr, OpNot, OpNor:
        return true
    default:
        return false
    }
}

func IsElementOperator(op Operator) bool {
    switch op {
    case OpExists, OpType:
        return true
    default:
        return false
    }
}

func IsEvaluationOperator(op Operator) bool {
    switch op {
    case OpRegex, OpMod:
        return true
    default:
        return false
    }
}

func IsArrayOperator(op Operator) bool {
    switch op {
    case OpAll, OpSize, OpElemMatch:
        return true
    default:
        return false
    }
}

func ValidateOperator(op Operator) bool {
    operator := Operator(op)
    return IsComparisonOperator(operator) || IsLogicalOperator(operator) || IsElementOperator(operator) || IsEvaluationOperator(operator) || IsArrayOperator(operator)
}