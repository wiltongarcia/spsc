# Timetracker

### The challenge

Using only the standard library, create a Go HTTP server that on each request responds with a counter of the total number of requests that it has received during the previous 60 seconds (moving window). The server should continue to the return the correct numbers after restarting it, by persisting data to a file.

### The Approach

The following Go libraries were used to being able to reach the Goal of this challenge:
```
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
```

When the service is started it checks the presence of the data file, if the file does not exist a new one is created, if the file does exist it's imported.

From this point on the HTTP Server is started, and it's ready to receive requests. When the endpoint is called it inserts an entry in the slice with the actual timestamp (Unix Timestamp).

The response of the GET request is a sum of all requests in the past 60 seconds:

```
	limit := now - 60
	for i := len(t.dataRange) - 1; i >= 0; i-- {
		if t.dataRange[i] > limit {
			total++
		} else {
			break
		}
	}
	return total
```

When the application receives a SIGTERM, the slice is lopped converted to a binary and it's written to the file:

```
func (t *tt) Write() error {
	t.file.Truncate(0)
	for i := 0; i < len(t.dataRange); i++ {
		b := make([]byte, 8)
		binary.PutVarint(b, t.dataRange[i])
		_, err := t.file.Write(b)
		if err != nil {
			return err
		}
	}
	t.file.Close()
	return nil
```

### Why is converted to a binary?

The intention of converting the data to a binary is to save disk space, and to delimited the size of the data, using the binary I can get chunks of 8 bytes (int64) and reconstruct the state.

### The integration test

The next part of the challenge for me consists in how to do a proper integration test for the solution itself, and to test the whole behaviour.

By default, the test is running against the application every second (this parameter can be configured using `--ms=<NUMBER OF MILISECONDS>`), the test start to show the result of the requests, in the case of the application being stopped and restarted the test is able to check the amount of requests in the last 60s by using the current timestamp and the last record in the data file.

![alt text](https://github.com/wiltongarcia/spsc/blob/master/image/animation.gif "Test DEMO")

