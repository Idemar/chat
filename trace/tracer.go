package trace

import (
	"fmt"
	"io"
)

// Tracer er grensesnittet som beskriver et objekt som er i stand til å
// spore hendelser gjennom hele koden.
type Tracer interface {
	Tracer(...any)
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Tracer(a ...any) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}
