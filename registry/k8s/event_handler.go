package k8s

import (
	v1 "k8s.io/api/core/v1"
)

import (
	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/tools"
)

type eventHandler struct {
	eventChan chan common.RegistryChangeEvent
	key       string
	stopper   chan struct{}
}

func newK8sEventHandler(ch chan common.RegistryChangeEvent, key string, stopper chan struct{}) *eventHandler {
	return &eventHandler{
		eventChan: ch,
		key:       key,
		stopper:   stopper,
	}
}

func (kr *eventHandler) add(obj interface{}) {
	log.Debug("on add event handler called")
	p, ok := obj.(*v1.Pod)
	if !ok {
		log.Warnf("pod-informer got object %T not *v1.Pod", obj)
		return
	}
	addr, ok := p.Labels[kr.key]
	if !ok {
		return
	}
	kr.eventChan <- *common.NewRegistryChangeEvent(
		common.RegistryAddEvent,
		tools.AddrLabel2Addr(addr),
	)
}

func (kr *eventHandler) update(oldObj, newObj interface{}) {
	log.Debug("on update event handler called")
	op, ok := oldObj.(*v1.Pod)
	if !ok {
		log.Warnf("pod-informer got object %T not *v1.Pod", oldObj)
		return
	}
	np, ok := newObj.(*v1.Pod)
	if !ok {
		log.Warnf("pod-informer got object %T not *v1.Pod", newObj)
		return
	}
	if op.GetResourceVersion() == np.GetResourceVersion() {
		return
	}
	addr, ok := np.Labels[kr.key]
	if !ok {
		return
	}
	kr.eventChan <- *common.NewRegistryChangeEvent(
		common.RegistryUpdateEvent,
		tools.AddrLabel2Addr(addr),
	)
}

func (kr *eventHandler) delete(obj interface{}) {
	log.Debug("on delete event handler called")
	p, ok := obj.(*v1.Pod)
	if !ok {
		log.Warnf("pod-informer got object %T not *v1.Pod", obj)
		return
	}
	addr, ok := p.Labels[kr.key]
	if !ok {
		return
	}
	kr.eventChan <- *common.NewRegistryChangeEvent(
		common.RegistryDeleteEvent,
		tools.AddrLabel2Addr(addr),
	)
}
