package whereparse

import (
	"errors"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

func checkQuery(query string) error {
	errorValues := []string{" LIMIT ", " ORDER BY "}
	for _, value := range errorValues {
		if strings.Contains(query, value) {
			return errors.New(fmt.Sprintf("The query contains a %s expression", value[1:len(value)-1]))
		}
	}
	return nil
}

func getExpression(query string) []string {
	expression := make([]string, 0)
	for len(query) != 0 || len(query) != 1 {
		spaceIdx := strings.IndexByte(query, ' ')
		if spaceIdx == -1 {
			spaceIdx = len(query)
		}
		expression = append(expression, query[:spaceIdx])
		if (len(query) - (len(expression[len(expression)-1]) + 1)) <= 0 {
			break
		}
		query = query[len(expression[len(expression)-1])+1:]
	}

	for i := range expression {
		expression[i] = strings.Trim(expression[i], "'")
		expression[i] = strings.Trim(expression[i], "\"")
	}
	return expression
}

func findFirstOper(src []string) int {
	operators := []string{"AND", "OR"}
	for i, j := range src {
		if j == operators[0] || j == operators[1] {
			return i
		}
	}
	return -1
}

func formRequest(src []string, checkValidity func(colName string, value interface{}) error) ([]sq.Sqlizer, error) { //src []string, qb *sq.SelectBuilder) []sq.Sqlizer {
	operIdx := findFirstOper(src)
	oper := src[operIdx]
	switch oper {
	case "AND":
		res := make(sq.And, 0)
		if operIdx < 3 {
			if err := checkValidity(src[0], nil); err != nil {
				return nil, err
			}
			res = append(res, sq.Expr(src[0]))
			src = src[2:]
		} else {
			switch src[1] {
			case "=":
				val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
				res = append(res, sq.Eq(val))
				break
			case "!=", "<>":
				val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
				res = append(res, sq.NotEq(val))
				break
			case ">":
				val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
				res = append(res, sq.Gt(val))
				break
			case ">=":
				val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
				res = append(res, sq.GtOrEq(val))
				break
			case "<":
				val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
				res = append(res, sq.Lt(val))
				break
			case "<=":
				val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
				res = append(res, sq.LtOrEq(val))
				break
			case "~", "!~":
				if err := checkValidity(src[0], src[2]); err != nil {
					return nil, err
				}
				res = append(res, sq.Expr(fmt.Sprintf("%s %s ?", src[0], src[1]), src[2]))
				break
			}
			src = src[4:]
		}
		operIdx = findFirstOper(src)
		if len(src) > 3 {
			switch src[operIdx] {
			case "AND":
				val, err := formRequest(src, checkValidity)
				if err!= nil {
					return nil, err
				}
				res = append(res, sq.And(val))
				break
			case "OR":
				val, err := formRequest(src, checkValidity)
				if err!= nil {
					return nil, err
				}
				res = append(res, sq.Or(val))
				break
			}
		} else {
			if len(src) < 3 {
				if err := checkValidity(src[0], nil); err != nil {
					return nil, err
				}
				res = append(res, sq.Expr(src[0]))
			} else {
				switch src[1] {
				case "=":
					val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
					res = append(res, sq.Eq(val))
					break
				case "!=", "<>":
					val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
					res = append(res, sq.NotEq(val))
					break
				case ">":
					val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
					res = append(res, sq.Gt(val))
					break
				case ">=":
					val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
					res = append(res, sq.GtOrEq(val))
					break
				case "<":
					val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
					res = append(res, sq.Lt(val))
					break
				case "<=":
					val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
					res = append(res, sq.LtOrEq(val))
					break
				case "~", "!~":
					if err := checkValidity(src[0], src[2]); err != nil {
						return nil, err
					}
					res = append(res, sq.Expr(fmt.Sprintf("%s %s ?", src[0], src[1]), src[2]))
					break
				}
			}
		}
		return res, nil

	case "OR":
		res := make(sq.Or, 0)
		if operIdx < 3 {
			if err := checkValidity(src[0], nil); err != nil {
				return nil, err
			}
			res = append(res, sq.Expr(src[0]))
			src = src[2:]
		} else {
			switch src[1] {
			case "=":
				val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
				res = append(res, sq.Eq(val))
				break
			case "!=", "<>":
				val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
				res = append(res, sq.NotEq(val))
				break
			case ">":
				val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
				res = append(res, sq.Gt(val))
				break
			case ">=":
				val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
				res = append(res, sq.GtOrEq(val))
				break
			case "<":
				val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
				res = append(res, sq.Lt(val))
				break
			case "<=":
				val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
				res = append(res, sq.LtOrEq(val))
				break
			case "~", "!~":
				if err := checkValidity(src[0], src[2]); err != nil {
					return nil, err
				}
				res = append(res, sq.Expr(fmt.Sprintf("%s %s ?", src[0], src[1]), src[2]))
				break
			}
			src = src[4:]
		}
		operIdx = findFirstOper(src)
		if checkLogOper(src) {
			switch src[operIdx] {
			case "AND":
				val, err := formRequest(src, checkValidity)
				if err!= nil {
					return nil, err
				}
				res = append(res, sq.And(val))
				break
			case "OR":
				val, err := formRequest(src, checkValidity)
				if err!= nil {
					return nil, err
				}
				res = append(res, sq.Or(val))
				break
			}
		} else {
			if len(src) < 3 {
				if err := checkValidity(src[0], nil); err != nil {
					return nil, err
				}
				res = append(res, sq.Expr(src[0]))
			} else {
				switch src[1] {
				case "=":
					val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
					res = append(res, sq.Eq(val))
					break
				case "!=", "<>":
					val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
					res = append(res, sq.NotEq(val))
					break
				case ">":val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
					res = append(res, sq.Gt(val))
					break
				case ">=":
					val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
					res = append(res, sq.GtOrEq(val))
					break
				case "<":
					val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
					res = append(res, sq.Lt(val))
					break
				case "<=":
					val, err := setExpression(src[:3], checkValidity)
				if err!= nil {
					return nil, err
				}
					res = append(res, sq.LtOrEq(val))
					break
				case "~", "!~":
					if err := checkValidity(src[0], src[2]); err != nil {
						return nil, err
					}
					res = append(res, sq.Expr(fmt.Sprintf("%s %s ?", src[0], src[1]), src[2]))
					break
				}
			}
		}
		return res, nil
	default:
		return nil, nil
	}
}

func setExpression(src []string, checkValidity func(colName string, value interface{}) error) (map[string]interface{}, error) {
	if err := checkValidity(src[0], src[2]); err != nil{
		return nil, err
	}
	
	switch src[1] {
	case "=":
		return sq.Eq{src[0]: src[2]}, nil
	case "!=", "<>":
		return sq.NotEq{src[0]: src[2]}, nil
	case ">":
		return sq.Gt{src[0]: src[2]}, nil
	case ">=":
		return sq.GtOrEq{src[0]: src[2]}, nil
	case "<":
		return sq.Lt{src[0]: src[2]}, nil
	case "<=":
		return sq.LtOrEq{src[0]: src[2]}, nil
	default:
		return nil, nil
	}
}

func getWhere(src string) int {
	return strings.Index(src, "WHERE")
}

func checkLogOper(src []string) bool {
	operators := []string{"AND", "OR"}
	for _, j := range src {
		if j == operators[0] || j == operators[1] {
			return true
		}
	}
	return false
}

func Parse(query string, qb sq.SelectBuilder, checkValidity func(colName string, value interface{}) error) (*sq.SelectBuilder, error) {
	if len(query) == 0 {
		return nil, errors.New("No query")
	}
	if err := checkQuery(query); err != nil {
		return nil, err
	}
	if err := getWhere(query); err != -1 {
		query = query[err+6:]
	}
	expressions := getExpression(query)

	if checkLogOper(expressions) {
		switch expressions[findFirstOper(expressions)] {
		case "AND":
			val, err := formRequest(expressions, checkValidity)
				if err!= nil {
					return nil, err
				}
			qb = qb.Where(sq.And(val))
			break
		case "OR":
			val, err := formRequest(expressions, checkValidity)
				if err!= nil {
					return nil, err
				}
			qb = qb.Where(sq.Or(val))
			break
		}
	} else {
		if len(expressions) < 3 {
			if err := checkValidity(expressions[0], nil); err != nil {
				return nil, err
			}
			qb = qb.Where(sq.Expr(expressions[0]))
		} else {
			switch expressions[1] {
			case "=":
				val, err := setExpression(expressions, checkValidity)
				if err!= nil {
					return nil, err
				}
				qb = qb.Where(sq.Eq(val))
				break
			case "!=", "<>":
				val, err := setExpression(expressions, checkValidity)
				if err!= nil {
					return nil, err
				}
				qb = qb.Where(sq.NotEq(val))
				break
			case ">":
				val, err := setExpression(expressions, checkValidity)
				if err!= nil {
					return nil, err
				}
				qb = qb.Where(sq.Gt(val))
				break
			case ">=":
				val, err := setExpression(expressions, checkValidity)
				if err!= nil {
					return nil, err
				}
				qb = qb.Where(sq.GtOrEq(val))
				break
			case "<":
				val, err := setExpression(expressions, checkValidity)
				if err!= nil {
					return nil, err
				}
				qb = qb.Where(sq.Lt(val))
				break
			case "<=":
				val, err := setExpression(expressions, checkValidity)
				if err!= nil {
					return nil, err
				}
				qb = qb.Where(sq.LtOrEq(val))
				break
			case "~", "!~":
				if err := checkValidity(expressions[0], expressions[2]); err != nil {
					return nil, err
				}
				qb = qb.Where(sq.Expr(expressions[0], expressions[2]))
				break
			}
		}
	}

	return &qb, nil
}
