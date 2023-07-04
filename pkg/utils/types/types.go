package types

import (
	"k8s.io/client-go/tools/cache"
)

const (
	FilterName    = "PodAnnotationsFilter"
	SchedulerName = "PodAnnotationsScheduler"
	Version       = "v0.0.1"
	VersionPath   = "/version"
	ApiPrefix     = "/pod-annotations-scheduler"
	FilterPrefix  = ApiPrefix + "/filter"
)

var (
	KeyFunc = cache.DeletionHandlingMetaNamespaceKeyFunc
)
