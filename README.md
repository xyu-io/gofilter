# gofilter
A data filter based on the go struct data type

### exp

```go
func TestSelfDataRule(t *testing.T) {
	data := []struct {
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
				CParams: []any{"c"},
				CType:   "Origin",
				CSymbol: "in",
			},
		},
	}

	rl, err := GenRule("test_rule_self", ruleEmpty)
	if err != nil {
		t.Error(err)
		return
	}

	for _, item := range data {
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
