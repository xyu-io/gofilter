package gofilter

import (
	"container/list"
	"errors"
	"fmt"
	"strings"
)

type Stack struct {
	l *list.List
}

func NewStack() *Stack {
	return &Stack{
		l: list.New(),
	}
}

func (s *Stack) Push(v interface{}) {
	s.l.PushBack(v)
}

func (s *Stack) Pop() (interface{}, error) {
	if s.l.Len() == 0 {
		return nil, errors.New("stack is empty")
	}
	e := s.l.Back()
	s.l.Remove(e)

	return e.Value, nil
}

func (s *Stack) IsEmpty() bool {
	return s.l.Len() == 0
}

func (s *Stack) ToExpress() string {
	var res = ""
	l := s.l.Len()
	var front any
	var types []int  // 暂存
	var fns []string // 暂存
	for i := 0; i < l; i++ {
		if e, err := s.Pop(); err == nil {
			//fmt.Println("---", len(types), len(fns), "|", e)
			switch e.(type) {
			case string:
				// 获取前一个的类型，追加后面
				front = e.(string)
				fns = append(fns, e.(string))
				if i == l-1 {
					if len(types) == 1 { // 1
						str := strings.Join(fns, fmt.Sprintf(" %s ", getSign(types[0])))
						res = fmt.Sprintf("%s", res+fmt.Sprintf("(%s)", str))
					} else {
						str := strings.Join(fns, fmt.Sprintf(" %s ", getSign(types[1])))
						res = res + fmt.Sprintf(" %s (%s) ", getSign(types[0]), str)
					}
				}
			case int:
				if _, ok := front.(int); !ok && len(types) > 0 {
					// 处理
					if len(types) == 2 { // 2 types[0]
						if len(fns) == 1 && res != "" {
							res = getOnlyOne(res, fns[0], types[0])
						} else {
							str := strings.Join(fns, fmt.Sprintf(" %s ", getSign(types[1])))
							res = res + fmt.Sprintf(" %s (%s) ", getSign(types[0]), str)
							if e.(int) != types[0] {
								res = fmt.Sprintf("(%s)", res)
							}
						}
					}
					if len(types) == 1 { // 1
						if len(fns) == 1 && res != "" {
							res = getOnlyOne(res, fns[0], types[0])
						} else {
							str := strings.Join(fns, fmt.Sprintf(" %s ", getSign(types[0])))
							res = fmt.Sprintf("%s", res+fmt.Sprintf(" (%s)", str))
						}
					}

					types = make([]int, 0)
					fns = make([]string, 0)
				}
				types = append(types, e.(int))
				front = e.(int)
			}
		}
	}

	//fmt.Println(res)
	return res
}

func getOnlyOne(res, fn string, sign int) string {
	str := fmt.Sprintf(" %s %s ", getSign(sign), fn)

	return fmt.Sprintf("%s", res+str)
}

func getSign(types int) string {
	var sign = ""
	switch types {
	case OR:
		sign = "||"
	case AND:
		sign = "&&"
	default:
		sign = "||" // 默认为OR
	}
	return sign
}
