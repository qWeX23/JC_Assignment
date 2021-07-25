package main

import (
	"github.com/qwex23/JC_Assignment/stats"
)

func main() {

	st := stats.NewStats()

	call1 := "{\"AcTion\":\"jump\", \"tIMe\":100}"
	call2 := "{\"action\":\"run\", \"time\":75}"
	call3 := "{\"action\":\"jump\", \"time\":200}"

	st.AddAction(call1)
	st.AddAction(call2)
	st.AddAction(call3)

	statsJson, err := st.GetStats()
	if err != nil {
		println("Bad news, we had an error!")
	}
	print(statsJson)

}
