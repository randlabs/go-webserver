# go-webserver

HTTP web server library for Go based on FastHttp

##### NOTE:

* This is a fork of the original [RandLabs.IO's webserver library](https://github.com/randlabs/go-webserver).
  May contain some modified functionality.

## Usage with example

```golang
package example

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	webserver "github.com/mxmauro/go-webserver/v2"
	"github.com/mxmauro/go-webserver/v2/middleware"
)

type testApiOutput struct {
	Status  string `json:"status"`
}

func main() {
	// Options struct has all the documentation
	srvOpts := webserver.Options{
		Address: "127.0.0.1",
		Port:    3000,
	}
	srv, err := webserver.Create(srvOpts)
	if err != nil {
		fmt.Printf("unable to create web server [%v]\n", err)
		return
	}

	// Add some middlewares
	srv.Use(middleware.DefaultCORS())

	// Setup a route
	srv.GET("/test", getTestApi)

	// Start server
	err = srv.Start()
	if err != nil {
		fmt.Printf("unable to start web server [%v]\n", err)
		return
	}

	fmt.Println("Server running. Press CTRL+C to stop.")

	// Wait for CTRL+C
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("Shutting down...")

	// Stop web server
	srv.Stop()
}

func getTestApi(req *webserver.RequestContext) error {
	// Prepare output
	output := testApiOutput{
		Status: "all systems operational",
    }

	// Encode and send output
	req.WriteJSON(output)
	req.Success()
	return nil
}
```

## License
See `LICENSE` file for details.
