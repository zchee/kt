// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	"fmt"
	"sync"

	"github.com/go-logr/logr"
	color "github.com/zchee/color/v2"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	ctrlevent "sigs.k8s.io/controller-runtime/pkg/event"
	ctrlhandler "sigs.k8s.io/controller-runtime/pkg/handler"
	ctrlpredicate "sigs.k8s.io/controller-runtime/pkg/predicate"
	ctrlreconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/zchee/kt/pkg/io"
)

const (
	namespaceFmt       = "%s %s%s » %s\n" // (+|-) Namespace/PodName » ContainerName
	nonNamespaceFmt    = "%s %s » %s\n"   // (+|-) PodName » ContainerName
	namespaceSeparator = "/"

	createMark = "+"
	deleteMark = "-"
)

// PredicatePodEventFilter filters events before they are provided to handler.EventHandlers.
type PredicatePodEventFilter struct {
	ioStreams    io.Streams
	ioMu         sync.Mutex
	log          logr.Logger
	states       *sync.Map
	isNamespaced bool
}

var _ ctrlpredicate.Predicate = (*PredicatePodEventFilter)(nil)

// Create implements predicate.Predicate.
func (e *PredicatePodEventFilter) Create(event ctrlevent.CreateEvent) bool {
	pod, ok := event.Object.(*corev1.Pod)
	if !ok {
		return true
	}

	mark := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	p, c := findColors(pod.Name)

	printFunc := func(pod *corev1.Pod, container *corev1.Container) {
		format := nonNamespaceFmt
		args := []interface{}{mark(createMark), p.SprintfFunc()(pod.Name), c.SprintfFunc()(container.Name)}

		if e.isNamespaced {
			format = namespaceFmt
			args = append([]interface{}{args[0], p.SprintfFunc()(pod.Namespace + namespaceSeparator)}, args[1:]...)
		}

		fmt.Fprintf(e.ioStreams.Out, format, args...)
	}

	// if pod.Status.Phase == corev1.PodRunning {
	for i, s := range pod.Status.InitContainerStatuses {
		e.log.Info("InitContainerStatuses", "s.State", s.State)
		if s.State.Running != nil {
			printFunc(pod, &pod.Spec.InitContainers[i])
		}
	}
	for i, s := range pod.Status.ContainerStatuses {
		// e.log.Info("ContainerStatuses", fmt.Sprintf("pod.Spec.ContainerStatuses[%d]", i), pod.Status.ContainerStatuses[i])
		e.log.Info("ContainerStatuses", "s.State", s.State)
		if s.State.Running != nil {
			printFunc(pod, &pod.Spec.Containers[i])
		}
	}
	// }

	return true
}

// Delete implements predicate.Predicate.
func (e *PredicatePodEventFilter) Delete(event ctrlevent.DeleteEvent) bool {
	pod, ok := event.Object.(*corev1.Pod)
	if !ok {
		return true
	}

	mark := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	p, c := findColors(pod.Name)

	printFunc := func(pod *corev1.Pod, container *corev1.Container) {
		format := nonNamespaceFmt
		args := []interface{}{mark(deleteMark), p.SprintfFunc()(pod.Name), c.SprintfFunc()(container.Name)}

		if e.isNamespaced {
			format = namespaceFmt
			args = append([]interface{}{args[0], p.SprintfFunc()(pod.Namespace + namespaceSeparator)}, args[1:]...)
		}

		e.ioMu.Lock()
		fmt.Fprintf(e.ioStreams.Out, format, args...)
		e.ioMu.Unlock()
	}

	for i, s := range pod.Status.InitContainerStatuses {
		// e.log.Info(fmt.Sprintf("pod.Spec.InitContainers[i]: %s", pod.Spec.InitContainers[i].Name), "pod.Status", pod.Status)
		e.log.Info("InitContainerStatuses", "s.State", s.State)
		if s.State.Terminated != nil {
			printFunc(pod, &pod.Spec.InitContainers[i])
		}
	}
	for i, s := range pod.Status.ContainerStatuses {
		// e.log.Info(fmt.Sprintf("pod.Spec.ContainerStatuses[i]: %s", pod.Status.ContainerStatuses[i].Name), "pod.Status", pod.Status)
		e.log.Info("ContainerStatuses", "s.State", s.State)
		if s.State.Terminated != nil {
			printFunc(pod, &pod.Spec.Containers[i])
		}
	}

	return true
}

// Update implements predicate.Predicate.
func (e *PredicatePodEventFilter) Update(event ctrlevent.UpdateEvent) bool {
	return true
}

// Generic implements predicate.Predicate.
func (e *PredicatePodEventFilter) Generic(event ctrlevent.GenericEvent) bool {
	return true
}

// PodEventHandler enqueues reconcile.Requests in response to only of pods events.
type PodEventHandler struct {
	ioStreams    io.Streams
	ioMu         sync.Mutex
	log          logr.Logger
	states       *sync.Map
	isNamespaced bool
}

var _ ctrlhandler.EventHandler = (*PodEventHandler)(nil)

// Create implements handler.EventHandler.
func (e *PodEventHandler) Create(event ctrlevent.CreateEvent, q workqueue.RateLimitingInterface) {
	if event.Meta == nil {
		e.log.Error(nil, "CreateEvent received with no metadata", "event", event)
		return
	}

	q.Add(ctrlreconcile.Request{NamespacedName: types.NamespacedName{
		Name:      event.Meta.GetName(),
		Namespace: event.Meta.GetNamespace(),
	}})
}

// Update implements handler.EventHandler.
func (e *PodEventHandler) Update(event ctrlevent.UpdateEvent, q workqueue.RateLimitingInterface) {
	if event.MetaOld == nil {
		e.log.Error(nil, "UpdateEvent received with no old metadata", "event", event)
	}
	q.Add(ctrlreconcile.Request{NamespacedName: types.NamespacedName{
		Name:      event.MetaOld.GetName(),
		Namespace: event.MetaOld.GetNamespace(),
	}})

	if event.MetaNew == nil {
		e.log.Error(nil, "UpdateEvent received with no new metadata", "event", event)
	}
	q.Add(ctrlreconcile.Request{NamespacedName: types.NamespacedName{
		Name:      event.MetaNew.GetName(),
		Namespace: event.MetaNew.GetNamespace(),
	}})
}

// Delete implements handler.EventHandler.
func (e *PodEventHandler) Delete(event ctrlevent.DeleteEvent, q workqueue.RateLimitingInterface) {
	if event.Meta == nil {
		e.log.Error(nil, "DeleteEvent received with no metadata", "event", event)
		return
	}

	q.Add(ctrlreconcile.Request{NamespacedName: types.NamespacedName{
		Name:      event.Meta.GetName(),
		Namespace: event.Meta.GetNamespace(),
	}})

	q.Add(ctrlreconcile.Request{NamespacedName: types.NamespacedName{
		Name:      event.Meta.GetName(),
		Namespace: event.Meta.GetNamespace(),
	}})
}

// Generic implements handler.EventHandler.
func (e *PodEventHandler) Generic(event ctrlevent.GenericEvent, q workqueue.RateLimitingInterface) {
	if event.Meta == nil {
		e.log.Error(nil, "GenericEvent received with no metadata", "event", event)
		return
	}
	q.Add(ctrlreconcile.Request{NamespacedName: types.NamespacedName{
		Name:      event.Meta.GetName(),
		Namespace: event.Meta.GetNamespace(),
	}})
}
