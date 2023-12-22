package broadcast

import (
	"bytes"
	"math"
	"pan/core"
	"time"
)

type Controller struct {
	service *Service
	network Net
}

// Handle ...
func (ctrl *Controller) Handle(ctx Context, next core.Next) {
	method := ctx.Method()
	if bytes.Equal([]byte("alive"), method) {
		ctrl.Alive(ctx, next)
		return
	}

	if bytes.Equal([]byte("dead"), method) {
		ctrl.Dead(ctx, next)
		return
	}
	next()
}

// BroadcastAlive ...
func (ctrl *Controller) BroadcastAlive() {
	payload, err := ctrl.service.GenerateAliveMessage()
	if err != nil {
		panic(err)
	}
	err = dispatch([]byte("alive"), payload, ctrl.network, 5)
	if err != nil {
		panic(err)
	}
}

// Alive ...
func (ctrl *Controller) Alive(ctx Context, next core.Next) {
	err := ctrl.service.RecvAliveMessage(ctx.Addr(), ctx.Body())
	if err != nil {
		panic(err)
	}
}

// BroadcastDead ...
func (ctrl *Controller) BroadcastDead() {
	payload, err := ctrl.service.GenerateDeadMessage()
	if err != nil {
		panic(err)
	}
	err = dispatch([]byte("dead"), payload, ctrl.network, 2)
	if err != nil {
		panic(err)
	}

}

// Dead ...
func (ctrl *Controller) Dead(ctx Context, next core.Next) {
	err := ctrl.service.RecvDeadMessage(ctx.Addr(), ctx.Body())
	if err != nil {
		panic(err)
	}
}

// NewController ...
func NewController(service *Service, network Net) *Controller {
	ctrl := new(Controller)
	ctrl.service = service
	ctrl.network = network
	return ctrl
}

// dispatch ...
func dispatch(method, body []byte, n Net, times int) (err error) {
	for i := 0; i < times; i++ {
		err = Dispatch(method, body, n)
		if err != nil {
			break
		}
		num := 1500 * math.Pow(1.5, float64(i))
		ms := int64(num)
		time.Sleep(time.Millisecond * time.Duration(ms))
	}
	return
}
