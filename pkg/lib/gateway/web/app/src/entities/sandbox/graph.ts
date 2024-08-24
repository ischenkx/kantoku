import {Task} from './flow'
import {MarkerType} from 'reactflow'

type TaskNode = {
    task: Task | null,
    subTasks: TaskNode[],
    visited: boolean
}

const makeGraph = (tasks: Task[]): TaskNode[] => {
    const idToTask: Record<string, Task> = {}
    for (const task of tasks) {
        idToTask[task.id] = task
    }

    const nodes: Record<string, TaskNode> = {}
    const initNode = (id: string, task: Task | null) => {
        let node = nodes[id]
        if (!node) {
            node = {
                task: task,
                subTasks: [],
                visited: false,
            }
            nodes[id] = node
        }

        if (!node.task) node.task = task

        return node
    }

    for (const task of tasks) {
        const node = initNode(task.id, task)

        if (task.parentId) {
            const parentNode = initNode(task.parentId, null)
            parentNode.subTasks.push(node)
        }
    }

    return Object.values(nodes)
}

export const tasksToReactflowGraph = (tasks: Task[]) => {
    const graph = makeGraph(tasks)

    const resourceToProvider: Record<string, { node: TaskNode, output: string }> = {}

    const _dfs = (node: TaskNode, handler: (node: TaskNode) => void) => {
        if (node.visited) return
        node.visited = true

        handler(node)

        for (const subNode of node.subTasks) {
            _dfs(subNode, handler)
        }
    }

    const dfs = (handler: (node: TaskNode) => void) => {
        for (const node of graph) {
            _dfs(node, handler)
        }

        for (const node of graph) node.visited = false
    }

    dfs((node: TaskNode) => {
        const task = node.task

        for (const output of task?.outputs || []) {
            resourceToProvider[output.resourceId] = {node: node, output: output.name}
        }
    })

    const nodes = []
    const edges = []

    dfs((node: TaskNode) => {
        const task = node.task

        const rfNode = {
            id: task?.id,
            data: {
                type: task?.type,
                status: task?.status,
                inputs: [],
                outputs: []
            },
            type: 'custom',
            position: {x: 0, y: 0}
        }

        if (task?.parentId) {
            edges.push({
                id: `${task.parentId}-${rfNode.id}__SUB_TASK`,
                source: task.parentId,
                target: rfNode.id,
                type: 'floating',
                animated: true,
                style: { stroke: 'rgb(158, 118, 255)', strokeWidth: 2 },
                markerEnd: {
                    type: MarkerType.ArrowClosed,
                    width: 10,
                    height: 10,
                    strokeWidth: 2,
                    color: 'rgb(158, 118, 255)',
                },
            })
        }

        for (const input of task?.inputs || []) {
            rfNode.data.inputs.push({
                name: input.name,
                type: input.type,
            })

            const provider = resourceToProvider[input.resourceId]
            if (!provider) continue

            const {node: sourceNode, output: sourceOutput} = provider

            edges.push({
                id: `${sourceNode.task?.id}__${sourceOutput}-${rfNode.id}__${input.name}`,
                source: sourceNode.task?.id,
                target: rfNode.id,
                sourceHandle: sourceOutput,
                targetHandle: input.name,
            })
        }

        for (const output of task?.outputs || []) {
            rfNode.data.outputs.push({
                name: output.name,
                type: output.type,
            })
        }

        nodes.push(rfNode)
    })

    return {nodes, edges}
}