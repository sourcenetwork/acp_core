import { CheckExplainGraph, CheckExplainNode, SearchResult } from '@/types/proto-js/sourcenetwork/acp_core/check_explain';
import { cn } from '@/utils/classnames';
import dagre from '@dagrejs/dagre';
import {
    BezierEdge,
    ConnectionLineType,
    Edge,
    EdgeProps,
    Handle,
    Node,
    NodeProps,
    Panel,
    Position,
    ReactFlow,
    ReactFlowProvider,
    useEdgesState,
    useNodesState,
    useOnSelectionChange,
    useReactFlow
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { cva } from 'class-variance-authority';
import { Ampersand, ArrowRight, Minus, Plus } from 'lucide-react';
import React, { memo, MouseEvent, useCallback, useEffect, useMemo, useState } from 'react';

// Flow Node/Edge Types
type NodeBaseType = 'root' | 'operatorNode' | 'pathNode';
type EdgeTypes = 'edge';

enum NodeTypeString {
    UNION_NODE = 'UNION_NODE',
    INTERSECTION_NODE = 'INTERSECTION_NODE',
    DIFF_NODE = 'DIFF_NODE',
    USERSET_NODE = 'USERSET_NODE',
    UNRECOGNIZED = 'UNRECOGNIZED',
}

type FlowNodeData = {
    id: string;
    label: string;
    reason: string;
    baseType: NodeBaseType;
    nodeType: NodeTypeString;
    result?: SearchResult;
}

type FlowEdgeData = {
    message: string;
    reason: string;
    nodeResult?: SearchResult;
    sourceNode: CheckExplainNode;
    targetNode: CheckExplainNode;
}

type PathFlowNode = Node<Partial<FlowNodeData>, 'pathNode'>;
type OperatorFlowNode = Node<Partial<FlowNodeData>, 'operatorNode'>;

type CustomFlowNode = PathFlowNode | OperatorFlowNode;
type CustomFlowEdge = Edge<Partial<FlowEdgeData>, EdgeTypes>;

const operatorIcons: Record<NodeTypeString, React.ElementType> = {
    [NodeTypeString.UNION_NODE]: Plus,
    [NodeTypeString.INTERSECTION_NODE]: Ampersand,
    [NodeTypeString.DIFF_NODE]: Minus,
    [NodeTypeString.USERSET_NODE]: ArrowRight,
    [NodeTypeString.UNRECOGNIZED]: Plus,
}

const dagreNodeDimensions: Record<NodeBaseType, { height: number, width: number }> = {
    'pathNode': {
        'height': 60,
        'width': 220,
    },
    'operatorNode': {
        'height': 30,
        'width': 30,
    },
    'root': {
        'height': 60,
        'width': 220,
    },
}

const getNodeDimensions = (baseType: NodeBaseType) => {
    return {
        width: dagreNodeDimensions[baseType]?.width || 220,
        height: dagreNodeDimensions[baseType]?.height || 60
    }
}

const castNodeType = (node: CheckExplainNode) => node.nodeType as unknown as NodeTypeString;
const getNodeBaseType = (node: CheckExplainNode) => castNodeType(node) === NodeTypeString.USERSET_NODE ? 'pathNode' : 'operatorNode';

// Layout function using Dagre
const getLayoutedElements = (nodes: CustomFlowNode[] = [], edges: CustomFlowEdge[] = []) => {
    const dagreGraph = new dagre.graphlib.Graph().setDefaultEdgeLabel(() => ({}));

    dagreGraph.setGraph({
        rankdir: 'TB',
        nodesep: 40,
        ranksep: 50,
    });

    nodes.forEach((node) => dagreGraph.setNode(node.id, getNodeDimensions(node.type)));
    edges.forEach((edge) => dagreGraph.setEdge(edge.source, edge.target));

    dagre.layout(dagreGraph);

    const newNodes = nodes.map((node) => {
        const nodeWithPosition = dagreGraph.node(node.id);
        const newNode: Node = {
            ...node,
            targetPosition: Position.Top,
            sourcePosition: Position.Bottom,
            position: {
                x: nodeWithPosition.x - getNodeDimensions(node.type).width / 2,
                y: nodeWithPosition.y - getNodeDimensions(node.type).height / 2,
            },
        };

        return newNode;
    });

    return { nodes: newNodes, edges };
};

const convertGraphToReactFlow = (checkGraph: CheckExplainGraph) => {
    const nodeMap = new Map(checkGraph.nodes.map(node => [node.id, node]));

    const reactFlowNodes: CustomFlowNode[] = checkGraph.nodes.map((node) => ({
        id: node.id,
        type: getNodeBaseType(node),
        data: {
            id: node.id,
            label: node.text || node.id,
            reason: node.detail,
            baseType: getNodeBaseType(node),
            nodeType: castNodeType(node),
            result: node.result,
        },
        zIndex: 2,
        position: { x: 0, y: 0 }, // Gets set by the layout function
        style: {
            width: getNodeDimensions(getNodeBaseType(node)).width,
            height: getNodeDimensions(getNodeBaseType(node)).height,
        },
    }));

    const reactFlowEdges: CustomFlowEdge[] = checkGraph.edges.map((edge) => {
        const targetNode = nodeMap.get(edge.destNodeId);
        const sourceNode = nodeMap.get(edge.sourceNodeId);
        const id = `e${edge.sourceNodeId}-${edge.destNodeId}`;

        const { authorized } = targetNode?.result ?? { authorized: false, explored: false };

        return {
            id,
            type: 'edge',
            source: edge.sourceNodeId,
            target: edge.destNodeId,
            zIndex: authorized ? 1 : 0,
            animated: authorized ? true : false,
            data: {
                id,
                sourceNode: sourceNode,
                targetNode: targetNode,
                message: edge.message,
                nodeResult: targetNode?.result,
            },
        };
    });

    return { nodes: reactFlowNodes, edges: reactFlowEdges };
};



const PathNodeComponent = memo((props: NodeProps<PathFlowNode>) => {
    const { label, result, reason } = props.data;
    const { selected } = props;
    const { authorized, explored } = result ?? { authorized: false, explored: false };

    const variants = cva('px-4 text-xs bg-editor shadow-lg/30 font-light tracking-wide text-foreground overflow-auto border-l-2 h-[60px] flex flex-col justify-center', {
        variants: {
            authorized: {
                true: 'border-src-secondary',
            },
            explored: {
                false: 'opacity-40 hover:opacity-100',
            },
            selected: {
                true: 'opacity-100 ring-1 ring-primary/50',
            },
        },
        defaultVariants: {
            authorized: false,
            explored: false,
            selected: false,
        },
        compoundVariants: [
            {
                authorized: true,
                explored: true,
                className: 'opacity-100',
            },
            {
                authorized: false,
                explored: true,
                className: 'border-src-error opacity-80 hover:opacity-100',
            },
        ],
    })

    return (
        <div className={variants({ authorized, explored, selected })}>
            <Handle type="target" position={Position.Top} />
            <div className="flex flex-col gap-1 font-mono">
                <div className="text-xs font-medium tracking-wider whitespace-nowrap overflow-hidden text-ellipsis">{label?.replace('relation node:', '').replace('(', '').replace(')', '')}</div>
                {reason && <div className="text-xs text-muted-foreground whitespace-nowrap overflow-hidden text-ellipsis">{reason}</div>}
            </div>
            <Handle type="source" position={Position.Bottom} />
        </div>
    );
});

const OperatorNodeComponent = memo((props: NodeProps<OperatorFlowNode>) => {
    const { result } = props.data;
    const { selected } = props;
    const { authorized, explored } = result ?? { authorized: false, explored: false };

    const operatorType = props.data.nodeType ?? NodeTypeString.UNION_NODE;
    const OperatorIcon = operatorIcons[operatorType] ?? Plus;

    const variants = cva('p-1 text-xs text-primary overflow-hidden border size-[30px] flex items-center justify-center text-center rounded-full', {
        variants: {
            authorized: {
                true: 'border-src-secondary',
            },
            explored: {
                true: 'opacity-80 hover:opacity-100',
                false: 'opacity-40 hover:opacity-100',
            },
            selected: {
                true: 'opacity-100 ring-1 ring-primary/50',
            },
        },
        defaultVariants: {
            authorized: false,
            explored: false,
            selected: false,
        },
    })

    return (
        <div className={variants({ authorized, explored, selected })}>
            <Handle type="target" position={Position.Top} />
            <OperatorIcon size={14} />
            <Handle type="source" position={Position.Bottom} />
        </div>
    );
});

const CustomEdgeComponent = memo((props: EdgeProps<CustomFlowEdge>) => {

    const { authorized, explored } = props.data?.nodeResult ?? { authorized: false, explored: false };

    const edgeStyles = {
        'unexplored': 'var(--color-gray-600)',
        'authorized': 'var(--src-secondary)',
        'unauthorized': 'var(--src-error)',
    }

    const stroke = edgeStyles[!explored ? 'unexplored' : authorized ? 'authorized' : 'unauthorized'];

    return <g>
        <BezierEdge {...props}
            style={{
                stroke,
                strokeWidth: 2,
                opacity: explored ? 1 : 0.3,
            }}
        />
    </g>
});

const NodeTypesMap = {
    pathNode: PathNodeComponent,
    operatorNode: OperatorNodeComponent,
};

const EdgeTypesMap = {
    edge: CustomEdgeComponent,
};

interface ExplainCheckGraphProps {
    explainGraph: CheckExplainGraph;
    className?: string;
    maximized?: boolean;
}

const ExplainCheckGraph: React.FC<ExplainCheckGraphProps> = ({ explainGraph, className, maximized = false }) => {
    const graph = useReactFlow<CustomFlowNode, CustomFlowEdge>();

    const { nodes: layoutedNodes, edges: layoutedEdges } = useMemo(() => {
        const { nodes, edges } = convertGraphToReactFlow(explainGraph);
        return getLayoutedElements(nodes, edges);
    }, [explainGraph]);

    const [nodes, setNodes, onNodesChange] = useNodesState(layoutedNodes);
    const [edges, setEdges, onEdgesChange] = useEdgesState(layoutedEdges);
    const [selectedNode, setSelectedNode] = useState<CustomFlowNode | null>(null);

    // Reset the graph view in the viewport
    const resetGraphView = useCallback(() => {
        setTimeout(() => {
            graph?.fitView({ padding: 0.1, duration: 700, maxZoom: 1, minZoom: 0.2 });
        }, 100);
    }, [graph]);

    useEffect(() => {
        setNodes(layoutedNodes);
        setEdges(layoutedEdges);

        resetGraphView();
    }, [layoutedNodes, layoutedEdges]);

    useEffect(() => {
        resetGraphView();
    }, [maximized]);

    const onChange = useCallback(({ nodes }: { nodes: CustomFlowNode[], edges: CustomFlowEdge[] }) => {
        setSelectedNode(nodes[0]);
    }, []);

    useOnSelectionChange<CustomFlowNode, CustomFlowEdge>({ onChange });

    const onNodeMouseEnter = useCallback((_e: MouseEvent, node: Node) => {
        setSelectedNode(node as CustomFlowNode);
    }, []);

    const onNodeMouseLeave = useCallback((_e: MouseEvent, _n: Node) => {
        setSelectedNode(null);
    }, []);

    return (
        <div className="h-full">
            <ReactFlow
                id="explain-check-graph"
                nodes={nodes}
                edges={edges}
                onNodesChange={onNodesChange}
                onEdgesChange={onEdgesChange}
                connectionLineType={ConnectionLineType.SmoothStep}
                nodeTypes={NodeTypesMap}
                edgeTypes={EdgeTypesMap}
                maxZoom={2}
                minZoom={0.2}
                fitView
                className={cn('!bg-graph-bg', className)}
                proOptions={{ hideAttribution: true, }}
                nodesDraggable={false}
                nodesConnectable={false}
                elevateEdgesOnSelect={true}
                onNodeMouseEnter={onNodeMouseEnter}
                onNodeMouseLeave={onNodeMouseLeave}
            >
                {selectedNode && <ExplainCheckSelectionPanel node={selectedNode} />}
            </ReactFlow>
        </div>
    );
};

const ExplainCheckResult = ({ explainGraph, className, maximized = false }: { explainGraph: CheckExplainGraph, className?: string, maximized?: boolean }) => {
    return (
        <ReactFlowProvider>
            <ExplainCheckGraph explainGraph={explainGraph} className={className} maximized={maximized} />
        </ReactFlowProvider>
    )
}

export default ExplainCheckResult;

const ExplainCheckSelectionPanel = ({ node }: { node: CustomFlowNode }) => {
    const { label, reason, result } = node.data;

    const { authorized, explored } = result ?? { authorized: false, explored: false };

    return (
        <Panel position="bottom-right" className="text-xs font-mono bg-editor shadow-lg/30 border border-divider overflow-auto w-full max-w-[600px] p-3">
            <div className='font-medium mb-2 border-b border-divider pb-1'>{label}</div>
            <div className='grid grid-cols-2 gap-2 text-muted-foreground'>

                {reason && <>
                    <div className='font-bold'>Reason:</div>
                    <div className=''>{reason.replace('reason: ', '')}</div>
                </>}

                <div className='font-bold'>Authorized:</div>
                <div className=''>{authorized ? 'Yes' : 'No'}</div>
                <div className='font-bold'>Explored:</div>
                <div className=''>{explored ? 'Yes' : 'No'}</div>
            </div>
        </Panel>
    )
}
