package util

/*
func AssertEventEmmited(t *testing.T, runtime runtime.RuntimeManager, event any) {
	var ev types.Event

	switch cast := event.(type) {
	case types.Event:
		ev = cast
	default:
		var err error
		ev, err = types.TypedEventToEvent(cast.(proto.Message))
		if err != nil {
			panic(err)
		}
	}

	for _, e := range ctx.EventManager().Events() {
		if reflect.DeepEqual(e, ev) {
			return
		}
	}
	t.Fatalf("EventManager did not emit wanted event: want %v", event)
}

*/
