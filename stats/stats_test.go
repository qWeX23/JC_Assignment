package stats

import "testing"

func TestStats_AddAction(t *testing.T) {
	type args struct {
		sampleString string
	}
	tests := []struct {
		name string
		s    *Stats
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.AddAction(tt.args.sampleString)
		})
	}
}

func Baseline() {
	call1 := "{\"action\":\"jump\", \"time\":100}"
	call2 := "{\"action\":\"run\", \"time\":75}"
	call3 := "{\"action\":\"jump\", \"time\":200}"

	st := NewStats()

	st.AddAction(call1)
	st.AddAction(call2)
	st.AddAction(call3)


	

}
