package trace

import (
	"fmt"
	"io"
)

// Tracer er grensesnittet som beskriver et objekt som er i stand til å
// spore hendelser gjennom hele koden.
type Tracer interface {
	Trace(...interface{})
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

// Off oppretter en Tracer som ignorerer kall til Trace.
func Off() Tracer {
	return &nilTracer{}
}
