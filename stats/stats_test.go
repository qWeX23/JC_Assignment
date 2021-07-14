package stats

import "testing"

func TestAddAction(t *testing.T) {
	//Arrange
	st := NewStats()
	call1 := "{\"action\":\"jump\", \"time\":100}"
	call2 := "{\"action\":\"run\", \"time\":75}"
	call3 := "{\"action\":\"jump\", \"time\":200}"
	//Act
	st.AddAction(call1)
	st.AddAction(call2)
	st.AddAction(call3)
	//Assert
	if st.GetStats() != "[{\"action\":\"jump\",\"avg\":150},{\"action\":\"run\",\"avg\":75}]" {
		t.Fatal("Base Test Case failed")
	}
}

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
