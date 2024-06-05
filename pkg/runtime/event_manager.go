package runtime

var _ EventManager = (*DefaultEventManager)(nil)

type DefaultEventManager struct {
	events []any
}

func (m *DefaultEventManager) EmitEvent(event any) error {
	m.events = append(m.events, event)
	return nil
}
