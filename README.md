# GoFilter

## Description

A data filter based on the go struct data type

+ Supports logical operation data filtering based on struct data types.

+ Such as IN / NIN

  ```
    //in  （contains）                包含
    //nin （not contains）            不包含
    //lt  （less than）               小于
    //le  （less than or equal to）   小于等于
    //eq  （equal to）                等于
    //ne  （not equal to）            不等于
    //ge  （greater than or equal to）大于等于
    //gt  （greater than）            大于
    // in 和 nin、eq、ne 用于string
    // 整形和浮点型可以使用全部
    ```

## Exp

```go
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
				CName:   "F1", // 注意，同一个规则内部不能重名，否则会被覆盖
				CParams: []any{1},
				CType:   "ID",
				CSymbol: gofilter.GT,
			},
			{
				CName:   "F2",
				CParams: []any{"AAA", "AA"},
				CType:   "Tag",
				CSymbol: gofilter.IN,
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

```

