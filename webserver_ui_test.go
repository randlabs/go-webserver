//go:build ui_test

package go_webserver_test

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	webserver "github.com/randlabs/go-webserver/v2"
	"github.com/randlabs/go-webserver/v2/helpers_test"
	"github.com/randlabs/go-webserver/v2/middleware"
)

// -----------------------------------------------------------------------------

func TestWebServerUI(t *testing.T) {
	//Create server
	srv := helpers_test.RunWebServer(t, func(srv *webserver.Server) error {
		// Add some middlewares
		srv.Use(middleware.DefaultCORS())
		srv.Use(middleware.DisableClientCache())

		// Add public files to server
		err := srv.ServeFiles("/", webserver.ServerFilesOptions{
			RootDirectory: helpers_test.GetWorkingDirectory(t) + "testdata/public",
		})
		if err != nil {
			return err
		}

		// Done
		return nil
	})
	defer srv.Stop()

	// Open default browser
	helpers_test.OpenBrowser()

	// Wait for CTRL+C
	fmt.Println("Server running. Press CTRL+C to stop.")

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-c:
	case <-time.After(5 * time.Minute):
	}
	fmt.Println("Shutting down...")
}
