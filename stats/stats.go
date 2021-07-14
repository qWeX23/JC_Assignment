package stats

import (
	"encoding/json"
	"sync"
)

type Sample struct {
	Action string `json:"action"`
	Time   uint64 `json:"time"`
}

type SampleAverage struct {
	Action  string `json:"action"`
	Average uint64 `json:"avg"`
}

type Average struct {
	NumSamples uint64 `json:"numSamples"`
	TotalTime  uint64 `json:"totalTime"`
}

type Stats struct {
	Averages map[string]*Average
	mu       sync.Mutex
}

func NewStats() Stats {
	return Stats{
		Averages: make(map[string]*Average),
	}
}

func printStats(s *Stats) {
	for a, v := range s.Averages {
		println(v.TotalTime / v.NumSamples)
		print(a)
	}
}

//todo err
func (s *Stats) GetStats() string {
	jsonString, _ := json.Marshal(s.getSampleAverageSlice())

	return string(jsonString)

}

//todo return error and threadsafey
func (s *Stats) getSampleAverageSlice() []SampleAverage {
	AveragesSlice := make([]SampleAverage, 0)

	s.mu.Lock()
	defer s.mu.Unlock()
	for action, average := range s.Averages {
		sampleAverage := SampleAverage{
			Action:  action,
			Average: average.TotalTime / average.NumSamples,
		}
		AveragesSlice = append(AveragesSlice, sampleAverage)

	}
	return AveragesSlice
}

func (s *Stats) AddAction(sampleString string) {
	var sample Sample
	//todo add some validation and error checking
	err := json.Unmarshal([]byte(sampleString), &sample)

	if err != nil {
		panic(err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Averages[sample.Action] == nil {
		s.Averages[sample.Action] = &Average{
			NumSamples: 1,
			TotalTime:  sample.Time,
		}
	} else {
		s.Averages[sample.Action].TotalTime += sample.Time
		s.Averages[sample.Action].NumSamples += 1
	}
}
