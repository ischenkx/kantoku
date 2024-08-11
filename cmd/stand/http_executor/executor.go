package http_executor

import (
	http_tasks "github.com/ischenkx/kantoku/cmd/stand/sample_project/net/http"
	"github.com/ischenkx/kantoku/cmd/stand/sample_project/test"
	"github.com/ischenkx/kantoku/pkg/core/services/executor"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/exe"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn_d"
)

func New() executor.Executor {
	router := exe.NewRouter()
	doExe := fn_d.NewExecutor[http_tasks.Do, http_tasks.DoInput, http_tasks.DoOutput](http_tasks.Do{})
	router.AddExecutor(doExe, "http.Do")

	randExe := fn_d.NewExecutor[test.RandFail, test.RandFailInput, test.RandFailOutput](test.RandFail{})
	router.AddExecutor(randExe, "test.RandFail")

	return &router
}
