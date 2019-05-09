package reactive

import (
	"reflect"
	"testing"
)

func Test_bufferSubject_Close(t *testing.T) {
	subject := NewBufferSubject(1)
	subject.Subscribe(func() {})

	subject.Close()

	casted := subject.(*bufferSubject)

	if casted.LastValues != nil {
		t.Error("LastValues does not equal nil")
	}

	if !reflect.DeepEqual(casted.Subscriptions, make(map[Subscription]interface{})) {
		t.Error("Subscriptions does not equal empty list")
	}

}

func TestBufferSubject_Pipe(t *testing.T) {
	subject := NewBufferSubject(1).(*bufferSubject)

	if subject != subject.Pipe() {
		t.Error("Empty pipe is different from original")
	}

	if subject == subject.Pipe(Take(1)) {
		t.Error("Take pipe resulted in original subject")
	}

	if subject != subject.Pipe(nil) {
		t.Error("Nil pipe is different from original")
	}
}
