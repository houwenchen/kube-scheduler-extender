package main

import (
	"flag"
	"strconv"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"

	"github.com/houwenchen/kube-scheduler-extender/pkg/controller"
	"github.com/houwenchen/kube-scheduler-extender/pkg/scheduler"
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

	// 实现平滑退出
	ctx := signals.SetupSignalContext()

	// 从 flag 或 .kube/config 或 集群内部的config 构造 kubeconfig
	clientConfig, err := util.BuildKubeConfig(*kubeconfig)
	if err != nil {
		klog.Fatal("build kubeconfig failed")
	}

	clientSet, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		klog.Fatal("build clientSet failed")
	}

	// 构造 informerFactory，为 controller 使用做准备
	sharedInformerFactory := informers.NewSharedInformerFactory(clientSet, syncPeriod)
	podInformer := sharedInformerFactory.Core().V1().Pods()
	nodeInformer := sharedInformerFactory.Core().V1().Nodes()

	controller, err := controller.NewController(clientSet, podInformer, nodeInformer)
	if err != nil {
		klog.Fatal("build controller failed")
	}

	// 运行 controller，监控 pod，pod 的 annotations 字段中的 nodeName 存在时，维护到 storage 中
	go controller.Run(ctx)

	// 注意： 这两行代码一定要在 lister 和 informer 创建之后执行
	// Informer 机制中可以看到，informerFactory start 的时候才会拉起 lister 和 informer，start 之后创建的 lister 和 informer 将不会启用
	sharedInformerFactory.Start(ctx.Done())
	sharedInformerFactory.WaitForCacheSync(ctx.Done())

	// 创建 scheduler 并运行，实现 filer 功能
	scheduler := scheduler.NewScheduler(controller.Storage)
	if err := scheduler.Run(":" + strconv.Itoa(*port)); err != nil {
		klog.Fatal("scheduler run failed")
	}
}
