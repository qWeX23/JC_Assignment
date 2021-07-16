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

//test some edge cases
func TestAddAction_BadJson(t *testing.T) {
	badJson := "{wd;;;]}"
	st := NewStats()

	err := st.AddAction(badJson)

	if err == nil {
		t.Fatal("Accepted Bad JSON!")
	}
}

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
