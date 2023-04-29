package app

import "kantoku"

//kantoku:app
type App struct {
	//kantoku:inject
	Kantoku *kantoku.Kantoku
}

func (app *App) SumRange(from, to int) int {
	sum := 0
	for i := from; i <= to; i++ {
		sum += i
	}
	return sum
}

func (app *App) Fibonacci(n int) int {
	if n <= 2 {
		return 1
	}

	return app.Fibonacci(n-1) + app.Fibonacci(n-2)
}
