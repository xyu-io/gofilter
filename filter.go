package gofilter

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

const (
	AND = 0
	OR  = 1
	//NOT
)

type FilterRule struct {
	rName   string              `json:"rname,omitempty"`
	rSender string              `json:"sender,omitempty"`
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
	CSymbol string // 逻辑运算符号（包含、大于、小于、不等于、等于等）
}

type Filter func(option any, data any) bool // 根据过滤

func GenRule(name, sender string, rule Rule) (*FilterRule, error) {
	var r = &FilterRule{
		rName:   name,
		rSender: sender,
		rFnMaps: make(map[string]RulePool),
		rExpr:   "", // fn1 AND  (fn2.1 OR fn2.0) AND  (fn3 AND fn2)
		Rule:    rule,
		execIns: NewGValuate(name),
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
		return nil, errors.New(fmt.Sprintf("%s 【%s】: %s", r.rSender, r.rExpr, err.Error()))
	}

	log.Printf("init succeed sender: %s, rule express: %s", sender, r.rExpr)
	return r, nil
}

func (f *FilterRule) expr() string {
	return f.rExpr
}

func (f *FilterRule) Name() string {
	return f.rName
}

func (f *FilterRule) Exec(data any) (bool, error) {

	flag, err := f.execIns.Exec(GetFnMaps(data, f.getFns()))
	if err != nil {
		return false, err
	}
	return flag, nil
}

func (f *FilterRule) Sender() string {
	return f.rSender
}

func (f *FilterRule) getFns() map[string]RulePool {
	return f.rFnMaps
}

func (f *FilterRule) genExp(s *Stack, rs Rule) {
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

// 生成expression表达式
//func genExp(s *Stack, rs []Rule) {
//	for _, rl := range rs {
//		// 多叶子
//		if len(rl.RChild) > 0 { // 有子节点，先处理子节点
//			genExp(s, rl.RChild)
//		}
//
//		// 入栈操作符，操作数
//		for _, r := range rl.RRules {
//			s.Push(r.CName)
//		}
//		s.Push(rl.RType)
//	}
//}

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
		mps[name] = AnyFind(fn.CParams, value)
	}

	return mps
}
