package stand

import (
	"github.com/ischenkx/kantoku/pkg/core/services/executor"
	"github.com/ischenkx/kantoku/pkg/lib/exe"
	"github.com/ischenkx/kantoku/pkg/lib/tasks/functional"
)

func MathExecutor() executor.Executor {
	router := exe.NewRouter()
	addExe := functional.NewExecutor[AddTask, MathInput, MathOutput](AddTask{})
	router.AddExecutor(addExe, addExe.Type())
	mulExe := functional.NewExecutor[MulTask, MathInput, MathOutput](MulTask{})
	router.AddExecutor(mulExe, mulExe.Type())
	divExe := functional.NewExecutor[DivTask, MathInput, MathOutput](DivTask{})
	router.AddExecutor(divExe, divExe.Type())
	sumExe := functional.NewExecutor[SumTask, SumInput, MathOutput](SumTask{})
	router.AddExecutor(sumExe, sumExe.Type())

	return &router
}
