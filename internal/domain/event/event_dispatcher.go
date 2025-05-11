package event

// Masalah: Event handling belum diimplementasikan
// Perbaikan: Implementasikan event dispatcher

// EventDispatcher menangani event handling
type EventDispatcher struct {
	handlers map[string][]EventHandler
}

// EventHandler adalah interface untuk handler event
type EventHandler interface {
	Handle(event Event)
}

// NewEventDispatcher membuat instance dispatcher baru
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

// Register mendaftarkan handler untuk tipe event tertentu
func (d *EventDispatcher) Register(eventType string, handler EventHandler) {
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// Dispatch mengirim event ke semua handler yang terdaftar
func (d *EventDispatcher) Dispatch(event Event) {
	if handlers, exists := d.handlers[event.GetType()]; exists {
		for _, handler := range handlers {
			go handler.Handle(event)
		}
	}
}
