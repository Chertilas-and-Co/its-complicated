import React, { useState, useRef, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import ForceGraph2D from "react-force-graph-2d";
import { forceCollide } from "d3-force";

function ForceGraph() {
    const graphRef = useRef(null);
    const [hoverNode, setHoverNode] = useState(null);
    const [graphData, setGraphData] = useState({ nodes: [], links: [] });
    const navigate = useNavigate();

    useEffect(() => {
        fetch('http://localhost:8080/graph-data')
            .then(res => res.json())
            .then(dataFromBackend => {
                const nodeMap = {};
                dataFromBackend.forEach(({ id_1, id_2, name_1, desc_1, name_2, desc_2 }) => {
                    if (!nodeMap[id_1]) nodeMap[id_1] = { id: id_1, name: name_1, description: desc_1, connections: 0 };
                    if (!nodeMap[id_2]) nodeMap[id_2] = { id: id_2, name: name_2, description: desc_2, connections: 0 };
                    nodeMap[id_1].connections++;
                    nodeMap[id_2].connections++;
                });

                Object.values(nodeMap).forEach(node => {
                    let totalK = 0;
                    dataFromBackend.forEach(({ id_1, id_2 }) => {
                        if (id_1 === node.id || id_2 === node.id) totalK++;
                    });
                    node.totalK = totalK;
                });

                const uniqueNodes = Object.values(nodeMap);

                const scaleSize = (k) => {
                    const minSize = 10;
                    const maxSize = 40;
                    return Math.min(maxSize, minSize + k * 5);
                };

                const nodes = uniqueNodes.map(node => ({
                    id: node.id,
                    name: node.name,
                    description: node.description,
                    size: scaleSize(node.totalK),
                    totalK: node.totalK,
                }));

                const BASE_DISTANCE = 1000;
                const ALPHA = 10;

                const links = dataFromBackend.map(({ id_1, id_2, subscribers_1, subscribers_2, common_subscribers }) => {
                    const subscribers_sum = subscribers_1 + subscribers_2;
                    const common_ratio = subscribers_sum > 0 ? common_subscribers / subscribers_sum : 0;
                    const distance = BASE_DISTANCE / (1 + ALPHA * common_ratio);
                    return {
                        source: id_1,
                        target: id_2,
                        color: "#bbb", // Softer color for links
                        distance,
                    };
                });

                setGraphData({ nodes, links });
            });
    }, []);

    useEffect(() => {
        if (graphRef.current && graphData.nodes.length > 0) {
            const graph = graphRef.current;
            graph.d3Force("link").distance(link => link.distance).iterations(1);
            graph.d3Force("charge").strength(node => -Math.abs(node.size * 100) / (1 + node.totalK));
            graph.d3Force("collide", forceCollide(node => node.size * 1.1).strength(0.2).iterations(1));
            graph.d3Force("center", null);
            graph.d3ReheatSimulation();
        }
    }, [graphData]);

    const handleNodeClick = (node) => {
        navigate(`/community/${node.id}`);
    };

    return (
        <div style={{ cursor: hoverNode ? 'pointer' : 'default' }}>
            <ForceGraph2D
                ref={graphRef}
                graphData={graphData}
                nodeVal={node => node.size || 10}
                onNodeClick={handleNodeClick}
                onNodeHover={setHoverNode}
                linkWidth={1}
                linkColor="color"
                linkOpacity={0.5}
                nodeCanvasObject={(node, ctx, globalScale) => {
                    const label = node.name;
                    const fontSize = 12 / globalScale;
                    const radius = node.size / 2 + 4;

                    // --- Circle ---
                    ctx.beginPath();
                    ctx.fillStyle = node === hoverNode ? '#ff8f00' : '#1f78b4'; // Orange on hover
                    ctx.arc(node.x, node.y, radius, 0, 2 * Math.PI, false);
                    ctx.fill();

                    // --- Text Label Above Node ---
                    ctx.textAlign = 'center';
                    ctx.textBaseline = 'bottom'; // Align text to be drawn above the Y coordinate
                    ctx.fillStyle = '#333'; // Dark color for readability
                    ctx.font = `${fontSize}px Sans-Serif`;
                    
                    // Position text above the circle with a small margin
                    ctx.fillText(label, node.x, node.y - radius - 3);
                }}
            />
            {hoverNode && (
                <div style={{
                    position: 'absolute',
                    top: '10px',
                    left: '10px',
                    backgroundColor: 'rgba(0, 0, 0, 0.8)',
                    color: 'white',
                    padding: '8px 12px',
                    borderRadius: '4px',
                    pointerEvents: 'none',
                    fontSize: '14px'
                }}>
                    <h3>{hoverNode.name}</h3>
                    <p>{hoverNode.description}</p>
                </div>
            )}
        </div>
    );
}

export default ForceGraph;