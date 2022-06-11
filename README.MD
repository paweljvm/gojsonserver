This is simple server which can be used as a mock server for frontend development.

Example:
```
package main

import (
	. "github.com/paweljvm/gojsonserver"
)

func main() {
	js := NewLocalJsonServer(4444, []RequestHandler{
		Get("/api", "api.json", 200, 300),
		GetProvider("/custom-logic", func(req *http.Request) {
			 if req.ContentLength > 0  {
				return "json-not-empty.json"
			 } else {
				return "json-empty.json"
			 }
		})
	})
	js.Start()
}
```