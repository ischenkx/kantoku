package rider

type Settings struct {
	Scalable bool
}

type Job struct {
	Type     string
	Param    any
	Settings Settings
}
