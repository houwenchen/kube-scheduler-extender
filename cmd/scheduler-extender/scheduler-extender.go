package main

import (
	"flag"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"

	"github.com/houwenchen/kube-scheduler-extender/pkg/controller"
	"github.com/houwenchen/kube-scheduler-extender/pkg/utils/signals"
	"github.com/houwenchen/kube-scheduler-extender/pkg/utils/util"
)

var (
	syncPeriod = 10 * time.Second
)

var (
	kubeconfig = flag.String("kubeconfig", "", "path of kubernetes config file")
	port       = flag.Int("port", 9999, "port of scheduler extender serves at")
)

func main() {
	flag.Parse()

	ctx := signals.SetupSignalContext()

	clientConfig, err := util.BuildKubeConfig(*kubeconfig)
	if err != nil {
		klog.Fatal("build kubeconfig failed")
	}

	clientSet, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		klog.Fatal("build clientSet failed")
	}

	sharedInformerFactory := informers.NewSharedInformerFactory(clientSet, syncPeriod)
	podInformer := sharedInformerFactory.Core().V1().Pods()
	nodeInformer := sharedInformerFactory.Core().V1().Nodes()

	sharedInformerFactory.Start(ctx.Done())
	sharedInformerFactory.WaitForCacheSync(ctx.Done())

	controller, err := controller.NewController(clientSet, podInformer, nodeInformer)
	if err != nil {
		klog.Fatal("build controller failed")
	}

	controller.Run(ctx)
}
