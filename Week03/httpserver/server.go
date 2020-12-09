package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	ctx, cancel := context.WithCancel(context.Background())
	eg, egCtx := errgroup.WithContext(context.Background())

	eg.Go(newHTTPServer(ctx, &wg, "Ironman", ":8080", ironmanHandler))
	eg.Go(newHTTPServer(ctx, &wg, "Thor", ":8081", thorHandler))
	eg.Go(newHTTPServer(ctx, &wg, "Cap.A", ":8082", captainAmericanHandler))

	go func() {
		<-egCtx.Done()
		cancel()
	}()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		cancel()
	}()

	if err := eg.Wait(); err != nil {
		fmt.Printf("error in the server goroutines: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("everything closed successfully")
}

func newHTTPServer(
	ctx context.Context,
	wg *sync.WaitGroup,
	name, addr string,
	handler http.HandlerFunc,
) func() error {
	return func() error {
		mux := http.NewServeMux()
		mux.HandleFunc("/", handler)
		server := &http.Server{Addr: addr, Handler: mux}
		errChan := make(chan error, 1)

		go func() {
			<-ctx.Done()
			shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(shutDownCtx); err != nil {
				errChan <- fmt.Errorf("error shutting down the %s server: %w", name, err)
			}
			fmt.Printf("the %s server is closed\n", name)
			close(errChan)
			wg.Done()
		}()

		fmt.Printf("the %s server is starting\n", name)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return fmt.Errorf("error starting the %s server: %w", name, err)
		}
		fmt.Printf("the %s server is closing\n", name)
		err := <-errChan
		wg.Wait()
		return err
	}
}

func ironmanHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`Hello, Tony Stark!`))
}

func thorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`Hello, Son of Odin!`))
}

func captainAmericanHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`Hello, Steve Rogers!`))
}
