package eventsourcingv1

type Command interface {
	GetBase() *CommandBase
}
