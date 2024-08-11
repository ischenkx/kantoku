package fn_d

import (
	"errors"
	"github.com/ischenkx/kantoku/pkg/core/resource"
	"github.com/samber/lo"
	"log"
)

// turns out we don't need that 🤷‍♂️

type GraphTask struct {
	Inputs  []resource.ID
	Outputs []resource.ID
}

type node struct {
	finished     bool
	processing   bool
	dependencies []*node
	task         GraphTask
}

func TopSortTasks(tasks []GraphTask) ([]GraphTask, error) {
	graph := newTaskGraph(tasks)
	sorted := make([]GraphTask, 0)
	for _, n := range graph {
		if !n.finished {
			err := dfs(n, func(n *node) {
				sorted = append(sorted, n.task)
			})
			if err != nil {
				return nil, err
			}
		}
	}
	return sorted, nil
}

func newTaskGraph(tasks []GraphTask) []*node {
	input2nodes := map[resource.ID][]*node{}
	nodes := lo.Map(tasks, func(t GraphTask, _ int) *node {
		n := node{
			finished:     false,
			processing:   false,
			dependencies: make([]*node, 0),
			task:         t,
		}
		for _, input := range t.Inputs {
			_, has := input2nodes[input]
			if !has {
				input2nodes[input] = make([]*node, 0)
			}
			input2nodes[input] = append(input2nodes[input], &n)
		}

		return &n
	})

	for _, n := range nodes {
		for _, out := range n.task.Outputs {
			for _, dep := range input2nodes[out] {
				n.dependencies = append(n.dependencies, dep)
			}
		}
	}
	return nodes
}

func dfs(cur *node, f func(n *node)) error {
	if cur.finished {
		log.Panic("entered finished node")
	}
	if cur.processing {
		return errors.New("cyclic dependency in tasks")
	}
	cur.processing = true

	for _, dep := range cur.dependencies {
		if dep.finished {
			continue
		}
		err := dfs(cur, f)
		if err != nil {
			return err
		}
	}

	f(cur)
	cur.finished = true
	cur.processing = false
	return nil
}

//func (g taskGraph) allocateResourcesForGraph(ctx context.Context, sys system.AbstractSystem) (map[*tasks.Future[any]]resource.Resource, error) {
//	future2resource := make(map[*tasks.Future[any]]resource.Resource)
//	for _, node := range g.nodes {
//		for _, input := range node.task.Inputs {
//			if _, has := future2resource[input]; input.HasResource() || has {
//				continue
//			}
//			future2resource[input] = resource.Resource{}
//		}
//		for _, output := range node.task.Outputs {
//			if _, has := future2resource[output]; output.HasResource() || has {
//				continue
//			}
//			future2resource[output] = resource.Resource{}
//		}
//	}
//	resources, err := sys.Resources().Alloc(ctx, len(future2resource))
//	if err != nil {
//		return nil, err
//	}
//	idx := 0
//	for fut := range future2resource {
//		future2resource[fut] = resource.Resource{ID: resources[idx]}
//		idx += 1
//	}
//	return future2resource, nil
//}
