import React from 'react';
import ForceGraph from '../Graph';
import './GraphPage.css';
import Navbar from '../Navbar';

function GraphPage() {
    return (
        <>
            <Navbar />
            <div className="graph-page-container">
                <h1 className="graph-title">Карта Сообществ</h1>
                <div className="graph-wrapper">
                    <ForceGraph />
                </div>
            </div>
        </>
    );
}

export default GraphPage;
