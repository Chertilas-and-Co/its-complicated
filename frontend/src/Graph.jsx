import React, { useState, useRef, useEffect } from "react";
import ForceGraph2D from "react-force-graph-2d";
import { forceCollide } from "d3-force";

const IMAGE_SIZE = 24;
const NODE_RELSIZE = IMAGE_SIZE;
const FORCE_MANYBODIES_STRENGTH = -(IMAGE_SIZE * 4);
const FORCE_COLLIDE_RADIUS = NODE_RELSIZE * 0.5;

const dataFromBackend = [
    { id_1: 0, id_2: 1, subscribers_1: 20, subscribers_2: 20, common_subscribers: 20, name_1: "Node 0", desc_1: "Описание узла 0", name_2: "Node 1", desc_2: "Описание узла 1" },
    { id_1: 0, id_2: 2, subscribers_1: 100, subscribers_2: 30, common_subscribers: 1, name_1: "Node 0", desc_1: "Описание узла 0", name_2: "Node 2", desc_2: "Описание узла 2" },
    { id_1: 2, id_2: 3, subscribers_1: 100, subscribers_2: 20, common_subscribers: 10, name_1: "Node 2", desc_1: "Описание узла 2", name_2: "Node 3", desc_2: "Описание узла 3" },
    { id_1: 3, id_2: 4, subscribers_1: 40, subscribers_2: 20, common_subscribers: 14, name_1: "Node 3", desc_1: "Описание узла 3", name_2: "Node 4", desc_2: "Описание узла 4" },
];

// Формируем узлы
const nodeMap = {};
dataFromBackend.forEach(({ id_1, id_2, name_1, desc_1, name_2, desc_2 }) => {
    if (!nodeMap[id_1]) nodeMap[id_1] = { id: id_1, name: name_1, description: desc_1, connections: 0 };
    if (!nodeMap[id_2]) nodeMap[id_2] = { id: id_2, name: name_2, description: desc_2, connections: 0 };
    nodeMap[id_1].connections++;
    nodeMap[id_2].connections++;
});

// Считаем totalK для масштабирования узлов (здесь вместо totalK можно использовать connections)
Object.values(nodeMap).forEach(node => {
    let totalK = 0;
    dataFromBackend.forEach(({ id_1, id_2 }) => {
        if (id_1 === node.id || id_2 === node.id) totalK++;
    });
    node.totalK = totalK;
});

const uniqueNodes = Object.values(nodeMap);

// Функция масштабирования размера узлов
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

// Параметры для расстояния ребер
const BASE_DISTANCE = 1000;
const ALPHA = 10;

const links = dataFromBackend.map(({ id_1, id_2, subscribers_1, subscribers_2, common_subscribers }) => {
    const subscribers_sum = subscribers_1 + subscribers_2;
    const common_ratio = subscribers_sum > 0 ? common_subscribers / subscribers_sum : 0;
    const distance = BASE_DISTANCE / (1 + ALPHA * common_ratio); // расстояние зависит от общего числа подписчиков
    return {
        source: id_1,
        target: id_2,
        color: "#000",
        distance,
    };
});

const graphData = { nodes, links };

function ForceGraph() {
    const graphRef = useRef(null);
    const [hoverNode, setHoverNode] = useState(null);

    useEffect(() => {
        if (graphRef.current) {
            graphRef.current
                .d3Force("link").distance(link => link.distance).iterations(1);

            graphRef.current
                .d3Force("charge")
                .strength(node => -Math.abs(FORCE_MANYBODIES_STRENGTH) / (1 + node.totalK));

            graphRef.current.d3Force(
                "collide",
                forceCollide(FORCE_COLLIDE_RADIUS).strength(0.1).iterations(1)
            );

            graphRef.current.d3Force("center", null);
            graphRef.current.d3ReheatSimulation();
        }
    }, []);

    const handleNodeDragStart = (node) => {
        if (!node) return;
        node.fx = node.x;
        node.fy = node.y;
    };

    const handleNodeDrag = (node) => {
        if (!node) return;
        node.fx = node.x;
        node.fy = node.y;
    };

    const handleNodeDragEnd = (node) => {
        if (!node) return;
        node.fx = null;
        node.fy = null;
    };

    return (
        <div style={{ position: "relative", height: 600, width: "100%" }}>
            <ForceGraph2D
                ref={graphRef}
                graphData={graphData}
                nodeVal={(node) => node.size || IMAGE_SIZE}
                linkCurvature="curvature"
                linkColor="color"
                linkWidth={2}
                linkOpacity={1}
                onNodeHover={(node) => setHoverNode(node)}
                onNodeDragStart={handleNodeDragStart}
                onNodeDrag={handleNodeDrag}
                onNodeDragEnd={handleNodeDragEnd}
                nodeCanvasObject={(node, ctx, globalScale) => {
                    const radius = 12 / globalScale;
                    const label = node.id.toString();
                    const fontSize = 12 / globalScale;

                    ctx.beginPath();
                    ctx.fillStyle = "#1f78b4";
                    ctx.strokeStyle = "#fff";
                    ctx.lineWidth = 1 / globalScale;
                    ctx.arc(node.x, node.y, radius, 0, 2 * Math.PI, false);
                    ctx.fill();
                    ctx.stroke();

                    ctx.fillStyle = "#fff";
                    ctx.font = `${fontSize}px Sans-Serif`;
                    ctx.textAlign = "center";
                    ctx.textBaseline = "middle";
                    ctx.fillText(label, node.x, node.y);

                    ctx.fillStyle = "#000";
                    ctx.font = `${fontSize}px Sans-Serif`;
                    ctx.textAlign = "center";
                    ctx.textBaseline = "bottom";
                    ctx.fillText(node.name, node.x, node.y - radius - 4);
                }}
                linkCanvasObjectMode={() => "replace"}
                linkCanvasObject={(link, ctx, globalScale) => {
                    if (typeof link.source === "string") return;

                    const src = link.source;
                    const tgt = link.target;

                    const radius = 12 / globalScale;
                    const dx = tgt.x - src.x;
                    const dy = tgt.y - src.y;
                    const dist = Math.sqrt(dx * dx + dy * dy);

                    const normX = dx / dist;
                    const normY = dy / dist;

                    const sourceX = src.x + normX * radius;
                    const sourceY = src.y + normY * radius;
                    const targetX = tgt.x - normX * radius;
                    const targetY = tgt.y - normY * radius;

                    ctx.beginPath();
                    ctx.lineWidth = 2 / globalScale;
                    ctx.strokeStyle = link.color || "#000";
                    ctx.moveTo(sourceX, sourceY);
                    ctx.lineTo(targetX, targetY);
                    ctx.lineJoin = "round";
                    ctx.lineCap = "round";
                    ctx.shadowBlur = 2;
                    ctx.shadowColor = link.color || "#000";
                    ctx.stroke();
                }}
            />
            {hoverNode && hoverNode.description && (
                <div
                    style={{
                        position: "absolute",
                        top: 100,
                        left: (hoverNode.x || 0) + 10,
                        backgroundColor: "rgba(255, 255, 255, 0.9)",
                        padding: "6px 8px",
                        borderRadius: 4,
                        boxShadow: "0 2px 5px rgba(0,0,0,0.15)",
                        pointerEvents: "none",
                        whiteSpace: "nowrap",
                        fontSize: 12,
                        zIndex: 10,
                        transform: "translate(-50%, -100%)",
                    }}
                >
                    {hoverNode.description}
                </div>
            )}
        </div>
    );
}

export default ForceGraph;
