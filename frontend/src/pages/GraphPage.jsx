import React from 'react';
import ForceGraph from '../Graph'; // Importing the graph component
import './GraphPage.css'; // Import the new styles

function GraphPage() {
    return (
        <div className="graph-page-container">
            <h1 className="graph-title">Карта Сообществ</h1>
            <div className="graph-wrapper">
                <ForceGraph />
            </div>
        </div>
    );
}

export default GraphPage;