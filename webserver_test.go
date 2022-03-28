package go_webserver

import (
	"fmt"
	"github.com/randlabs/go-webserver/middleware"
	"github.com/randlabs/go-webserver/request"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// -----------------------------------------------------------------------------

func TestWebServer(t *testing.T) {
	//Create server
	srvOpts := Options{
		Address: "127.0.0.1",
		Port:    3000,
	}
	srv, err := Create(srvOpts)
	if err != nil {
		t.Errorf("unable to create web server [%v]", err)
		return
	}

	// Add some middlewares
	srv.Use(middleware.DefaultCORS())
	srv.Use(middleware.DisableCacheControl())

	// Add public files to server
	var workDir string

	workDir, err = os.Getwd()
	if err != nil {
		t.Errorf("unable to get current directory [%v]", err)
		return
	}
	if !strings.HasSuffix(workDir, string(os.PathSeparator)) {
		workDir += string(os.PathSeparator)
	}

	srv.ServeFiles("/", ServerFilesOptions{
		RootDirectory: workDir + "testdata/public",
	})

	// Add a dummy api function
	srv.POST("/api/version", renderApiVersion)

	// Add also profile output
	srv.AddProfilerHandlers("/debug/", nil)

	// Start server
	err = srv.Start()
	if err != nil {
		t.Errorf("unable to start web server [%v]", err)
		return
	}

	// Open default browser
	openBrowser("http://" + srvOpts.Address + ":" + strconv.Itoa(int(srvOpts.Port)) + "/")
	//openBrowser("http://" + srvOpts.Address + ":" + strconv.Itoa(int(srvOpts.Port)) + "/debug/")

	// Wait for CTRL+C
	fmt.Println("Server running. Press CTRL+C to stop.")

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	select {
	case <-c:
	case <-time.After(5 * time.Minute):
	}
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

func renderApiVersion(req *request.RequestContext) error {
	req.WriteString(`{ "version": "1.0.0" }`)
	return nil
}
