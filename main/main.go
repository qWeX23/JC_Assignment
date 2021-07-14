package main

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/qwex23/JC_Assignment/stats"
)

type Sample struct {
	Action string `json:"action"`
	Time   int64  `json:"time"`
}

func main() {

	st := stats.NewStats()
	for i := 0; i < 100000; i++ {
		a := stats.Sample{
			Action: "jump",
			Time:   uint64(rand.Intn(10000000)),
		}

		b, _ := json.Marshal(a)

		go st.AddAction(string(b))

		if i%2 == 0 {
			//go println(st.GetStats())
		}

	}

	call1 := "{\"action\":\"jump\", \"time\":100}"
	call2 := "{\"action\":\"run\", \"time\":75}"
	call3 := "{\"action\":\"jump\", \"time\":200}"
	st.AddAction(call1)
	st.AddAction(call2)
	st.AddAction(call3)

	fmt.Scanln()
	print(st.GetStats())

}
