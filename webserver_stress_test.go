package go_webserver_test

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	webserver "github.com/randlabs/go-webserver/v2"
	"github.com/randlabs/go-webserver/v2/helpers_test"
	"github.com/randlabs/go-webserver/v2/middleware"
)

// -----------------------------------------------------------------------------

func TestWebServerStress(t *testing.T) {
	//Create server
	srv := helpers_test.RunWebServer(t, func(srv *webserver.Server) error {
		// Add some middlewares
		srv.Use(middleware.DisableClientCache())

		// Done
		return nil
	})
	defer srv.Stop()

	// Start request workers and main context
	successCounter := int32(0)
	failCounter := int32(0)

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	for idx := 0; idx < runtime.GOMAXPROCS(0); idx++ {
		wg.Add(1)

		go func(ctx context.Context) {
			for {
				_, _, err := helpers_test.QueryApiVersion(false, nil, nil, []int{200})
				if err == nil {
					atomic.AddInt32(&successCounter, 1)
				} else {
					atomic.AddInt32(&failCounter, 1)
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

	t.Logf("Processed %v requests (%v succeeded) in %d seconds",
		atomic.LoadInt32(&successCounter)+atomic.LoadInt32(&failCounter),
		atomic.LoadInt32(&successCounter), 5)
}
