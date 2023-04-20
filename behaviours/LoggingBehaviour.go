package behaviours

import (
	"context"
	"kando-backend/log"
	"kando-backend/mediator"
	"reflect"
)

func LoggingBehaviour(request any, ctx context.Context, next mediator.Next) {
	requestType := reflect.TypeOf(request)
	log.Logger.Debugf("Executing %v", requestType.Name())
	next()
	log.Logger.Debugf("Executed %v", requestType.Name())
}
