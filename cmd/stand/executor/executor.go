package executor

import (
	http_tasks "github.com/ischenkx/kantoku/cmd/stand/sample_project/net/http"
	"github.com/ischenkx/kantoku/cmd/stand/sample_project/recursive"
	"github.com/ischenkx/kantoku/cmd/stand/sample_project/scraper"
	"github.com/ischenkx/kantoku/cmd/stand/sample_project/test"
	"github.com/ischenkx/kantoku/pkg/core/services/executor"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/exe"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/fn"
	"log/slog"
	"reflect"
	"strings"
)

func New() executor.Executor {
	router := exe.NewRouter()

	//doExe := fn.NewExecutor[http_tasks.Do, http_tasks.DoInput, http_tasks.DoOutput](http_tasks.Do{})
	//router.AddExecutor(doExe, "http.Do")
	//
	//randExe := fn.NewExecutor[*test.RandFail, test.RandFailInput, test.RandFailOutput](&test.RandFail{})
	//router.AddExecutor(randExe, "test.RandFail")

	registerExecutor[http_tasks.Do](router)
	registerExecutor[*test.RandFail](router)
	registerExecutor[recursive.A](router)
	registerExecutor[recursive.B](router)
	registerExecutor[scraper.Scrape](router)
	registerExecutor[scraper.DownloadPage](router)
	registerExecutor[scraper.ParsePage](router)
	registerExecutor[scraper.ExtractImages](router)

	return router
}

func registerExecutor[F fn.AbstractFunction[I, O], I, O any](router *exe.Router) {
	var task F
	typ := reflect.TypeFor[F]()
	if typ.Kind() == reflect.Ptr {
		reflect.ValueOf(&task).Elem().Set(reflect.New(typ.Elem()))

		v := reflect.ValueOf(task)
		e := v.Elem()
		nv := reflect.New(typ.Elem()).Elem()

		e.Set(nv)
	} else {
		reflect.ValueOf(&task).Elem().Set(reflect.New(typ).Elem())
	}

	exec := fn.NewExecutor(task)

	strtype := exec.Type()
	strtype = strings.TrimPrefix(strtype, "github.com/ischenkx/kantoku/cmd/stand/")

	slog.Info("registering",
		"type", strtype)

	router.AddExecutor(exec, strtype)
}
