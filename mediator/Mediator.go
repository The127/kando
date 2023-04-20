package mediator

import (
	"context"
	"kando-backend/log"
	"kando-backend/utils"
	"reflect"
)

type Mediator struct {
	behaviours    map[reflect.Type][]behaviourInfo
	handlers      map[reflect.Type]handlerInfo
	eventHandlers map[reflect.Type][]eventHandlerInfo
}

func NewMediator() *Mediator {
	return &Mediator{
		handlers:   map[reflect.Type]handlerInfo{},
		behaviours: map[reflect.Type][]behaviourInfo{},
	}
}

type EventHandlerFunc[TEvent any] func(TEvent, context.Context) error

type eventHandlerInfo struct {
	eventType        reflect.Type
	eventHandlerFunc func(any, context.Context) error
}

func RegisterEventHandler[TEvent any](m *Mediator, eventHandler EventHandlerFunc[TEvent]) {
	eventType := utils.TypeOf[TEvent]()

	eventHandlers, ok := m.eventHandlers[eventType]
	if !ok {
		eventHandlers = []eventHandlerInfo{}
	}

	eventHandlers = append(eventHandlers, eventHandlerInfo{
		eventType: eventType,
		eventHandlerFunc: func(event any, ctx context.Context) error {
			return eventHandler(event.(TEvent), ctx)
		},
	})

	m.eventHandlers[eventType] = eventHandlers
}

type Next func()

type BehaviourFunc[TRequest any] func(TRequest, context.Context, Next)

type behaviourInfo struct {
	requestType   reflect.Type
	behaviourFunc func(any, context.Context, Next)
}

func RegisterBehaviour[TRequest any](m *Mediator, behaviour BehaviourFunc[TRequest]) {
	requestType := utils.TypeOf[TRequest]()

	behaviours, ok := m.behaviours[requestType]
	if !ok {
		behaviours = []behaviourInfo{}
	}

	behaviours = append(behaviours, behaviourInfo{
		requestType: requestType,
		behaviourFunc: func(request any, ctx context.Context, next Next) {
			behaviour(request.(TRequest), ctx, next)
		},
	})

	m.behaviours[requestType] = behaviours
}

type HandlerFunc[TRequest any, TResponse any] func(TRequest, context.Context) (TResponse, error)

func RegisterHandler[TRequest any, TResponse any](m *Mediator, handler HandlerFunc[TRequest, TResponse]) {
	m.handlers[utils.TypeOf[TRequest]()] = handlerInfo{
		requestType:  utils.TypeOf[TRequest](),
		responseType: utils.TypeOf[TResponse](),
		handlerFunc: func(request any, ctx context.Context) (any, error) {
			return handler(request.(TRequest), ctx)
		},
	}
}

type handlerInfo struct {
	requestType  reflect.Type
	responseType reflect.Type
	handlerFunc  func(any, context.Context) (any, error)
}

func SendEvent[TEvent any](m *Mediator, event TEvent, ctx context.Context) error {
	eventType := utils.TypeOf[TEvent]()

	eventHandlers, ok := m.eventHandlers[eventType]
	if !ok {
		log.Logger.Debugf("Could not find any event handlers for %s", eventType.Name())
		return nil
	}

	for _, eventHandler := range eventHandlers {
		err := eventHandler.eventHandlerFunc(event, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func Send[TResponse any](m *Mediator, request any, ctx context.Context) (TResponse, error) {
	requestType := reflect.TypeOf(request)

	var step Next

	handler, ok := m.handlers[requestType]
	if !ok {
		log.Logger.Fatalf("Could not find any registered handler for %s", requestType.Name())
	}

	responseType := utils.TypeOf[TResponse]()
	if handler.responseType != responseType {
		log.Logger.Fatalf("wrong response type %s was used for request %s, expected response type %s",
			responseType.Name(), requestType.Name(), handler.responseType.Name())
	}

	var response any
	var err error

	step = func() {
		response, err = handler.handlerFunc(request, ctx)
	}

	var behaviours []behaviourInfo

	for key := range m.behaviours {
		if requestType.AssignableTo(key) {
			behaviours = append(behaviours, m.behaviours[key]...)
		}
	}

	for i := len(behaviours) - 1; i >= 0; i-- {
		behaviour := behaviours[i]
		prev := step
		step = func() {
			behaviour.behaviourFunc(request, ctx, prev)
		}
	}

	step()

	return response.(TResponse), err
}
