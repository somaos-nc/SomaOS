import React from 'react';
import { HardwareState } from '../App';
import { LogicBlock } from './LogicBlock';
import { TopologicalFlowWindow } from './TopologicalFlowWindow';
import { Activity, Thermometer, Cpu, Radio, Plus, Minus, BarChart3, Database, Binary } from 'lucide-react';

export const Dashboard = ({ hwState }: { hwState: HardwareState }) => {

  const triggerDPR = async (action: 'spawn' | 'collapse') => {
    try {
      await fetch('http://localhost:8081/api/dpr', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ action })
      });
    } catch (e) {
      console.error("DPR Injection Failed", e);
    }
  };

  // Convert histogram map to sorted array for display
  const histogramData = Object.entries(hwState.state_histogram || {})
    .sort((a, b) => b[1] - a[1])
    .slice(0, 8); // Top 8 states

  // Format the 64-bit register as a binary string
  const ghzStateString = hwState.register.toString(2).padStart(hwState.active_cells, '0');

  return (
    <div className="dashboard">
      <div className="telemetry-panel">
        <div className="stat-card">
          <div className="stat-header">
            <Thermometer size={20} className="text-red-400" />
            <h3>XADC Thermal Load</h3>
          </div>
          <div className="stat-value">{hwState.thermal_load.toFixed(2)} °C</div>
        </div>

        <div className="stat-card">
          <div className="stat-header">
            <Radio size={20} className="text-green-400" />
            <h3>SPHY Phase Tuner</h3>
          </div>
          <div className="stat-value">Φ {(hwState.phase_field / Math.PI).toFixed(3)} π rad</div>
        </div>

        <div className="stat-card">
          <div className="stat-header">
            <Cpu size={20} className="text-blue-400" />
            <h3>Topology State (d=2^{hwState.active_cells})</h3>
          </div>
          <div className="stat-value" style={{ fontSize: '0.8rem', fontWeight: 'bold', wordBreak: 'break-all', fontFamily: 'monospace' }}>
            |{ghzStateString}⟩ GHZ
          </div>
        </div>
        
        <div className="stat-card dpr-panel">
          <div className="stat-header">
            <Activity size={18} className="text-yellow-400" />
            <h3>DPR Intent Engine</h3>
          </div>
          <div className="dpr-controls" style={{ display: 'flex', gap: '10px', marginTop: '10px' }}>
             <button onClick={() => triggerDPR('spawn')} className="dpr-btn" style={{flex: 1, padding: '8px', background: '#1e293b', color: '#00ffcc', border: '1px solid #00ffcc', cursor: 'pointer', display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '5px', borderRadius: '4px'}}>
               <Plus size={14} /> Spawn
             </button>
             <button onClick={() => triggerDPR('collapse')} className="dpr-btn" disabled={hwState.active_cells <= 8} style={{flex: 1, padding: '8px', background: '#1e293b', color: '#ff5555', border: '1px solid #ff5555', cursor: 'pointer', display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '5px', opacity: hwState.active_cells <= 8 ? 0.5 : 1, borderRadius: '4px'}}>
               <Minus size={14} /> Collapse
             </button>
          </div>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 350px', gap: '2rem', flex: 1, minHeight: '500px' }}>
        <div className="visualization-panel" style={{ position: 'relative' }}>
          <h2>{hwState.routing_mode === 'station' ? 'Electronic Hyper-Braid: 64-Cell Station' : 'Electronic Entanglement Bus: 8-Cell Macro-Cube'}</h2>
          <LogicBlock state={hwState} />
          <TopologicalFlowWindow state={hwState} />
        </div>

        {/* Topological Diagnostics Panel */}
        <div className="visualization-panel" style={{ background: '#0f172a', border: '1px solid #1e293b' }}>
          <h2 style={{ color: '#94a3b8', fontSize: '0.9rem', textAlign: 'left', padding: '1rem', borderBottom: '1px solid #1e293b' }}>
            <div className="flex items-center gap-2">
              <BarChart3 size={16} /> TOPOLOGICAL MANIFOLD DIAGNOSTICS
            </div>
          </h2>
          
          <div className="p-4 flex flex-col gap-6">
            {/* Entropy & Coherence metrics */}
            <div className="grid grid-cols-2 gap-4">
              <div className="bg-black/40 p-3 rounded border border-gray-800">
                <div className="text-[10px] text-gray-500 font-bold uppercase mb-1">Shannon Entropy (H)</div>
                <div className="text-xl font-mono text-blue-300">{hwState.shannon_entropy?.toFixed(4)}</div>
              </div>
              <div className="bg-black/40 p-3 rounded border border-gray-800">
                <div className="text-[10px] text-gray-500 font-bold uppercase mb-1">Coherence (T2)</div>
                <div className="text-xl font-mono text-green-300">{hwState.coherence_time?.toFixed(2)} ms</div>
              </div>
            </div>

            {/* Histogram */}
            <div>
              <div className="text-[10px] text-gray-500 font-bold uppercase mb-3 flex items-center gap-2">
                <Database size={12}/> Topological State Distribution (Top 8)
              </div>
              <div className="flex flex-col gap-2">
                {histogramData.length > 0 ? histogramData.map(([state, count]) => (
                  <div key={state} className="flex items-center gap-2">
                    <span className="font-mono text-[10px] text-gray-400 w-8">{state}</span>
                    <div className="flex-1 h-2 bg-gray-900 rounded-full overflow-hidden">
                      <div 
                        className="h-full bg-blue-500/60 shadow-[0_0_8px_rgba(59,130,246,0.5)]" 
                        style={{ width: `${Math.min(100, (count / 128) * 100)}%`, transition: 'width 0.3s ease' }}
                      />
                    </div>
                    <span className="font-mono text-[10px] text-gray-500 w-6 text-right">{count}</span>
                  </div>
                )) : (
                  <div className="text-center py-8 text-gray-600 text-xs italic">Awaiting manifold telemetry...</div>
                )}
              </div>
            </div>

            {/* Manifold Hex Sample */}
            <div className="mt-2">
              <div className="text-[10px] text-gray-500 font-bold uppercase mb-2 flex items-center gap-2">
                <Binary size={12}/> Raw Manifold Vector (Sample)
              </div>
              <div className="bg-black p-2 rounded font-mono text-[9px] text-blue-400/70 break-all leading-tight border border-blue-900/20">
                {hwState.manifold ? `${hwState.manifold.substring(0, 128)}...` : '0x00000000000000...'}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};