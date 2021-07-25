# JumpCloud Software Engineer Programming Assignment

[![qWeX23](https://circleci.com/gh/qWeX23/JC_Assignment.svg?style=svg)](https://app.circleci.com/pipelines/github/qWeX23/JC_Assignment)
 
The stats module has the ability to track the average  time for a given action. It works by calculating a running average of the time of each action as they are added. [Specification](https://github.com/qWeX23/JC_Assignment/blob/main/Software%20Engineer%20-%20Backend%20Assignment.pdf)
 
---
## Usage
 
Stats Hello World 
```
package main
 
import (
    "github.com/qwex23/JC_Assignment/stats"
)
 
func main() {
    st := stats.NewStats()
 
    call1 := "{\"action\":\"jump\", \"time\":100}"
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
```
will output similar to 
 
`[{"action":"jump","avg":150},{"action":"run","avg":75}]`
 
 
---
## Downloading and Running the Code
 
You should have the following installed 
 
[go](https://golang.org/dl/) Developed and tested on `go version go1.16.5 windows/amd64` but there is not reason to believe at time of writing that go 1.16.X for any OS would be incompatible, However, they are not tested. go.exe will need to be updated into the path.
 
[git](https://git-scm.com/downloads) 
 
Open a terminal or cmd in the desired directory and run the following commands
 
`git clone https://github.com/qWeX23/JC_Assignment.git`
 
`cd JC_assignment/stats`
 
Testing the module
 
`go test -v`
 
To use the module standalone
 
`cd ../main`
 
Open main.go in the text editor and add your custom code to use the module
 
`go run .`
 
Will run the main program.
 
Compile the code 
 
`go build`
 
Will compile the main program in to main.exe
 
---
## Using the module
 
`go get github.com/qwex23/JC_Assignment/stats`
 
add the following import to your code 
```
import (
    "github.com/qwex23/JC_Assignment/stats"
)
```
---
 
## Design
 
### Map Vs Slice 
 
This implementation uses a map with key of the given action and value of a struct that contains the total number of actions and the running total of time units. From this we can calculate the running average by dividing the two values. 
 
Map was chosen for its low memory footprint and fast lookup.
 
An alternative implementation could use a slice. This would hold each action input as a struct in the array (in memory). Upon the `getStats()` call, the program could then calculate the average of every action by making one or more passes through the slice. This would provide the most extensibility. The cost of this would be both memory usage and lookup time for averaging. 
 
 A map with key of action and value of slice where the slice is each action input could be used to reduce the lookup time, but have little effect on the memory usage. 
 
### Mutex for unsafe operations
 
A mutex was chosen for this implementation to ensure thread safety. This allowed for simple implementation and guaranteed thread safety for the shared memory operations. An Alternative implementation could be to use a database engine, or to investigate more into thread safe data structures in golang
 
---
 
## Assumptions
 
- No other statistics would be needed from the program
- No persistence of input is necessary
- JSON is case insensitive to go's standard
- The values passed for `time` will be relatively small in number of values or size of values. Because the program calculates the total of all values per `action`, there is a possibility of overflowing uint64 (18446744073709551615). Assuming the use case specified in the document, uint64 would have adequate headroom for the total of all specified values. The program was designed under that assumption. A mitigation for this would be to use the [cumulative moving average function](https://en.wikipedia.org/wiki/Moving_average). CMA uses the last value and the total number of values to calculate the new average. Implementing this would allow for max uint64 number of times with value that is valid uint64.
- The `time` value cannot be negative
- The average returned will be an integer approximation based on go's rounding rules
- The order of the action averages in the return of `GetStats()` is unimportant. Adding a sort before we Marshal the final slice would fix this at the cost of higher runtime complexity 
 
 
---
## Performance
 
The Benchmark tests can be run from the stats directory using the command: (note on Windows this only ran successfully in powershell and not CMD)
 
`go test -bench=$.`
 
Benchmarking test provided key insight into the performance of the module. The Results showed that there was an increased compute time based on the number of unique actions in the core map. The design decision was made to use a hashmap with key of action (string) and value of the average values. The thought at the time was that there would be quick inserts in the `AddActions()` call, and this would be advantageous for a hypothetical real world use case. The speed that the map could have provided would make up for the lengthier `GetStats()`, because hypothetically there would be more adds than gets. The Results of the test back up this hypothesis. 
 
```
 go test -bench=$. 
goos: windows
goarch: amd64
pkg: github.com/qwex23/JC_Assignment/stats
cpu: AMD Ryzen 9 3900X 12-Core Processor
BenchmarkGetStatsSmall-24                         733446              1489 ns/op
BenchmarkGetStatsMega-24                              39          30190587 ns/op
BenchmarkGetStatsSmall_direct-24                 3347928               357.2 ns/op
BenchmarkGetStatsMega_direct-24                      100          13974519 ns/op
BenchmarkAddActionMega-24                         799914              1632 ns/op
BenchmarkAddActionSmall-24                        799999              1542 ns/op
BenchmarkAddActionMega_direct-24                 2580138               448.3 ns/op
BenchmarkAddActionSmall_direct-24                2927943               409.3 ns/op
BenchmarkGetStats_direct_highvolume-24          12369030                97.23 ns/op
BenchmarkGetStats_direct_lowvolume-24           12302934                96.90 ns/op
BenchmarkAddActionMega_direct_highvolume-24     46109864                25.89 ns/op
BenchmarkAddActionMega_direct_lowvolume-24      46153490                26.07 ns/op
PASS
ok      github.com/qwex23/JC_Assignment/stats   18.094s
```
 
The Benchmarks show that the performance of `GetStats()` is directly dependant on the numer of unique actions. (see GetStatsSmall vs. GetStatsMega)
 
We can also see that the amount of samples of the same action name has no effect on the performance of either `AddAction()` or `GetStats()`. (see tests with _highvolume vs _lowvolume) 
 
During development of the benchmarks, it is found that map lookup in go for string keys is at [best O(N log N )](https://stackoverflow.com/questions/29677670/what-is-the-big-o-performance-of-maps-in-golang) because of the preprocessing necessary. This could be the cause of the slightly larger processing time for the larger datasets. 

