// See the LICENSE file for license details.

package go_webserver

// -----------------------------------------------------------------------------

type loggerBridge struct {
}

// -----------------------------------------------------------------------------

func newSilentLogger() *loggerBridge {
	return &loggerBridge{}
}

func (l *loggerBridge) Printf(_ string, _ ...interface{}) {
}
