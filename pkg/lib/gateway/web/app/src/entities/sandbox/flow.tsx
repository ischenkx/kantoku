import React, {useCallback} from 'react';
import ReactFlow, {
    addEdge, Background,
    BackgroundVariant,
    Handle, MiniMap,
    Position,
    SelectionMode,
    useEdgesState,
    useNodesState
} from 'reactflow';
import 'reactflow/dist/style.css';
import './styles.css'
import {ThemeProvider} from "styled-components";
import { theme } from "antd";


const PipelineNodeInputs = ({inputs}) => {
    return (
        <div style={{
            display: 'flex',
            alignItems: 'center',
            flexDirection: 'column',
            justifyContent: 'space-evenly',
            width: '100%',
            height: '100%',
            // position: 'relative'
        }}>
            {
                inputs.map(input => {
                    console.log(input);
                    return <div style={{
                        // position: 'relative',
                        display: 'flex',
                        width: '100%',
                        flexDirection: 'row',
                        alignItems: 'center',
                        // justifyContent: 'center',
                        fontSize: '13px'
                    }}>
                        <Handle
                            key={input}
                            id={input}
                            type="target"
                            position={Position.Left}
                            // id="input"
                            style={{
                                background: '#555',
                                position: 'relative',
                                top: 'auto',
                                transform: 'none',
                            }}
                        />
                        <span>{input}</span>
                    </div>;
                    {/*    Input*/}
                    {/*</div>*/}
                })
            }
        </div>
    )
}

const PipelineNodeOutputs = ({outputs}) => {
    return (
        <div style={{
            display: 'flex',
            alignItems: 'self-end',
            flexDirection: 'column',
            justifyContent: 'space-evenly',
            width: '100%',
            height: '100%',
            position: 'relative'
        }}>
            {
                outputs.map(output => {
                    return <div style={{
                        marginLeft: 'auto',
                        display: 'flex',
                        width: '100%',
                        flexDirection: 'row',
                        alignItems: 'center',
                        justifyContent: 'flex-end',
                        fontSize: '13px',
                    }}>
                        <span>{output}</span>
                        <Handle
                            key={output}
                            id={output}
                            type="source"
                            position={Position.Right}
                            // id="input"
                            className={'pipeline-node__handle'}
                            style={{
                                background: '#555',
                                position: 'relative',
                                top: 'auto',
                                transform: 'none',
                                minWidth: '10px'
                            }}
                        />
                    </div>;
                    {/*    Input*/
                    }
                    {/*</div>*/}
                })
            }
        </div>
    )
}

const PipelineNode = ({ data }) => {
    console.log(data)
    return (
        <div style={{
            position: 'relative',
            minWidth: '150px',
            height: '100px',
            // border: '1px solid black',
            padding: '10px',
            borderRadius: '10px',
            display: 'flex',
            justifyContent: 'space-between',
            color: '#777',
            backgroundColor: '#AAA'

        }}>
            <div
                style={{
                    position: 'relative',
                    minWidth: '20%',
                    maxWidth: '28%',
                    width: 'fit-content',
                    height: '100%',
                }}
            >
                <PipelineNodeInputs inputs={data.inputs || []} />
            </div>
            <div
                style={{
                    position: 'relative',
                    minWidth: '40%',
                    maxWidth: '60%',
                    width: 'fit-content',
                    height: '100%',
                    paddingLeft: '15px',
                    paddingRight: '15px',
                    boxSizing: 'border-box',
                }}
            >
                {data.label}
            </div>
            <div
                style={{
                    position: 'relative',
                    minWidth: '20%',
                    maxWidth: '30%',
                    width: 'fit-content',
                    height: '100%',
                }}
            >
                <PipelineNodeOutputs outputs={data.outputs || []} />
            </div>
        </div>
    );
};


const initialNodes = [
    {
        id: '1',
        data: {
            label: 'Sum',
            inputs: [
                'A',
                'B'
            ],
            outputs: [
                'R',
            ]
        },
        position: {x: 150, y: 0},
        type: 'custom'
    },
    {
        id: '2',
        data: {
            label: 'Divide',
            inputs: [
                'A',
                'B'
            ],
            outputs: [
                'R',
            ]
        },
        position: {x: 75, y: 0},
        type: 'custom'
    },
    {
        id: '3',
        data: {
            label: 'Multiply',
            inputs: [
                'A',
                'B'
            ],
            outputs: [
                'Result',
            ]
        },
        position: {x: 120, y: 0},
        type: 'custom'
    },
];

const initialEdges = [
    // { id: 'e1-2', source: '1', target: '2' },
    // { id: 'e1-3', source: '1', target: '3' },
];

const panOnDrag = [1, 2];


const nodeTypes = {
    // resizer: ResizerNode,
    custom: PipelineNode
};


function Flow() {
    const {token} = theme.useToken()
    const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
    const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

    const onConnect = useCallback(
        (connection) => setEdges((eds) => addEdge(connection, eds)),
        [setEdges]
    );

    return (
        <ThemeProvider theme={{ antd: token, base: { color: "mediumseagreen" } }}>
            <ReactFlow
                nodes={nodes}
                edges={edges}
                nodeTypes={nodeTypes}
                onNodesChange={onNodesChange}
                onEdgesChange={onEdgesChange}
                onConnect={onConnect}
                panOnScroll
                selectionOnDrag
                panOnDrag={panOnDrag}
                selectionMode={SelectionMode.Partial}
            >
                <Background variant={BackgroundVariant.Dots} />
                <MiniMap/>
            </ReactFlow>
        </ThemeProvider>

    );
}

export default Flow;