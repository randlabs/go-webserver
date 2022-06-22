package go_webserver_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	webserver "github.com/randlabs/go-webserver"
	"github.com/randlabs/go-webserver/middleware"
	"github.com/randlabs/go-webserver/request"
)

// IMPORTANT NOTE: Tests are intended to be executed separately.

// -----------------------------------------------------------------------------

type versionApiOutput struct {
	Version string `json:"version"`
}

// -----------------------------------------------------------------------------

func TestWebServer(t *testing.T) {
	//Create server
	srvOpts := webserver.Options{
		Address: "127.0.0.1",
		Port:    3000,
	}
	srv, err := webserver.Create(srvOpts)
	if err != nil {
		t.Errorf("unable to create web server [%v]", err)
		return
	}

	// Add some middlewares
	srv.Use(middleware.DefaultCORS())
	srv.Use(middleware.DisableClientCache())

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

	_ = srv.ServeFiles("/", webserver.ServerFilesOptions{
		RootDirectory: workDir + "testdata/public",
	})

	// Add a dummy api function
	srv.POST("/api/version", renderApiVersion)

	// Add also profile output
	srv.ServeDebugProfiles("/debug/")

	// Start server
	err = srv.Start()
	if err != nil {
		t.Errorf("unable to start web server [%v]", err)
		return
	}

	// Open default browser
	openBrowser("http://" + srvOpts.Address + ":" + strconv.Itoa(int(srvOpts.Port)) + "/")

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

func TestWebServerStress(t *testing.T) {
	//Create server
	srvOpts := webserver.Options{
		Address: "127.0.0.1",
		Port:    3000,
	}
	srv, err := webserver.Create(srvOpts)
	if err != nil {
		t.Errorf("unable to create web server [%v]", err)
		return
	}

	// Add some middlewares
	srv.Use(middleware.DisableClientCache())

	// Add a dummy api function
	srv.GET("/api/version", renderApiVersion)

	// Start server
	err = srv.Start()
	if err != nil {
		t.Errorf("unable to start web server [%v]", err)
		return
	}

	var counter int32

	// Start request workers and main context
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	for idx := 0; idx < runtime.GOMAXPROCS(0); idx++ {
		wg.Add(1)

		go func(ctx context.Context) {
			url := "http://" + srvOpts.Address + ":" + fmt.Sprint(srvOpts.Port) + "/api/version"

			for {
				req, err2 := http.NewRequest(http.MethodGet, url, nil)
				if err2 == nil {
					var resp *http.Response

					reqCtx, reqCtxCancel := context.WithTimeout(ctx, 5*time.Second)

					resp, _ = http.DefaultClient.Do(req.WithContext(reqCtx))
					if resp != nil {
						_ = resp.Body.Close()
					}

					reqCtxCancel()

					atomic.AddInt32(&counter, 1)
				}

				select {
				case <-ctx.Done():
					wg.Done()
					return
				default:
				}
			}
		}(ctx)
	}

	// Run
	time.Sleep(5 * time.Second)

	// Stop workers
	cancel()
	wg.Wait()

	t.Logf("Processed %v requests", atomic.LoadInt32(&counter))

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
	output := versionApiOutput{
		Version: "1.0.0",
	}
	req.WriteJSON(output)
	req.Success()
	return nil
}
