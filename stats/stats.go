package stats

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sync"
)

//Model for input object, has json names included
type Sample struct {
	Action string `json:"action"`
	Time   uint64 `json:"time"`
}

//model for output object, has json names included
type SampleAverage struct {
	Action  string `json:"action"`
	Average uint64 `json:"avg"`
}

//internal model for calculating average
type Average struct {
	NumSamples uint64 `json:"numSamples"`
	TotalTime  uint64 `json:"totalTime"`
}

//primary struct for use in calculating averages. SS
type Stats struct {
	Averages map[string]*Average
	mu       sync.Mutex
}

//creates new stats struct
func NewStats() Stats {
	return Stats{
		Averages: make(map[string]*Average),
	}
}

//returns all of the averages as a json array
func (s *Stats) GetStats() (string, error) {
	sliceAvg := s.getSampleAverageSlice()
	if recover() != nil {
		return "", fmt.Errorf("error getting sample array")
	}
	jsonString, err := json.Marshal(sliceAvg)
	return string(jsonString), err
}

//traverses the stats map, calulates the averages and returns them as an array
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

//adds the json sample to the stats object
func (s *Stats) AddAction(sampleString string) error {
	var sample Sample
	//todo add some validation
	err := json.Unmarshal([]byte(sampleString), &sample)
	if err != nil {
		return errors.New("JSON String is invalid-> " + err.Error())
	}
	s.addAction(sample)
	return err
}

//takes the sample and adds to the average struct of the corresponding action
// creates new action in stats if non is available
func (s *Stats) addAction(sample Sample) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Averages[sample.Action] == nil {
		s.Averages[sample.Action] = &Average{
			NumSamples: 1,
			TotalTime:  sample.Time,
		}
	} else {
		if math.MaxUint64-s.Averages[sample.Action].TotalTime < sample.Time {
			return fmt.Errorf("adding Sample with time %d will overflow unint64 with current time total for %s as %d", sample.Time, sample.Action, s.Averages[sample.Action].TotalTime)
		}

		s.Averages[sample.Action].TotalTime += sample.Time
		s.Averages[sample.Action].NumSamples += 1
	}
	return nil
}
