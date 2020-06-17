package ctx

import (
	"context"
	"io"

	orascontext "github.com/deislabs/oras/pkg/context"
	"github.com/sirupsen/logrus"
)

// Context retrieves a fresh context.
// disable verbose logging coming from ORAS (unless debug is enabled)
func Context(out io.Writer, debug bool) context.Context {
	if !debug {
		return orascontext.Background()
	}
	ctx := orascontext.WithLoggerFromWriter(context.Background(), out)
	orascontext.GetLogger(ctx).Logger.SetLevel(logrus.DebugLevel)
	return ctx
}
