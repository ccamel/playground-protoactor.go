package eventsourcingv1

type Event interface {
	GetBase() *EventBase
}
