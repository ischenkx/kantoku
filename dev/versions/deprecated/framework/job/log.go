package job

type Log struct {
	Spawned []string
}

func (log *Log) Merge(other *Log) {
	log.Spawned = append(log.Spawned, other.Spawned...)
}
