import React, { useState, useRef, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom"; // Import for navigation
import ForceGraph2D from "react-force-graph-2d";
import { forceCollide } from "d3-force";

const NODE_RELSIZE = 10;
const FORCE_MANYBODIES_STRENGTH = -100;
const FORCE_COLLIDE_RADIUS = NODE_RELSIZE * 1.5;

function ForceGraph() {
    const graphRef = useRef(null);
    const [graphData, setGraphData] = useState({ nodes: [], links: [] });
    const [hoverNode, setHoverNode] = useState(null);
    const [highlightNodes, setHighlightNodes] = useState(new Set());
    const [highlightLinks, setHighlightLinks] = useState(new Set());
    const [communityImage, setCommunityImage] = useState(null);
    const navigate = useNavigate(); // Get navigate function

    useEffect(() => {
        const img = new Image();
        img.src = '/community.jpg';
        img.onload = () => setCommunityImage(img);
    }, []);

    useEffect(() => {
        fetch('http://localhost:8080/graph-data')
            .then(res => res.json())
            .then(dataFromBackend => {
                // dataFromBackend now contains { nodes: [], links: [] }
                const { nodes: backendNodes, links: backendLinks } = dataFromBackend;

                // Create a map for quick lookup and to add connection count
                const nodeMap = {};
                backendNodes.forEach(node => {
                    nodeMap[node.id] = {
                        id: node.id,
                        name: node.name,
                        description: "Участников: " + node.size,
                        connections: 0,
                        size: node.size,
                    };
                });

                // Calculate connections based on links
                backendLinks.forEach(link => {
                    if (nodeMap[link.id_1]) nodeMap[link.id_1].connections++;
                    if (nodeMap[link.id_2]) nodeMap[link.id_2].connections++;
                });

                const scaleSize = (initialSize, connections) => {
                    const minSize = 8;
                    const maxSize = 30;
                    // Scale based on initial size (subscribers) and connections
                    return Math.min(maxSize, minSize + initialSize * 0.5 + connections * 1);
                };

                const nodes = Object.values(nodeMap).map(node => ({
                    ...node,
                    size: scaleSize(node.size, node.connections),
                }));

                const links = backendLinks.map(link => ({
                    source: link.id_1,
                    target: link.id_2,
                    value: link.common_subscribers,
                }));

                setGraphData({ nodes, links });
            });
    }, []);

    useEffect(() => {
        if (graphRef.current) {
            graphRef.current.d3Force("charge").strength(FORCE_MANYBODIES_STRENGTH);
            graphRef.current.d3Force("collide", forceCollide(FORCE_COLLIDE_RADIUS));
            graphRef.current.d3Force("link").distance(link => 150 / (1 + link.value * 0.1));
        }
    }, [graphData]);

    const handleNodeHover = (node) => {
        setHoverNode(node);
        if (node) {
            const newHighlightNodes = new Set([node]);
            const newHighlightLinks = new Set();
            graphData.links.forEach(link => {
                if (link.source.id === node.id || link.target.id === node.id) {
                    newHighlightLinks.add(link);
                    newHighlightNodes.add(link.source);
                    newHighlightNodes.add(link.target);
                }
            });
            setHighlightNodes(newHighlightNodes);
            setHighlightLinks(newHighlightLinks);
        } else {
            setHighlightNodes(new Set());
            setHighlightLinks(new Set());
        }
    };

    // --- NEW: Click handler for navigation ---
    const handleNodeClick = (node) => {
        navigate(`/community/${node.id}`);
    };

    const nodeCanvasObject = useCallback((node, ctx, globalScale) => {
        const size = node.size;
        const isHighlighted = highlightNodes.size > 0 && !highlightNodes.has(node);
        const opacity = isHighlighted ? 0.1 : 1;

        ctx.globalAlpha = opacity;

        if (highlightNodes.has(node)) {
            ctx.shadowBlur = 20;
            ctx.shadowColor = 'rgba(74, 144, 226, 0.8)';
        } else {
            ctx.shadowBlur = 5;
            ctx.shadowColor = 'rgba(0, 0, 0, 0.2)';
        }

        if (communityImage) {
            ctx.save();
            ctx.beginPath();
            ctx.arc(node.x, node.y, size / 2, 0, 2 * Math.PI, true);
            ctx.clip();
            ctx.drawImage(communityImage, node.x - size / 2, node.y - size / 2, size, size);
            ctx.restore();
        }

        ctx.beginPath();
        ctx.arc(node.x, node.y, size / 2, 0, 2 * Math.PI, false);
        ctx.strokeStyle = `rgba(74, 144, 226, ${opacity})`;
        ctx.lineWidth = 1 / globalScale;
        ctx.stroke();

        ctx.shadowBlur = 0;

        const label = node.name;
        const fontSize = 12 / globalScale;
        ctx.font = `bold ${fontSize}px Sans-Serif`;
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillStyle = `rgba(0, 0, 0, ${opacity})`;
        ctx.fillText(label, node.x, node.y + size / 2 + fontSize * 0.9);
        
        ctx.globalAlpha = 1;
    }, [communityImage, highlightNodes]);

    const linkCanvasObject = useCallback((link, ctx, globalScale) => {
        const isHighlighted = highlightLinks.size > 0 && !highlightLinks.has(link);
        const opacity = isHighlighted ? 0.05 : 0.5;

        ctx.globalAlpha = opacity;
        ctx.lineWidth = 1 / globalScale;
        ctx.strokeStyle = 'grey';

        ctx.beginPath();
        ctx.moveTo(link.source.x, link.source.y);
        ctx.lineTo(link.target.x, link.target.y);
        ctx.stroke();
        
        ctx.globalAlpha = 1;
    }, [highlightLinks]);

    return (
        <div style={{ position: "relative", flexGrow: 1, width: "100%", background: '#f0f2f5', borderRadius: '8px' }}>
            <ForceGraph2D
                ref={graphRef}
                graphData={graphData}
                nodeVal="size"
                nodeCanvasObject={nodeCanvasObject}
                linkCanvasObject={linkCanvasObject}
                linkWidth={0}
                onNodeHover={handleNodeHover}
                onNodeClick={handleNodeClick} // --- ADDED: Click handler ---
                onNodeDragEnd={node => {
                    // --- FIXED: Release node back into simulation ---
                    node.fx = null;
                    node.fy = null;
                    graphRef.current.d3ReheatSimulation(); // Wake up simulation
                }}
                linkDirectionalParticles={link => highlightLinks.has(link) ? 2 : 0}
                linkDirectionalParticleWidth={2}
                linkDirectionalParticleSpeed={() => 0.01}
                cooldownTicks={100}
                onEngineStop={() => graphRef.current.zoomToFit(400, 100)}
            />
            {hoverNode && (
                <div
                    style={{
                        position: "absolute",
                        top: 10,
                        left: 10,
                        backgroundColor: "rgba(255, 255, 255, 0.9)",
                        padding: "10px",
                        borderRadius: 8,
                        boxShadow: "0 4px 12px rgba(0,0,0,0.15)",
                        pointerEvents: "none",
                        fontSize: 14,
                        zIndex: 10,
                    }}
                >
                    <h3 style={{ margin: 0, marginBottom: 5 }}>{hoverNode.name}</h3>
                    <p style={{ margin: 0 }}>{hoverNode.description}</p>
                    <p style={{ margin: 0, marginTop: 5 }}>Связей: {hoverNode.connections}</p>
                </div>
            )}
        </div>
    );
}

export default ForceGraph;