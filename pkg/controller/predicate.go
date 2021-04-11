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
	ctrlevent "sigs.k8s.io/controller-runtime/pkg/event"
	ctrlpredicate "sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/zchee/kt/pkg/options"
	"github.com/zchee/kt/pkg/stdio"
)

const (
	namespaceFmt       = "%s %s%s » %s\n" // (+|-) Namespace/PodName » ContainerName
	nonNamespaceFmt    = "%s %s » %s\n"   // (+|-) PodName » ContainerName
	namespaceSeparator = "/"
)

const (
	createPodMark                 = "+"
	createPodAttr color.Attribute = color.FgHiGreen + color.Bold

	deletePodMark                 = "-"
	deletePodAttr color.Attribute = color.FgHiRed + color.Bold + color.Concealed
)

// PredicateEventFilter filters events before they are provided to handler.EventHandlers.
type PredicateEventFilter struct {
	ioStreams    stdio.Streams
	ioMu         sync.Mutex
	log          logr.Logger
	isNamespaced bool
	query        *options.Query
}

var _ ctrlpredicate.Predicate = (*PredicateEventFilter)(nil)

func (e *PredicateEventFilter) filterQuery(pod *corev1.Pod, state *corev1.ContainerStatus) bool {
	for i := range e.query.ExcludeQuery {
		if e.query.ExcludeQuery[i].MatchString(pod.Name) {
			return false // matched ExcludeQuery
		}
	}

	if e.query.ExcludeContainerQuery != nil && e.query.ExcludeContainerQuery.MatchString(pod.Name) {
		return false // matched ExcludeContainerQuery
	}

	if !e.query.ContainerState.Match(state.State) {
		return false // not matched ContainerStatus
	}

	if e.query.ContainerQuery.MatchString(pod.Name) {
		return true // matched ContainerQuery
	}

	return true
}

func (e *PredicateEventFilter) printFunc(marker string, pod *corev1.Pod, container *corev1.Container) {
	var (
		attr color.Attribute
		p, c *color.Color
	)

	switch marker {
	case createPodMark:
		attr = createPodAttr
		p, c = findColors(pod.Name)
	case deletePodMark:
		attr = deletePodAttr
	}

	mark := color.New(attr).SprintFunc()

	format := nonNamespaceFmt
	args := []interface{}{mark(marker), p.SprintfFunc()(pod.Name)}
	if container != nil {
		args = append(args, c.SprintfFunc()(container.Name))
	}

	if e.isNamespaced {
		format = namespaceFmt
		args = append([]interface{}{args[0], p.SprintfFunc()(pod.Namespace + namespaceSeparator)}, args[1:]...)
	}

	e.ioMu.Lock()
	fmt.Fprintf(e.ioStreams.Out, format, args...)
	e.ioMu.Unlock()
}

// Create implements predicate.Predicate.
func (e *PredicateEventFilter) Create(event ctrlevent.CreateEvent) bool {
	pod := event.Object.(*corev1.Pod)
	e.log.V(1).Info("PredicateEventFilter.Create", "pod", pod)

	if !e.query.PodQuery.MatchString(pod.Name) {
		return false // skip if not matched PodQuery
	}

	for i := range pod.Status.InitContainerStatuses {
		state := pod.Status.InitContainerStatuses[i]
		if !e.filterQuery(pod, &state) {
			return false
		}

		if state.State.Running != nil {
			e.printFunc(createPodMark, pod, &pod.Spec.InitContainers[i])
		}
	}
	for i := range pod.Status.ContainerStatuses {
		state := pod.Status.ContainerStatuses[i]
		if !e.filterQuery(pod, &state) {
			return false
		}

		if state.State.Running != nil {
			e.printFunc(createPodMark, pod, &pod.Spec.Containers[i])
		}
	}

	return true
}

// Delete implements predicate.Predicate.
func (e *PredicateEventFilter) Delete(event ctrlevent.DeleteEvent) bool {
	pod := event.Object.(*corev1.Pod)
	e.log.V(1).Info("PredicateEventFilter.Delete", "pod", pod)

	if !e.query.PodQuery.MatchString(pod.Name) {
		return false // skip if not matched PodQuery
	}

	for i := range pod.Status.InitContainerStatuses {
		state := pod.Status.InitContainerStatuses[i]
		if !e.filterQuery(pod, &state) {
			return false
		}
		if state.State.Terminated == nil {
			e.printFunc(deletePodMark, pod, nil)
		}
	}
	for i := range pod.Status.ContainerStatuses {
		state := pod.Status.ContainerStatuses[i]
		if !e.filterQuery(pod, &state) {
			return false
		}
		if state.State.Terminated == nil {
			e.printFunc(deletePodMark, pod, nil)
		}
	}

	return true
}

// Update implements predicate.Predicate.
func (e *PredicateEventFilter) Update(event ctrlevent.UpdateEvent) bool {
	podQueryFn := func(pod *corev1.Pod) bool {
		return e.query.PodQuery.MatchString(pod.Name) // not skip if matched PodQuery
	}

	podOld := event.ObjectOld.(*corev1.Pod)
	podNew := event.ObjectNew.(*corev1.Pod)
	e.log.Info("PredicateEventFilter.Update", "podOld", podOld, "podNew", podNew)

	if podQueryFn(podOld) {
		return false
	}

	if podQueryFn(podNew) {
		return false
	}

	return true
}

// Generic implements predicate.Predicate.
func (e *PredicateEventFilter) Generic(event ctrlevent.GenericEvent) bool {
	return false
}
