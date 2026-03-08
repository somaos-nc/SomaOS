import React, { useEffect, useState } from 'react'
import { Dashboard } from './components/Dashboard'

export type HardwareState = {
  q0: boolean;
  q1: boolean;
  thermal_load: number;
  phase_field: number;
};

function App() {
  const [hwState, setHwState] = useState<HardwareState>({
    q0: false,
    q1: false,
    thermal_load: 35.0,
    phase_field: 0.0
  });

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
        <h1>SomaOS</h1>
        <p>Geometric Virtualization Visualizer | FPGA Hardware Link</p>
      </header>
      <main className="main-content">
        <Dashboard hwState={hwState} />
      </main>
    </div>
  )
}

export default App
