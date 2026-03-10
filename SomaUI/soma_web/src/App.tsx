import React, { useEffect, useState } from 'react'
import { Dashboard } from './components/Dashboard'
import { ClojureVIDE } from './components/ClojureVIDE'

export type HardwareState = {
  c0: boolean;
  c1: boolean;
  c2: boolean;
  c3: boolean;
  c4: boolean;
  c5: boolean;
  c6: boolean;
  c7: boolean;
  thermal_load: number;
  phase_field: number;
  active_cells: number;
};

function App() {
  const [hwState, setHwState] = useState<HardwareState>({
    c0: false,
    c1: false,
    c2: false,
    c3: false,
    c4: false,
    c5: false,
    c6: false,
    c7: false,
    thermal_load: 35.0,
    phase_field: 0.0,
    active_cells: 8
  });

  const [isIdeOpen, setIsIdeOpen] = useState(false);

  useEffect(() => {
    const fetchState = async () => {
      try {
        const res = await fetch('http://localhost:8081/api/state');
        if (res.ok) {
          const data = await res.json();
          setHwState(data);
        }
      } catch (err) {
        // Silently fail if server is down, keep showing old state
      }
    };

    // Poll server at 10Hz to match hardware simulator
    const interval = setInterval(fetchState, 100);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="app-container">
      <header className="header">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <h1>SomaOS</h1>
            <p>Geometric Virtualization Visualizer | FPGA Hardware Link</p>
          </div>
          <button 
            onClick={() => setIsIdeOpen(true)}
            style={{ 
              backgroundColor: '#2563eb', 
              color: 'white', 
              padding: '0.5rem 1rem', 
              borderRadius: '0.25rem', 
              fontSize: '0.875rem', 
              fontWeight: 'bold',
              border: 'none',
              cursor: 'pointer'
            }}
          >
            OPEN CLOJUREV IDE
          </button>
        </div>
      </header>
      <main className="main-content">
        <Dashboard hwState={hwState} />
      </main>

      {isIdeOpen && <ClojureVIDE onClose={() => setIsIdeOpen(false)} />}
    </div>
  )
}

export default App
