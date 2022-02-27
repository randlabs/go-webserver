package go_webserver

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"testing"
)

// -----------------------------------------------------------------------------

func TestWebServer(t *testing.T) {
	srvOpts := Options{
		Address: "127.0.0.1",
		Port:    3000,
	}
	srv, err := Create(srvOpts)
	if err != nil {
		t.Errorf("unable to create web server [%v]", err)
		return
	}

	srv.AddProfilerHandlers("/debug/", nil)

	// Start server
	err = srv.Start()
	if err != nil {
		t.Errorf("unable to start web server [%v]", err)
		return
	}

	// Open default browser
	openBrowser("http://" + srvOpts.Address + ":" + strconv.Itoa(int(srvOpts.Port)) + "/debug/")

	fmt.Println("Server running. Press CTRL+C to stop.")

	// Wait for CTRL+C
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("Shutting down...")

	// Stop web server
	srv.Stop()
}

// -----------------------------------------------------------------------------

func openBrowser(url string) {
	switch runtime.GOOS {
	case "linux":
		_ = exec.Command("xdg-open", url).Start()
	case "windows":
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		_ = exec.Command("open", url).Start()
	}
}
