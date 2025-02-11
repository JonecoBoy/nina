package event

import (
	"errors"
	"time"
)

var ErrHandlerAlreadyRegistered = errors.New("handler already registered")

// evento em si
type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() interface{} // dados do evento

}

// operacao que  executa o evento
type EventHandlerInterface interface {
	Handle(event EventInterface) // para rodar uma operacao de evento precisa ter o evento comigo
}

// quem gerencia os eventos e dispacha as operações

type EventDispatcherInterface interface {
	Register(eventName string, handler EventHandlerInterface) error // quando o evento acontecer execute o handler.
	Dispatch(event EventInterface) error                            // faz com que o evento aconteça e faz com que os handlers sejam executados
	Remove(eventName string, handler EventHandlerInterface) error   // remove um evento da fila
	Has(eventName string, handler EventHandlerInterface) bool       // verifica se tem um event name com esse handler para gente
	Clear() error                                                   // vai limpar o event dispatcher matando todos os eventos registrados lá dentro.
}

type EventDispatcher struct {
	// map where we will have the key/name as string (event name) and it can have several EventHandlerInterface
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	// if it already exists, return an error
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	// check if the event is already registered
	if _, ok := ed.handlers[eventName]; ok {
		// if ok is true, iterate through all the event handlers
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}
	// if not registered, add to the map.
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

func (ed *EventDispatcher) Clear() {
	// simply recreate the struct with zeroed values
	ed.handlers = make(map[string][]EventHandlerInterface)
}

func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	// must go through all handlers within all events and see if the handler already exists
	if _, ok := ed.handlers[eventName]; ok {
		// ok{} means if ok=true
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return true
			}
		}
	}
	return false
}

func (ev *EventDispatcher) Dispatch(event EventInterface) error {
	// if there are handlers registered with this name, enter each one of them and execute the Handle method passing the event that was called. I know it has the handle method because they implement the interface
	if handlers, ok := ev.handlers[event.GetName()]; ok {
		for _, handler := range handlers {
			// we add the goroutine here so it won't be synchronous, that is, wait for one handler to execute the event for the other to execute
			go handler.Handle(event)
		}
	}
	return nil
}
