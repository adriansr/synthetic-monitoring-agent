//go:build !windows

package checks

import (
	"context"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

// installSignalHandler installs a signal handler for SIGUSR1.
//
// The returned context's Done channel is closed if the signal is
// delivered. To make it simpler to determine if the signal was
// delivered, a value of 1 is written to the location pointed to by the
// returned int32 pointer.
//
// If the provided context's Done channel is closed before the signal is
// delivered, the signal handler is removed and the returned context's
// Done channel is closed, too. It's the callers responsibility to
// cancel the provided context if it's no longer interested in the
// signal.
func installSignalHandler(ctx context.Context) (context.Context, *int32) {
	sigCtx, cancel := context.WithCancel(ctx)

	fired := new(int32)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGUSR1)

	go func() {
		select {
		case <-sigCh:
			atomic.StoreInt32(fired, 1)
			cancel()
		case <-ctx.Done():
		}
		signal.Stop(sigCh)
	}()

	return sigCtx, fired
}
