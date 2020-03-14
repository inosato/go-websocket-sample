package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("tracer is nil")
	} else {
		tracer.Trace("hello trace package")
		if buf.String() != "hello trace package\n" {
			t.Errorf("unexpect string: %s", buf.String())
		}
	}
}

func TestOff(t *testing.T) {
	var silentReacer Tracer = Off()
	silentReacer.Trace("data")
}
