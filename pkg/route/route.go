package route

import (
	"net/http"

	"github.com/houwenchen/kube-scheduler-extender/pkg/scheduler"
	"github.com/julienschmidt/httprouter"
)

func checkBody(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "http request body can't be empty", 400)
		return
	}
}

func FilterRouter(filter scheduler.Filter) httprouter.Handle {
	return nil
}
