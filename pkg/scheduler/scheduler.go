package scheduler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"k8s.io/klog/v2"

	"github.com/houwenchen/kube-scheduler-extender/pkg/route"
	"github.com/houwenchen/kube-scheduler-extender/pkg/scheduler/extender"
	"github.com/houwenchen/kube-scheduler-extender/pkg/storage"
	"github.com/houwenchen/kube-scheduler-extender/pkg/utils/types"
)

type Scheduler struct {
	name   string
	router *httprouter.Router
	filter *extender.Filter
}

func NewScheduler(storage *storage.Storage) *Scheduler {
	s := &Scheduler{
		name:   types.SchedulerName,
		router: httprouter.New(),
		filter: extender.NewFilter(storage),
	}

	s.router.GET(types.VersionPath, route.VersionRoute)
	s.router.POST(types.FilterPrefix, route.FilterRouter(*s.filter))

	return s
}

func (s *Scheduler) Run(addr string) error {
	klog.Info("scheduler start")

	return http.ListenAndServe(addr, s.router)
}
