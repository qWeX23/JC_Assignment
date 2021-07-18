package stats

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"
)

//test basic functionality from strings
func TestAddAction(t *testing.T) {
	//Arrange
	st := NewStats()
	call1 := "{\"action\":\"jump\",\"time\":100}"
	call2 := "{\"action\":\"run\", \"time\":75}"
	call3 := "{\"action\":\"jump\", \"time\":200}"
	//Act
	err1 := st.AddAction(call1)
	err2 := st.AddAction(call2)
	err3 := st.AddAction(call3)
	//Assert

	if err1 != nil {
		t.Error(err1.Error())
	}
	if err2 != nil {
		t.Error(err2.Error())
	}
	if err3 != nil {
		t.Error(err3.Error())
	}

	//since the ordering of the items is unimportant, we should marshal
	//to ensure we are validating the values not the exact string.
	want := "[{\"action\":\"jump\",\"avg\":150},{\"action\":\"run\",\"avg\":75}]"

	wantStruct := make([]SampleAverage, 2)
	json.Unmarshal([]byte(want), &wantStruct)

	statsJson, err := st.GetStats()
	if err != nil {
		t.Errorf("error from get stats %s", err.Error())
	}
	haveStruct := make([]SampleAverage, 2)
	json.Unmarshal([]byte(statsJson), &haveStruct)

	//Terribly inneficient loop, but only for 4 passes
	for _, w := range wantStruct {
		for _, h := range haveStruct {
			if w.Action == h.Action {
				if w.Average != h.Average {
					t.Errorf("action of %s has wrong average %d want %d", w.Action, h.Average, w.Average)
				}
			}
		}
	}

}

//test basic funcitonality from the core code
func TestAddAction_Sample(t *testing.T) {
	//Arrange
	st := NewStats()
	call1 := Sample{
		Action: "jump",
		Time:   1,
	}
	call2 := Sample{
		Action: "run",
		Time:   0,
	}
	call3 := Sample{
		Action: "jump",
		Time:   3,
	}

	//Act
	st.addAction(call1)
	st.addAction(call2)
	st.addAction(call3)

	//Assert
	expectedTotalJumpTime := uint64(4)
	foundTotalJumpTime := st.Averages["jump"].TotalTime
	if foundTotalJumpTime != expectedTotalJumpTime {
		t.Fatalf("TotalTime calculation incorrect, expected %d but found %d", foundTotalJumpTime, expectedTotalJumpTime)
	}

	expectedTotalRunTime := uint64(0)
	fountTotalRunTime := st.Averages["run"].TotalTime
	if fountTotalRunTime != expectedTotalRunTime {
		t.Fatalf("TotalTime calculation incorrect, expected %d but found %d", fountTotalRunTime, expectedTotalRunTime)
	}

	expectedJumpCount := uint64(2)
	foundJumpCount := st.Averages["jump"].NumSamples
	if foundJumpCount != expectedJumpCount {
		t.Fatalf("Number of Samples calculation is incorrect, expected %d but found %d", expectedJumpCount, foundJumpCount)
	}

	expectedRunCount := uint64(2)
	foundRunCount := st.Averages["jump"].NumSamples
	if foundRunCount != expectedRunCount {
		t.Fatalf("Number of Samples calculation is incorrect, expected %d but found %d", expectedRunCount, foundRunCount)
	}
}

//test concurrecny
func TestAddAction_Concurrent(t *testing.T) {
	st := NewStats()

	for i := 0; i < 10; i++ {

		call1 := Sample{
			Action: "jump",
			Time:   10,
		}
		jsonString, _ := json.Marshal(call1)
		t.Run(fmt.Sprintf("Concurrent Test %d", i), func(t *testing.T) {
			t.Parallel()
			addErr := st.AddAction(string(jsonString))
			if addErr != nil {
				t.Error(addErr.Error())
			}
			statsJson, geterr := st.GetStats()
			if geterr != nil {
				t.Error(geterr.Error())
			}
			haveStruct := make([]SampleAverage, 1)
			json.Unmarshal([]byte(statsJson), &haveStruct)

			if haveStruct[0].Average != 10 {
				t.Fatalf("Concurrency error!")
			}

		})
	}
}

//test some edge cases, confirm error is thrown after bad json is passed
func TestAddAction_BadJson(t *testing.T) {
	badJson := "{wd;;;]}"
	st := NewStats()

	err := st.AddAction(badJson)

	if err == nil {
		t.Fatal("Accepted Bad JSON!")
	}
}

//verify we get an error when the numbers get too big
func TestAddAction_IntOverflow(t *testing.T) {
	st := NewStats()

	call1 := Sample{
		Action: "jump",
		Time:   math.MaxUint64,
	}
	call2 := Sample{
		Action: "jump",
		Time:   uint64(1),
	}
	err1 := st.addAction(call1)
	err2 := st.addAction(call2)

	if err1 != nil {
		t.Fatalf("Cannot Add uint64 max as a time")
	}
	if err2 == nil {
		t.Fatalf("TotalTime for jump exceeded maxuint64")
	}
}

//=====================Benchmark Tests====================//

//helper to make new stats with numActions number of unique actions
//used by most tests to show complexity increase for unique actions
func makeStatsWithUniqueActions(numActions int) (st Stats) {
	st = NewStats()
	for i := 0; i < numActions; i++ {
		actionName := fmt.Sprintf("Action%d", i)
		sample := Sample{
			Action: actionName,
			Time:   math.MaxUint64,
		}
		st.addAction(sample)
	}
	return st
}

//helper to make stats with numActions entries into a single action
//used by XX_highvolume and XX_lowvolume tests to show complexity in terms of unique actions
func makeStatsOneAction(numActions int) (st Stats) {
	st = NewStats()
	actionName := "action"
	for i := 0; i < numActions; i++ {
		sample := Sample{
			Action: actionName,
			Time:   uint64(i),
		}
		st.addAction(sample)
	}
	return st
}

//get benchmarks from get stats with only 10 unique actions
func BenchmarkGetStatsSmall(b *testing.B) {
	numActions := 10
	st := makeStatsWithUniqueActions(numActions)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st.GetStats()
	}

}

//benchmark get stats with 1,000,000 unique actions
func BenchmarkGetStatsMega(b *testing.B) {
	numActions := 1000000
	st := makeStatsWithUniqueActions(numActions)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st.GetStats()
	}
}

//benchmark get stats with 10 unique records to underlying call, not JSON
func BenchmarkGetStatsSmall_direct(b *testing.B) {
	st := makeStatsWithUniqueActions(10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st.getSampleAverageSlice()
	}

}

//benchmark get stats with 1,000,000 unique records to underlying call, not JSON
func BenchmarkGetStatsMega_direct(b *testing.B) {
	numActions := 1000000
	st := makeStatsWithUniqueActions(numActions)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st.getSampleAverageSlice()
	}
}

//benchmark add action with 1,000,000 unique actions
func BenchmarkAddActionMega(b *testing.B) {
	numActions := 1000000
	st := makeStatsWithUniqueActions(numActions)
	actionName := "action"
	sample := Sample{
		Action: actionName,
		Time:   math.MaxUint64,
	}
	var jsonString []byte
	jsonString, _ = json.Marshal(sample)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		st.AddAction(string(jsonString))
	}
}

//benchmark add action with 10 unique actions
func BenchmarkAddActionSmall(b *testing.B) {
	numActions := 10
	st := makeStatsWithUniqueActions(numActions)
	actionName := "action"
	sample := Sample{
		Action: actionName,
		Time:   math.MaxUint64,
	}
	var jsonString []byte
	jsonString, _ = json.Marshal(sample)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		st.AddAction(string(jsonString))
	}
}

// benchmark add action underlying call with 1,000,000 unique actions
func BenchmarkAddActionMega_direct(b *testing.B) {
	numActions := 1000000
	st := makeStatsWithUniqueActions(numActions)
	sample := Sample{
		Action: "actionName",
		Time:   math.MaxUint64,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st.addAction(sample)
	}
}

//benchmark add action underlying call with 10 unique actions
func BenchmarkAddActionSmall_direct(b *testing.B) {
	numActions := 10
	st := makeStatsWithUniqueActions(numActions)
	sample := Sample{
		Action: "actionName",
		Time:   math.MaxUint64,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		st.addAction(sample)
	}
}

// The following tests demonstrate the constant complexity of operations on a sigle action

//benchmark underlying getstats call with one unique action updated 1,000,000 times
func BenchmarkGetStats_direct_highvolume(b *testing.B) {
	numActions := 1000000
	st := makeStatsOneAction(numActions)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st.getSampleAverageSlice()
	}
}

//benchmark underlying getstats call with one unique action updated 10 times
func BenchmarkGetStats_direct_lowvolume(b *testing.B) {
	numActions := 10
	st := makeStatsOneAction(numActions)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st.getSampleAverageSlice()
	}
}

//benchmark underlying addAction call with one unique action updated 1,000,000 times
func BenchmarkAddAction_direct_highvolume(b *testing.B) {
	numActions := 1000000
	st := makeStatsOneAction(numActions)
	sample := Sample{
		Action: "action",
		Time:   uint64(1),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st.addAction(sample)
	}
}

//benchmark underlying addAction call with one unique action updated 10 times
func BenchmarkAddAction_direct_lowvolume(b *testing.B) {
	numActions := 10
	st := makeStatsOneAction(numActions)
	sample := Sample{
		Action: "action",
		Time:   uint64(1),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st.addAction(sample)
	}
}
