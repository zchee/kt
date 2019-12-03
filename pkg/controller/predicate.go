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

	"github.com/zchee/kt/pkg/io"
	"github.com/zchee/kt/pkg/options"
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
	ioStreams    io.Streams
	ioMu         sync.Mutex
	log          logr.Logger
	isNamespaced bool
	query        *options.Query
}

var _ ctrlpredicate.Predicate = (*PredicateEventFilter)(nil)

func (e *PredicateEventFilter) filterQuery(state corev1.ContainerStatus) bool {
	if !e.query.ContainerQuery.MatchString(state.Name) {
		return false // not matched ContainerQuery
	}

	if e.query.ExcludeContainerQuery != nil && e.query.ExcludeContainerQuery.MatchString(state.Name) {
		return false // matched ExcludeContainerQuery
	}

	if !e.query.ContainerState.Match(state.State) {
		return false // not matched ContainerStatus
	}

	return true
}

func (e *PredicateEventFilter) printFunc(marker string, pod *corev1.Pod, container corev1.Container) {
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
	args := []interface{}{mark(marker), p.SprintfFunc()(pod.Name), c.SprintfFunc()(container.Name)}

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

	if !e.query.PodQuery.MatchString(pod.Name) {
		return false // skip if not matched PodQuery
	}

	for i, s := range pod.Status.InitContainerStatuses {
		if !e.filterQuery(s) {
			continue
		}

		if s.State.Running != nil {
			e.printFunc(createPodMark, pod, pod.Spec.InitContainers[i])
		}
	}
	for i, s := range pod.Status.ContainerStatuses {
		if !e.filterQuery(s) {
			continue
		}

		if s.State.Running != nil {
			e.printFunc(createPodMark, pod, pod.Spec.Containers[i])
		}
	}

	return true
}

// Delete implements predicate.Predicate.
func (e *PredicateEventFilter) Delete(event ctrlevent.DeleteEvent) bool {
	pod := event.Object.(*corev1.Pod)

	if !e.query.PodQuery.MatchString(pod.Name) {
		return false // skip if not matched PodQuery
	}

	for i, s := range pod.Status.InitContainerStatuses {
		if s.State.Terminated != nil {
			e.printFunc(deletePodMark, pod, pod.Spec.InitContainers[i])
		}
	}
	for i, s := range pod.Status.ContainerStatuses {
		if s.State.Terminated != nil {
			e.printFunc(deletePodMark, pod, pod.Spec.Containers[i])
		}
	}

	return false
}

// Update implements predicate.Predicate.
func (e *PredicateEventFilter) Update(event ctrlevent.UpdateEvent) bool {
	podOld := event.ObjectOld.(*corev1.Pod)
	if !e.query.PodQuery.MatchString(podOld.Name) {
		return false // skip if not matched PodQuery
	}

	podNew := event.ObjectNew.(*corev1.Pod)
	if !e.query.PodQuery.MatchString(podNew.Name) {
		return false // skip if not matched PodQuery
	}

	return true
}

// Generic implements predicate.Predicate.
func (e *PredicateEventFilter) Generic(event ctrlevent.GenericEvent) bool {
	return false
}
