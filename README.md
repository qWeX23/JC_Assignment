# JumpCloud Software Engineer Programming Assignment

The stats module has the ability to track the average  time for a given action. It works by calculating a running average of the time of each action as they are added. [Specification](https://github.com/qWeX23/JC_Assignment/blob/main/Software%20Engineer%20-%20Backend%20Assignment.pdf)

---
Usage:

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

[go](https://golang.org/dl/) Developed and tested on `go version go1.16.5 windows/amd64` but there is not reason to believe at time of writing that go 1.16.X for any OS would be incompatible, However, they are not tested. 

[git](https://git-scm.com/downloads) 

Open a terminal or cmd in the desired directory and run the following commands

`git clone https://github.com/qWeX23/JC_Assignment.git`

`cd JC_assignment/stats`

Compile the code 

`go build`

Testing the code

`go test -v`

To use the module standalone

`cd ../main`

Open main.go in the text editor and 

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

An alternative implementation could use a slice. This would hold each action input as a struct in the array (in memory). Upon the `getStats()` call, the program could then calculate the average of every action by making one or more passes throught the slice. This would provide the most extensability. The cost of this would be both memory usage and lookup time for averaging. 

 A map with key of action and value of slice where the slice is each action input could be used to reduce the lookup time, but have little effect on the memory usage. 

### Mutex for unsafe operations

A mutex was chosen for this implementation to ensure thread safety. This allowed for simple implementation and guaranteed thread safety for the shared memory operations. An Alternative implementation could be to use a database engine, or to investigate more into thread safe data structures in golang

---

## Assumptions

- No other statsistics would be needed from the program
- No persistance of input is necessary
- JSON is case insensative to go's standard
- The values passed for `time` will be relatively small in number of values or size of values. Because the program caluclates the total of all values per `action`, there is a possiblity of overflowing uint64 (18446744073709551615). Assuming the use case specified in the document, uint64 would have adequate headroom for the total of all specified values. The program was designed under that assumption. A mitigation for this would be to use the [cumulative moving average function](https://en.wikipedia.org/wiki/Moving_average). CMA uses the last value and the total number of values to calculate the new average. Implementing this would allow for max uint64 number of times with value that is valid uint64.
- The `time` value cannot be negative
- The average returned will be an integer approximation based on go's rounding rules
- TODO. MORE.


---
## Performance

TODO add Benchmark unit test for add and get
