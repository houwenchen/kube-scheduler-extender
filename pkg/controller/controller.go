package controller

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	coreinformer "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	corelister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

const (
	maxReQueue = 10
)

var (
	KeyFunc = cache.DeletionHandlingMetaNamespaceKeyFunc
)

type Controller struct {
	client kubernetes.Interface

	podLister  corelister.PodLister
	nodeLister corelister.NodeLister

	podInformerSynced  cache.InformerSynced
	nodeInformerSynced cache.InformerSynced

	queue workqueue.RateLimitingInterface
}

func NewController(client kubernetes.Interface, podInformer coreinformer.PodInformer, nodeInformer coreinformer.NodeInformer) (*Controller, error) {
	c := &Controller{
		client: client,
		queue:  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "scheduler controller"),
	}

	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addPod,
		UpdateFunc: c.updatePod,
		DeleteFunc: c.deletePod,
	})

	c.podLister = podInformer.Lister()
	c.nodeLister = nodeInformer.Lister()
	c.podInformerSynced = podInformer.Informer().HasSynced
	c.nodeInformerSynced = nodeInformer.Informer().HasSynced

	return c, nil
}

func (c *Controller) Run(ctx context.Context) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	if !cache.WaitForCacheSync(ctx.Done(), c.podInformerSynced) {
		klog.Fatalf("podInformer isn't synced")
	}

	if !cache.WaitForCacheSync(ctx.Done(), c.nodeInformerSynced) {
		klog.Fatalf("nodeInformer isn't synced")
	}

	for i := 0; i < 5; i++ {
		go wait.Until(c.worker, time.Second, ctx.Done())
	}

	<-ctx.Done()
}

func (c *Controller) worker() {
	for c.processNextItem() {
	}
}

func (c *Controller) processNextItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}

	defer c.queue.Done(item)

	err := c.syncPod(item.(string))
	if err != nil {
		c.handleErr(item, err)
	}

	return true
}

func (c *Controller) addPod(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		klog.InfoS("can't convert %v to *v1.Pod", obj)
		return
	}

	klog.InfoS("start add pod to queue: %s", pod.Name)
	c.queue.Add(pod)
}

func (c *Controller) updatePod(oldObj interface{}, newObj interface{}) {
	oldPod, ok := oldObj.(*corev1.Pod)
	if !ok {
		klog.InfoS("can't convert %v to *v1.Pod", oldObj)
		return
	}

	curPod, ok := newObj.(*corev1.Pod)
	if !ok {
		klog.InfoS("can't convert %v to *v1.Pod", newObj)
		return
	}

	klog.InfoS("start update pod to queue: %s", oldPod.Name)
	c.queue.Add(curPod)
}

func (c *Controller) deletePod(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		klog.InfoS("can't convert %v to *v1.Pod", obj)
		return
	}

	klog.InfoS("start delete pod to queue: %s", pod.Name)
	c.queue.Add(pod)
}

func (c *Controller) handleErr(item interface{}, err error) {
	if err == nil {
		return
	}

	if c.queue.NumRequeues(item) <= maxReQueue {
		c.queue.AddRateLimited(item)
		return
	}

	c.queue.Forget(item)
}

func (c *Controller) syncPod(key string) error {
	return nil
}
