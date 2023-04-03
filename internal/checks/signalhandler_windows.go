//go:build windows

package checks

import (
	"context"
)

var neverFired int32 = 0

// installSignalHandler installs a signal handler for SIGUSR1.
//
// This is a no-op for Windows.
func installSignalHandler(ctx context.Context) (context.Context, *int32) {
	return ctx, &neverFired
}
