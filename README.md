# Hertz-accessLog

Hertz middleware to get access log.

## Usage

### Start using it

Download [hertz-accessLog](https://github.com/FlameMida/accessLog) by using:

```sh
go get -u github.com/FlameMida/accessLog
```

Import following in your code:

```go
import "github.com/FlameMida/accessLog" // hertz-accessLog middleware 
```

### Quick start

Now assume you have implemented a simple api as following:

```go
func main() {
    h := server.Default()
    h.Use(accessLog.Logger())
    h.Spin()
}
```
