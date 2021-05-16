package handler

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/hotstar/kubetimer/pkg/event"
	"github.com/hotstar/kubetimer/pkg/k8s"
	"github.com/hotstar/kubetimer/pkg/metrics"
)

type PodConditionType string

type Handler interface {
	Handle(e event.Event) error
}

type Timer struct {
	clientset kubernetes.Interface
}

func NewTimer() *Timer {
	var clientset kubernetes.Interface
	if _, err := rest.InClusterConfig(); err != nil {
		clientset = k8s.GetClientOutOfCluster()
	} else {
		clientset = k8s.GetClient()
	}

	return &Timer{clientset}
}

func (t *Timer) Handle(e event.Event) error {
	zap.S().Debugf("Event: %v", e)
	return t.process(e)
}

func (t *Timer) process(e event.Event) error {
	pod, err := t.clientset.CoreV1().Pods(e.Namespace).Get(e.Name, meta_v1.GetOptions{})
	if err != nil {
		zap.S().Errorf("Get pod: %v", err)
	}
	deployment := t.getOwnerControllerName(pod, e.Namespace)
	podConditions := make(map[core_v1.PodConditionType]core_v1.PodCondition)
	for _, condition := range pod.Status.Conditions {
		podConditions[condition.Type] = condition
	}

	if podReady, ok := podConditions[core_v1.PodReady]; !ok || (ok && podReady.Status == core_v1.ConditionFalse) {
		return fmt.Errorf("pod is not yet in %v state", core_v1.PodReady)
	}

	if podScheduled, ok := podConditions[core_v1.PodScheduled]; ok && podScheduled.Status == core_v1.ConditionTrue {
		scheduleTime := podScheduled.LastTransitionTime.Sub(pod.ObjectMeta.CreationTimestamp.Time).Seconds()
		metrics.PodScheduledTime.WithLabelValues(e.Namespace, deployment).Observe(scheduleTime)
	}
	if podInitialized, ok := podConditions[core_v1.PodInitialized]; ok && podInitialized.Status == core_v1.ConditionTrue {
		scheduleTime := podInitialized.LastTransitionTime.Sub(pod.ObjectMeta.CreationTimestamp.Time).Seconds()
		metrics.PodInitTime.WithLabelValues(e.Namespace, deployment).Observe(scheduleTime)
	}
	if containersReady, ok := podConditions[core_v1.ContainersReady]; ok && containersReady.Status == core_v1.ConditionTrue {
		scheduleTime := containersReady.LastTransitionTime.Sub(pod.ObjectMeta.CreationTimestamp.Time).Seconds()
		metrics.PodContainersReadyTime.WithLabelValues(e.Namespace, deployment).Observe(scheduleTime)
	}
	if podReady, ok := podConditions[core_v1.PodReady]; ok && podReady.Status == core_v1.ConditionTrue {
		scheduleTime := podReady.LastTransitionTime.Sub(pod.ObjectMeta.CreationTimestamp.Time).Seconds()
		metrics.PodReadyTime.WithLabelValues(e.Namespace, deployment).Observe(scheduleTime)
	}
	return nil
}

func (t *Timer) getOwnerControllerName(obj interface{}, namespace string) string {
	md := k8s.GetObjectMetaData(obj)
	if len(md.OwnerReferences) < 1 {
		return md.Name
	}
	owner := md.OwnerReferences[0]
	obj, err := t.getK8sObject(owner.Kind, owner.Name, namespace)
	if err != nil {
		zap.S().Infof("Failed to process type of %s of kind %s, falling back to current owner: %s. Error: %v", owner.Name, owner.Kind, owner.Name, err)
		return owner.Name
	}

	return t.getOwnerControllerName(obj, namespace)
}

func (t *Timer) getK8sObject(kind, name, namespace string) (interface{}, error){
	switch kind {
	case "Deployment":
		return t.clientset.AppsV1().Deployments(namespace).Get(name, meta_v1.GetOptions{})
	case "ReplicaSet":
		return t.clientset.AppsV1().ReplicaSets(namespace).Get(name, meta_v1.GetOptions{})
	case "Job":
		return t.clientset.BatchV1().Jobs(namespace).Get(name, meta_v1.GetOptions{})
	case "CronJob":
		return t.clientset.BatchV1beta1().CronJobs(namespace).Get(name, meta_v1.GetOptions{})
	default:
		return nil, errors.New("unsupported object kind")
	}
}
