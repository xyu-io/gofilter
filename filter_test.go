package gofilter

import (
	"encoding/json"
	"fmt"

	"math/rand"
	"testing"
	"time"
)

type data struct {
	DTime  int64
	DLevel int
	DType  string
	DMsg   any
}

// GenData 生成模拟数据
func genData(col int) []data {
	var dat = make([]data, 0)
	for i := 0; i < col; i++ {
		time.Sleep(time.Millisecond)
		dat = append(dat, data{
			DTime: time.Now().Unix(),
			DLevel: func() int {
				arr := []int{0, 1, 2, 3}
				randomIndex := rand.Intn(len(arr))
				return arr[randomIndex]
			}(),
			DType: func() string {
				arr := []string{"info", "normal", "important", "danger"}
				randomIndex := rand.Intn(len(arr))
				return arr[randomIndex]
			}(),
			DMsg: "dataTime-" + fmt.Sprintf("%d", time.Now().UnixNano()),
		})
	}

	return dat
}

func TestValuate(t *testing.T) {
	gvIns := NewGValuate("test")
	err := gvIns.Eval("(((a == 1) && (b || c))) || !d")
	if err != nil {
		t.Error(err)
	}
	parameters := map[string]interface{}{
		"a": 12,
		"b": false,
		"c": true,
		"d": true,
	}
	result, err := gvIns.Exec(parameters)
	if err != nil {
		t.Error(err)
	}
	t.Log(result)
}

func TestSampleRule(t *testing.T) {
	rule := Rule{
		RSign: AND,
		RType: OR,
		RRules: []RulePool{
			{
				CName:   "fn1_1",
				CParams: []any{"important"},
				CType:   "DType", //
			},
			{
				CName:   "fn1_2",
				CParams: []any{"danger"},
				CType:   "DType", // 来自于数据字段
			},
		},
		RChild: &Rule{
			//RSign: AND,
			RType: OR,
			RRules: []RulePool{
				{
					CName:   "fn2_1",
					CParams: []any{0},
					CType:   "DLevel", //
				},
				{
					CName:   "fn2_2",
					CParams: []any{2},
					CType:   "DLevel", // 来自于数据字段
				},
			},
		},
	}

	rl, err := NewFilter("test_rule", rule)
	if err != nil {
		t.Error(err)
		return
	}

	dat := genData(50)

	for _, item := range dat {
		flag, err := rl.Exec(item) //rl.ExecIns.Exec(GetFnMaps(*item, rl.getFns()))
		if err != nil {
			t.Error(err)
			continue
		}
		if flag {
			t.Logf("-->tag 【%s】 msg: {%s, %d, %s}", rl.Tag(), item.DType, item.DLevel, item.DMsg)
		}
	}

}

func TestSliceRule(t *testing.T) {
	rule := Rule{
		RSign: AND,
		RType: OR,
		RRules: []RulePool{
			{
				CName:   "fn1_1",
				CParams: []any{"important"},
				CType:   "DType", //
			},
			{
				CName:   "fn1_2",
				CParams: []any{"danger"},
				CType:   "DType", // 来自于数据字段
			},
		},
		RChild: &Rule{
			//RSign: AND,
			RType: OR,
			RRules: []RulePool{
				{
					CName:   "fn2_1",
					CParams: []any{0},
					CType:   "DLevel", //
				},
				{
					CName:   "fn2_2",
					CParams: []any{2},
					CType:   "DLevel", // 来自于数据字段
				},
			},
		},
	}

	rl, err := NewFilter("test_rule", rule)
	if err != nil {
		t.Error(err)
		return
	}

	dat := genData(50)

	res, err := rl.ExecWithSlice(dat)
	if err != nil {
		t.Error(err)
		return
	}
	bytes, err := json.Marshal(res)
	if err != nil {
		t.Error(err)
		return
	}

	var output []data
	err = json.Unmarshal(bytes, &output)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("output: %+v", output)
}

func TestSampleRuleOfPtr(t *testing.T) {
	rule := Rule{
		RSign: AND,
		RType: OR,
		RRules: []RulePool{
			{
				CName:   "fn1_1",
				CParams: []any{"important"},
				CType:   "DType", //
			},
			{
				CName:   "fn1_2",
				CParams: []any{"danger"},
				CType:   "DType", // 来自于数据字段
			},
		},
	}

	rl, err := NewFilter("test_rule_ptr", rule)
	if err != nil {
		t.Error(err)
		return
	}

	dat := []*data{
		&data{
			DTime:  time.Now().UnixNano(),
			DLevel: 1,
			DType:  "important",
			DMsg:   "test-ptr-important",
		},
		{
			DTime:  time.Now().UnixNano(),
			DLevel: 2,
			DType:  "danger",
			DMsg:   "test-ptr-danger",
		},
		{
			DTime:  time.Now().UnixNano(),
			DLevel: 0,
			DType:  "info",
			DMsg:   "test-ptr-info",
		},
	}

	for _, item := range dat {
		flag, err := rl.Exec(item) //rl.ExecIns.Exec(GetFnMaps(*item, rl.getFns()))
		if err != nil {
			t.Error(err)
			continue
		}
		if flag {
			t.Logf("-->tag 【%s】 msg: {%s, %d, %s}", rl.Tag(), item.DType, item.DLevel, item.DMsg)
		}
	}

}

func TestComplexRule(t *testing.T) {
	rule1 := Rule{
		RSign: AND,
		RType: AND,
		RRules: []RulePool{
			{
				CName:   "fn1_1",
				CParams: []any{"info"},
				CType:   "DType", //
			},
		},
		RChild: &Rule{
			RSign: AND,
			RType: OR,
			RRules: []RulePool{
				{
					CName:   "fn2_1",
					CParams: []any{2, 3},
					CType:   "DLevel", //
				},
			},
		},
	}

	rule2 := Rule{
		RSign: AND,
		RType: OR,
		RRules: []RulePool{
			{
				CName:   "fn1_1",
				CParams: []any{"important"},
				CType:   "DType", //
			},
			{
				CName:   "fn1_2",
				CParams: []any{"danger"},
				CType:   "DType", //
			},
		},
		RChild: &Rule{
			RType: OR,
			RRules: []RulePool{
				{
					CName:   "fn2_1",
					CParams: []any{1, 2, 3},
					CType:   "DLevel", //
				},
			},
		},
	}

	rule3 := Rule{
		RSign: AND,
		RType: OR,
		RRules: []RulePool{
			{
				CName:   "fn3_1",
				CParams: []any{"danger"},
				CType:   "DType", //
			},
		},
	}

	rls := []Rule{rule1, rule2, rule3}
	rIns := make([]*Filter, 0)
	for index, r := range rls {
		rl, err := NewFilter(fmt.Sprintf("test_rule_%d", index), r)
		if err != nil {
			t.Error(err)
			continue
		}
		rIns = append(rIns, rl)
	}

	dat := genData(50)

	for _, item := range dat {
		for _, rl := range rIns {
			flag, err := rl.Exec(item)
			if err != nil {
				t.Error(err)
				continue
			}
			if flag {
				t.Logf("-->tag 【%s】 msg: {%s, %d, %s}", rl.Tag(), item.DType, item.DLevel, item.DMsg)
			}
		}

	}

}

func TestEmptyRule(t *testing.T) {
	ruleEmpty := Rule{
		RSign:  OR,
		RType:  OR,
		RRules: []RulePool{},
	}

	rl, err := NewFilter("test_rule_empty", ruleEmpty)
	if err != nil {
		t.Error(err)
		return
	}

	dat := genData(5)

	for _, item := range dat {
		flag, err := rl.Exec(item)
		if err != nil {
			t.Error(err)
			continue
		}
		if flag {
			t.Logf("-->tag 【%s】 msg: {%s, %d, %s}", rl.Tag(), item.DType, item.DLevel, item.DMsg)
		}
	}

}

func TestSelfDataRule(t *testing.T) {
	dat := []struct {
		ID     int
		Msg    string
		Tag    string
		Origin string
	}{
		{
			ID:     1,
			Msg:    "hello-1",
			Tag:    "AAA",
			Origin: "a",
		},
		{
			ID:     2,
			Msg:    "hello-2",
			Tag:    "AA",
			Origin: "b",
		},
		{
			ID:     3,
			Msg:    "hello-3",
			Tag:    "AAA",
			Origin: "b",
		},
		{
			ID:     4,
			Msg:    "hello-4",
			Tag:    "AAA",
			Origin: "c",
		},
	}

	ruleEmpty := Rule{
		RType: OR,
		RRules: []RulePool{
			{
				CName:   "F2",
				CParams: []any{1},
				CType:   "ID",
				CSymbol: GT,
			},
		},
	}

	rl, err := NewFilter("test_rule_self", ruleEmpty)
	if err != nil {
		t.Error(err)
		return
	}

	for _, item := range dat {
		flag, err := rl.Exec(item)
		if err != nil {
			t.Error(err)
			continue
		}
		if flag {
			t.Logf("-->msg: {%v, %v, %v,%v}", item.ID, item.Tag, item.Origin, item.Msg)
		}
	}
}
