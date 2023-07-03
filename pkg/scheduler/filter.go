package scheduler

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	extenderv1 "k8s.io/kube-scheduler/extender/v1"

	"github.com/houwenchen/kube-scheduler-extender/pkg/storage"
	"github.com/houwenchen/kube-scheduler-extender/pkg/utils/types"
)

type Filter struct {
	storage *storage.Storage
}

func (f *Filter) Handler(args extenderv1.ExtenderArgs) *extenderv1.ExtenderFilterResult {
	if args.Pod == nil {
		return &extenderv1.ExtenderFilterResult{
			Error: fmt.Errorf("the filed pod of ExtenderArgs can't be empty").Error(),
		}
	}

	fulleKey, err := types.KeyFunc(args.Pod)
	if err != nil {
		return &extenderv1.ExtenderFilterResult{
			Error: fmt.Errorf("parse Obj to storage key failed").Error(),
		}
	}

	// pod annotations 中存在 nodeName 的处理逻辑：
	nodeName, exist := f.storage.GetPodOfStorage(fulleKey)
	if exist {
		if args.NodeNames != nil {
			for _, v := range *args.NodeNames {
				if v == nodeName {
					return &extenderv1.ExtenderFilterResult{
						NodeNames: &[]string{nodeName},
					}
				}
			}
		}

		if args.Nodes != nil {
			for _, node := range args.Nodes.Items {
				if node.Name == nodeName {
					return &extenderv1.ExtenderFilterResult{
						Nodes: &corev1.NodeList{
							Items: []corev1.Node{
								node,
							}},
					}
				}
			}
		}

		return &extenderv1.ExtenderFilterResult{
			Error: fmt.Errorf("nodeName fetch failed").Error(),
		}
	}

	return &extenderv1.ExtenderFilterResult{
		Nodes:     args.Nodes,
		NodeNames: args.NodeNames,
	}
}
