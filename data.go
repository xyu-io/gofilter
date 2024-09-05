package gofilter

import (
	"fmt"
	"math/rand"
	"time"
)

type Data struct {
	DTime  int64
	DLevel int
	DType  string
	DMsg   any
}

// GenData 生成模拟数据
func GenData(col int) []Data {
	var data = make([]Data, 0)
	for i := 0; i < col; i++ {
		time.Sleep(time.Millisecond)
		data = append(data, Data{
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

	return data
}
