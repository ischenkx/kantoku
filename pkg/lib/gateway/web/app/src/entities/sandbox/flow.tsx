import React, {ReactElement, useCallback, useContext, useEffect, useState} from 'react'
import Dagre from '@dagrejs/dagre'
import {
    addEdge,
    Background,
    BackgroundVariant,
    Handle,
    MiniMap,
    Panel,
    Position,
    ReactFlow,
    ReactFlowProvider,
    SelectionMode,
    useEdgesState,
    useNodesInitialized,
    useNodesState,
    useReactFlow
} from '@xyflow/react'

// you also need to adjust the style import
import '@xyflow/react/dist/style.css'

// or if you just want basic styles
import '@xyflow/react/dist/base.css'

import './styles.css'
import {ThemeProvider} from 'styled-components'
import {Button, Collapse, Input, Space, theme} from 'antd'
import {
    ApartmentOutlined,
    CheckCircleOutlined,
    CloseCircleOutlined,
    FunctionOutlined,
    LoadingOutlined,
    ReloadOutlined
} from '@ant-design/icons'
import FloatingConnectionLine from './FloatingConnectionLine'
import FloatingEdge from './FloatingEdge'
import {tasksToReactflowGraph} from './graph'
import {APIWrapper} from '../../providers/common'
import {ColorModeContext} from '../../contexts/color-mode'
import {useParsed} from '@refinedev/core'

export type TaskResource = {
    name: string
    resourceId: string
    type: string
}

export type Task = {
    id: string
    inputs: TaskResource[]
    outputs: TaskResource[]
    type: string
    status: string
    parentId: string
    info: Record<string, any>
}

type Resource = {
    name: string
    type: string

}

type Node = {
    type: string
    status: string
    inputs: Resource[]
    outputs: Resource[]
}

const PipelineNodeInputs = ({inputs}: { inputs: Resource[] }) => {
    return (
        <div className={'pipeline__handle_set'}>
            {
                inputs.map(input => {
                    return <div className={'pipeline__handle_container'}>
                        <Handle
                            key={input.name}
                            id={input.name}
                            type='target'
                            position={Position.Left}
                            className={'pipeline__handle pipeline__input_handle'}
                        />
                        <span className={'pipeline__label pipeline__input_label'}>
                            {input.name}
                            <span className={'pipeline__node_handle_type'} style={{marginLeft: '4px'}}>
                                {input.type}
                            </span>
                        </span>
                    </div>
                })
            }
        </div>
    )
}

const PipelineNodeOutputs = ({outputs}: { outputs: Resource[] }) => {
    return (
        <div className={'pipeline__handle_set pipeline__output_handle_set'}>
            {
                outputs.map(output => {
                    return <div className={'pipeline__handle_container'}>
                        <span className={'pipeline__label pipeline__output_label'}>
                            {output.name}
                            <span className={'pipeline__node_handle_type'} style={{marginLeft: '4px'}}>
                                {output.type}
                            </span>
                        </span>
                        <Handle
                            key={output.name}
                            id={output.name}
                            type='source'
                            position={Position.Right}
                            className={'pipeline__handle pipeline__output_handle'}
                        />
                    </div>
                })
            }
        </div>
    )
}

const PipelineNodeStatus = ({node}: { node: Node }) => {
    const statuses: Record<string, { label: string, icon: ReactElement }> = {
        'ok': {
            label: 'Done',
            icon: <CheckCircleOutlined style={{color: '#52c41a'}}/>
        },
        'failed': {
            label: 'Failed',
            icon: <CloseCircleOutlined style={{color: '#E24329', fontSize: 20,}}/>
        },
        'running': {
            label: 'Running...',
            icon: <LoadingOutlined style={{color: '#6FA8DC'}}/>
        },
    }

    return (
        <div style={{display: 'flex', width: '100%'}}>
            <div style={{width: '30%', display: 'flex', alignItems: 'center'}}>
                <div style={{display: 'flex', alignItems: 'center'}}>
                    {statuses[node.status].icon}
                    <span className={'pipeline__node_status_label'}>
                        {statuses[node.status].label}
                    </span>
                </div>
            </div>
        </div>
    )
}

const PipelineNode = ({data: node}: { data: Node }) => {
    console.log('at pipeline node!')
    return (
        <div className={'pipeline__node'}>
            <div className={'pipeline__node_header'}>
                <div style={{marginLeft: 10}}>
                    <FunctionOutlined/>
                </div>
                <div className={'pipeline__node_header_text'} style={{marginLeft: 10}}>
                    {node.type}
                </div>
            </div>
            <div style={{marginLeft: 8, marginRight: 8, marginTop: 10}}>
                <div style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                }}>
                    <div
                        style={{
                            position: 'relative',
                            minWidth: '20%',
                            maxWidth: '50%',
                            width: 'fit-content',
                            height: '100%',
                            paddingRight: '20px',
                        }}
                    >
                        <PipelineNodeInputs inputs={node.inputs || []}/>
                    </div>
                    <div
                        style={{
                            position: 'relative',
                            minWidth: '20%',
                            maxWidth: '50%',
                            width: 'fit-content',
                            height: '100%',
                            marginLeft: 'auto',
                            paddingLeft: '20px',
                        }}
                    >
                        <PipelineNodeOutputs outputs={node.outputs || []}/>
                    </div>
                    {/*<div>*/}
                    {/*    <PipelineNodeSubtaskHandles/>*/}
                    {/*</div>*/}
                </div>
                <div style={{marginTop: 3}}>
                    <PipelineNodeStatus node={node}/>
                </div>
            </div>
        </div>
    )
}

type RFNode<T> = {
    id: string
    data: T,
    position: { x: number, y: number },
    type: string,
}

const panOnDrag = [1, 2]


const nodeTypes = {
    // resizer: ResizerNode,
    custom: PipelineNode
}

const edgeTypes = {
    floating: FloatingEdge,
}


const getLayoutedElements = (nodes, edges, options) => {
    const g = new Dagre.graphlib.Graph().setDefaultEdgeLabel(() => ({}))
    g.setGraph({rankdir: options.direction})

    console.log('edges:', edges, nodes)

    const nodeIndex = Object.fromEntries(nodes.map(node => [node.id, node]))

    edges.forEach((edge) => {
        if (edge.type === 'floating') return

        const source = nodeIndex[edge.source]
        const target = nodeIndex[edge.target]

        const sourceIndex = source.data.outputs.findIndex(output => output.name === edge.sourceHandle)
        const targetIndex = target.data.inputs.findIndex(input => input.name === edge.targetHandle)

        edge.sortPair = [sourceIndex, targetIndex]
    })

    const sorter = (e1, e2) => {
        if (!e1.sortPair || !e2.sortPair) return 0

        const pair1 = e1.sortPair, pair2 = e2.sortPair

        for (let i = 0; i < pair1.length; i++) {
            if (pair1[i] < pair2[i]) return -1
            if (pair1[i] > pair2[i]) return 1
        }
        return 0
    }

    const reverseSorter = (...args) => -sorter(...args)

    edges.sort(sorter)
    // edges.sort(reverseSorter)

    edges.forEach((edge) => g.setEdge(edge.source, edge.target))


    nodes.forEach((node) =>
        g.setNode(node.id, {
            ...node,
            width: node.measured?.width ?? 0,
            height: node.measured?.height ?? 0,
        }),
    )

    Dagre.layout(g)

    return {
        nodes: nodes.map((node) => {
            const position = g.node(node.id)
            // We are shifting the dagre node position (anchor=center center) to the top left
            // so it matches the React Flow node anchor point (top left).
            const x = position.x - (node.measured?.width ?? 0) / 2
            const y = position.y - (node.measured?.height ?? 0) / 2

            return {...node, position: {x, y}}
        }),
        edges,
    }
}

interface ControlPanelParams {
    contextId: string;

    updateContext(contextId: string): void;

    fixLayout(): void;

    reload(): void;

    fitView(): void;
}

const ControlPanel = ({contextId, updateContext, fixLayout, reload, fitView}: ControlPanelParams) => {
    const [contextIdInputValue, setContextIdInputValue] = useState<string>(contextId)

    return <>
        <Collapse style={{width: 360}}>
            <Collapse.Panel key={'1'} header={'Control Panel'} style={{background: 'rgba(0, 0, 0, 0.7)'}}>
                <Space direction={'vertical'}>
                    <Space>
                        Context:
                        <Input placeholder='Context Id'
                               defaultValue={contextIdInputValue}
                               onChange={(e) => setContextIdInputValue(e.target.value)}/>
                        <Button onClick={() => {
                            updateContext(contextIdInputValue)
                            console.log('updating context', contextIdInputValue)
                        }}>Apply</Button>
                    </Space>

                    <Space>
                        <Button icon={<ApartmentOutlined/>} onClick={() => fixLayout()}>Fix Layout</Button>
                        <Button onClick={() => fitView()}>Fit View</Button>
                        <Button icon={<ReloadOutlined/>} onClick={() => reload()}>Reload</Button>
                    </Space>
                </Space>
            </Collapse.Panel>
        </Collapse>
    </>
}

const Flow = () => {
    const {token} = theme.useToken()
    const {mode} = useContext(ColorModeContext)

    const {params} = useParsed()

    const {fitView, getNodes, getEdges} = useReactFlow()
    const [nodes, setNodes, onNodesChange] = useNodesState([])
    const [edges, setEdges, onEdgesChange] = useEdgesState([])
    const nodesInitialized = useNodesInitialized({
        includeHiddenNodes: false,
    })

    const [shouldFetch, setShouldFetch] = useState<boolean>(true)
    const [contextId, setContextId] = useState<string>(params.context_id as string || '')

    const [shouldFitView, setShouldFitView] = useState<boolean>(true)
    const [shouldLayout, setShouldLayout] = useState<boolean>(true)

    useEffect(() => {
        if (!shouldFetch) return

        if (!contextId) {
            setShouldFetch(false)
            return
        }

        console.log('fetching:', contextId)

        APIWrapper.getSpecifications().then((specifications) => {
            APIWrapper.getTasksByContext(contextId).then(apiTasks => {
                setShouldFetch(false)
                console.log('got', apiTasks)

                const tasks = apiTasks.map((t): Task => {
                    let status = 'running'
                    const rawStatus = t.info['status'] as string
                    const rawSubstatus = t.info['sub_status'] as string

                    if (rawStatus === 'finished') status = rawSubstatus

                    const spec = specifications.find((spec) => {
                        return spec.id === t.info['type']
                    })


                    return {
                        id: t.id,
                        inputs: t.inputs.map((id, idx) => {
                            const naming = (spec?.io.inputs.naming || []).find(
                                (nm) => nm.index === idx,
                            )

                            const typing = (spec?.io.inputs.types || []).find(
                                (tp) => tp.index === idx,
                            )

                            return {
                                name: naming?.name || `#${idx}`,
                                resourceId: id,
                                type: typing?.type.name || 'n/a',
                            }
                        }),
                        outputs: t.outputs.map((id, idx) => {
                            const naming = (spec?.io.outputs.naming || []).find(
                                (nm) => nm.index === idx,
                            )

                            const typing = (spec?.io.outputs.types || []).find(
                                (tp) => tp.index === idx,
                            )

                            return {
                                name: naming?.name || `#${idx}`,
                                resourceId: id,
                                type: typing?.type.name || 'n/a',
                            }
                        }),
                        type: t.info['type'],
                        status: status,
                        parentId: t.info['context_parent_id'],
                        info: t.info,
                    }
                })

                const reactflowGraph = tasksToReactflowGraph(tasks)

                setNodes(reactflowGraph.nodes)
                setEdges(reactflowGraph.edges)


                setShouldLayout(true)
                setShouldFitView(true)
                // onLayout('LR')
            })
        })
    })

    useEffect(() => {
        if (nodesInitialized) {
            if (shouldLayout) {
                onLayout('LR')
                setShouldLayout(false)
            } else if (shouldFitView) {
                setTimeout(fitView)
                setShouldFitView(false)
            }
        }
    })

    // const onInit = useCallback(() => {
    //     onLayout('LR')
    //     // setShouldFitView(true)
    // })


    const onConnect = useCallback(
        (connection) => setEdges((eds) => {
            return addEdge(connection, eds)
        }),
        [setEdges]
    )

    const onLayout = (direction) => {
        const nodes = getNodes()
        const edges = getEdges()
        const layouted = getLayoutedElements(nodes, edges, {direction})

        setNodes([...layouted.nodes])
        setEdges([...layouted.edges])

        console.log('layed out')

        setTimeout(fitView)
    }


    // if (shouldFetch) return

    return (
        <ThemeProvider theme={{antd: token, base: {color: 'mediumseagreen'}}}>
            <ReactFlow
                nodes={nodes}
                edges={edges}
                nodeTypes={nodeTypes}
                edgeTypes={edgeTypes}
                onNodesChange={onNodesChange}
                onEdgesChange={onEdgesChange}
                onConnect={onConnect}
                // onInit={onInit}
                panOnScroll
                selectionOnDrag
                panOnDrag={panOnDrag}
                selectionMode={SelectionMode.Partial}
                style={{borderRadius: 20}}
                proOptions={{hideAttribution: true}}
                fitView
                minZoom={0.01}
                maxZoom={4}
                connectionLineComponent={FloatingConnectionLine}
            >
                <Background
                    style={{backgroundColor: mode === 'dark' ? '#262626' : '#D9D9D9'}}
                    color={mode === 'dark' ? '#D9D9D9' : '#262626'}
                    variant={BackgroundVariant.Dots}
                />
                <MiniMap position='top-right' style={{borderRadius: 20, overflow: 'hidden'}}/>

                <Panel position={'top-left'}>
                    <ControlPanel
                        contextId={contextId}
                        updateContext={(ctx: string) => {
                            setContextId(ctx)
                            setShouldFetch(true)

                        }}
                        fixLayout={() => onLayout('LR')}
                        reload={() => setShouldFetch(true)}
                        fitView={fitView}
                    />
                </Panel>

            </ReactFlow>
        </ThemeProvider>
    )
}

const ProviderFlow = () => {
    return (
        <ReactFlowProvider>
            <Flow/>
        </ReactFlowProvider>
    )
}

export default ProviderFlow