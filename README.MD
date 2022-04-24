This is simple server which can be used as a mock server for frontend development.

Example:
```
package main

import (
	. "github.com/paweljvm/gojsonserver"
)

func main() {
	js := NewLocalJsonServer(4444, []RequestHandler{
		GET("/api", "api.json", 200, 300),
	})
	js.Start()
}
```