package controller

import (
	"time"

	"golang.org/x/time/rate"
	coreinformer "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	corelister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
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
	rateLimiter := workqueue.NewMaxOfRateLimiter(
		workqueue.NewItemExponentialFailureRateLimiter(5*time.Millisecond, 1000*time.Second),
		&workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(100), 500)},
	)
	c := &Controller{
		client: client,
		queue:  workqueue.NewNamedRateLimitingQueue(rateLimiter, "controller queue"),
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

func (c *Controller) Run() {

}

func (c *Controller) addPod(obj interface{}) {

}

func (c *Controller) updatePod(oldObj interface{}, newObj interface{}) {

}

func (c *Controller) deletePod(obj interface{}) {

}
