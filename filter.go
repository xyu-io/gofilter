package gofilter

import (
	"errors"
	"fmt"
	"strings"
)

const (
	AND = 0
	OR  = 1
)

type Filter struct {
	rTag    string              `json:"rtag,omitempty"`
	rExpr   string              `json:"expression,omitempty"`
	rFnMaps map[string]RulePool `json:"fnmaps,omitempty"`
	Rule    Rule                `json:"rules,omitempty"`

	execIns *GValuate // 规则执行器
}

type Rule struct {
	RSign  int        `json:"sign,omitempty"`   // 用于RRules和RChild的逻辑运算符
	RType  int        `json:"rtype,omitempty"`  // 用于RChild之间的逻辑运算符
	RChild *Rule      `json:"rchild,omitempty"` // children级别rules []Rule不好在rule之间逻辑运算
	RRules []RulePool `json:"rrules,omitempty"` // 本级rules
}

type RulePool struct {
	CName   string // 函数标识
	CParams []any  // 参数
	CType   string // 字段名称
	CSymbol int    // 逻辑运算符号（包含、大于、小于、不等于、等于等）
}

type FilterFunc func(option any, data any) bool // 根据过滤

func NewFilter(tag string, rule Rule) (*Filter, error) {
	var r = &Filter{
		rTag:    tag,
		rFnMaps: make(map[string]RulePool),
		rExpr:   "", // fn1 AND  (fn2.1 OR fn2.0) AND  (fn3 AND fn2)
		Rule:    rule,
		execIns: NewGValuate(tag),
	}

	stack := NewStack()
	r.genExp(stack, r.Rule)
	expr := stack.ToExpress()
	if expr == "" {
		expr = "true"
	}
	r.rExpr = strings.TrimSpace(expr)

	err := r.execIns.Eval(r.expr())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s 【%s】: %s", r.rTag, r.rExpr, err.Error()))
	}

	//log.Printf("init succeed. tag: %s, rule express: %s", r.rTag, r.rExpr)
	return r, nil
}

func (f *Filter) expr() string {
	return f.rExpr
}

func (f *Filter) Exec(data any) (bool, error) {

	flag, err := f.execIns.Exec(GetFnMaps(data, f.getFns()))
	if err != nil {
		return false, err
	}
	return flag, nil
}

func (f *Filter) Tag() string {
	return f.rTag
}

func (f *Filter) getFns() map[string]RulePool {
	return f.rFnMaps
}

func (f *Filter) genExp(s *Stack, rs Rule) {
	// 入栈操作符，操作数
	for _, r := range rs.RRules {
		if _, ok := f.rFnMaps[r.CName]; !ok {
			f.rFnMaps[r.CName] = r
		}
		s.Push(r.CName)
	}
	s.Push(rs.RType)
	if rs.RChild != nil { // 有子节点，先处理子节点
		s.Push(rs.RSign)
		f.genExp(s, *rs.RChild)
	}

}

// GetDataField 通过反射获取结构体字段名称, data类型必须是struct
func GetDataField(data any, ctype string) (string, any) {
	if ctype == "" {
		return "", nil
	}

	if field, value, flag := dealStructPtr(data, ctype); flag && value != nil {
		return field, value
	}

	return dealStruct(data, ctype)
}

// GetFnMaps 重要，处理规则和数据的关联部分
func GetFnMaps(item any, fs map[string]RulePool) map[string]interface{} {
	var mps = make(map[string]interface{})
	if _, ok := item.(struct{}); ok {
		return mps
	}
	for name, fn := range fs {
		if len(fn.CParams) == 0 || fn.CType == "" {
			mps[name] = func() bool { return true }()
			continue
		}
		_, value := GetDataField(item, fn.CType)
		if value == nil {
			continue
		}
		mps[name] = logicSelection(fn, value)
	}

	return mps
}

// 逻辑运算
func logicSelection(rp RulePool, value any) bool {
	var flag = true
	switch rp.CSymbol {
	case IN:
		flag = AnyFind(rp.CParams, value)
	case NIN:
		flag = !AnyFind(rp.CParams, value)
	case EQ: // 等于
		equal, err := AnyEqual(rp.CParams, value)
		if err != nil {
			return true // 不兼容类型不做处理
		}
		flag = equal
	case NE: // 不等于
		equal, err := AnyEqual(rp.CParams, value)
		if err != nil {
			return true // 不兼容类型不做处理
		}
		flag = !equal
	case GT: // 大于
		gt, err := AnyGreaterThan(rp.CParams, value)
		if err != nil {
			return true // 不兼容类型不做处理
		}
		flag = gt
	case LT: // 小于
		lt, err := AnyLessThan(rp.CParams, value)
		if err != nil {
			return true // 不兼容类型不做处理
		}
		flag = lt
	case LE: // 小于等于 等价于 不大于
		gt, err := AnyGreaterThan(rp.CParams, value)
		if err != nil {
			return true // 不兼容类型不做处理
		}
		flag = !gt
	case GE: // 大于等于 等价于 不小于
		lt, err := AnyLessThan(rp.CParams, value)
		if err != nil {
			return true // 不兼容类型不做处理
		}
		flag = lt
	default:
		// 不支持的类型
		flag = false
		panic("unhandled default case")
	}

	return flag
}
