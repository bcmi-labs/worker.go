package worker

// Trace is an object that keep tracks of the duration of an action.
// When the object is created a timer starts, and whenever Mark is called
// the object will record the duration of the action described
type Trace interface {
	Mark(action string)
}

// Tracer is an object that can spawn traces. Scope is a generic string
// to assign to the trace in order to organize them
type Tracer interface {
	New(scope string) Trace
}
