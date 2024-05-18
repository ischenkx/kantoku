import { useCallback } from 'react';
import ReactFlow, {
    addEdge,
    Handle,
    NodeResizer,
    Position,
    SelectionMode,
    useEdgesState,
    useNodesState
} from 'reactflow';
import 'reactflow/dist/style.css';

const initialNodes = [
    {
        id: '1',
        data: { label: 'Node 1' },
        position: { x: 150, y: 0 },
        // type: 'resizer'
    },
    {
        id: '2',
        data: { label: 'Node 2' },
        position: { x: 0, y: 150 },
    },
    {
        id: '3',
        data: { label: 'Node 3' },
        position: { x: 300, y: 150 },
    },
];

const initialEdges = [
    { id: 'e1-2', source: '1', target: '2' },
    { id: 'e1-3', source: '1', target: '3' },
];

const panOnDrag = [1, 2];

function ResizerNode({ data }) {
    return (
        <>
            {/*<Node />*/}
            <div style={{ padding: 10 }}>{data.label}</div>
            <div
                style={{
                    display: 'flex',
                    position: 'absolute',
                    justifyContent: 'space-evenly',
                    flexDirection: 'column',
                    left: 0,
                    top: 0,
                    height: '100%',
                }}
            >
                {/*<Handle type="target" position={Position.Left} />*/}

                <Handle
                    style={{ position: 'relative', left: 0, transform: 'none' }}
                    id="a"
                    type="target"
                    position={Position.Left}
                />
                <Handle
                    style={{ position: 'relative', left: 0, transform: 'none' }}
                    id="b"
                    type="target"
                    position={Position.Left}
                />
            </div>
        </>
    );
}

const nodeTypes = {
    resizer: ResizerNode,
};


function Flow() {
    const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
    const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

    const onConnect = useCallback(
        (connection) => setEdges((eds) => addEdge(connection, eds)),
        [setEdges]
    );

    return (
        <div style={{ width: '100%', height: '100%', background: 'white' }}>
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
            />
        </div>

    );
}

export default Flow;