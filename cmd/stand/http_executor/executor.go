package http_executor

import (
	http_tasks "github.com/ischenkx/kantoku/local/deprecated/sample_project/net/http"
	"github.com/ischenkx/kantoku/pkg/core/services/executor"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/exe"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/functional"
)

func New() executor.Executor {
	router := exe.NewRouter()
	doExe := functional.NewExecutor[*http_tasks.Do, http_tasks.DoInput, http_tasks.DoOutput](&http_tasks.Do{})
	router.AddExecutor(doExe, "http.Do")

	return &router
}
