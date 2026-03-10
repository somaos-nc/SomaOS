import React, { useState } from 'react';
import { FileCode, Play, Terminal, ChevronRight, Save, FolderOpen } from 'lucide-react';

const EXAMPLES = {
  'grovers_search.cljv': `(ns ClojureV.qurq)

(defn-ai grover_oracle [clk rst_n in]
  (let [target 0xABCDEF]
    (if (= in target)
      (qurq/phi-scale out in -1.0)
      (qurq/assign out in))))

(defn-ai grover_diffusion [clk rst_n in]
  (let [mean (qurq/read-average in)]
    (qurq/phi-scale out (- (* 2 mean) in))))

(defn-fractal grovers_search [clk rst_n in]
  (loop [depth 16 signal in]
    (if (zero? depth)
      (qurq/assign out signal)
      (let [marked (grover_oracle clk rst_n signal)]
        (recur (dec depth) (grover_diffusion clk rst_n marked))))))`,

  'shors_factorization.cljv': `(ns ClojureV.qurq)

(defn-ai modular_exponentiation [clk rst_n base exp mod]
  (let [phi_base (qurq/phi-scale base)]
    (qurq/mod-exp out phi_base exp mod)))

(defn-fractal quantum_fourier_transform [clk rst_n in]
  (loop [q_idx 0 data in]
    (if (= q_idx 8)
      (qurq/assign out data)
      (let [h_gate (qurq/hadamard data q_idx)
            cp_gate (qurq/controlled-phase h_gate q_idx)]
        (recur (inc q_idx) cp_gate)))))

(defn-fractal shors_factorization [clk rst_n n]
  (let [a (qurq/select-coprime n)
        x (qurq/superposition 8)
        f_x (modular_exponentiation clk rst_n a x n)]
    (let [qft_result (quantum_fourier_transform clk rst_n f_x)]
      (qurq/collapsed-period out qft_result))))`,

  'bell_state_entanglement.cljv': `(ns ClojureV.qurq)

(defn-ai BellState [clk rst_n in]
  "Creating a maximally entangled state (Psi+) using a sum-split braid"
  (let [h_gate (qurq/hadamard in)
        cnot_gate (qurq/sum-split h_gate in)]
    (qurq/assign out cnot_gate)))`
};

export const ClojureVIDE = ({ onClose }: { onClose: () => void }) => {
  const [activeFile, setActiveFile] = useState('grovers_search.cljv');
  const [code, setCode] = useState(EXAMPLES['grovers_search.cljv']);
  const [terminalOutput, setTerminalOutput] = useState<string[]>(['SomaOS ClojureV IDE v1.2.0 initialized.', 'Ready for topological synthesis...']);
  const [isCompiling, setIsCompiling] = useState(false);

  const handleRun = () => {
    setIsCompiling(true);
    setTerminalOutput(prev => [...prev, `> Compiling ${activeFile}...`]);
    
    // Simulate compilation stages
    setTimeout(() => {
      setTerminalOutput(prev => [...prev, '> Lexer: 104 tokens generated.']);
      setTimeout(() => {
        setTerminalOutput(prev => [...prev, '> Parser: AST generated (Depth: 12).']);
        setTimeout(() => {
          setTerminalOutput(prev => [...prev, '> Verilog Emitter: sphy_core.v synthesized.']);
          setTimeout(() => {
            setTerminalOutput(prev => [...prev, '> DPR Engine: Streaming x8C Virtual Infinity bitstream...']);
            setTimeout(() => {
              setTerminalOutput(prev => [...prev, '[SUCCESS] GHZ state manifest on the Möbius manifold (d=256).']);
              setIsCompiling(false);
            }, 800);
          }, 600);
        }, 600);
      }, 500);
    }, 400);
  };

  const selectFile = (filename: string) => {
    setActiveFile(filename);
    setCode(EXAMPLES[filename as keyof typeof EXAMPLES]);
    setTerminalOutput(prev => [...prev, `> Opened ${filename}`]);
  };

  return (
    <div className="ide-overlay">
      <div className="ide-container">
        <header className="ide-header">
          <div className="flex items-center gap-2">
            <FileCode className="text-blue-400" size={20} />
            <span className="font-mono text-sm font-bold">ClojureV IDE</span>
          </div>
          <div className="flex items-center gap-4">
            <button 
              onClick={handleRun}
              disabled={isCompiling}
              className={`flex items-center gap-1 px-3 py-1 rounded text-xs font-bold transition-all ${isCompiling ? 'bg-gray-700 text-gray-400' : 'bg-green-600 hover:bg-green-500 text-white'}`}
            >
              <Play size={14} fill="currentColor" />
              {isCompiling ? 'COMPILING...' : 'RUN SYNTHESIS'}
            </button>
            <button onClick={onClose} className="text-gray-400 hover:text-white transition-colors">
              ✕
            </button>
          </div>
        </header>

        <div className="ide-main">
          {/* Sidebar */}
          <div className="ide-sidebar">
            <div className="p-2 border-b border-gray-800 flex items-center gap-2 text-gray-400">
              <FolderOpen size={14} />
              <span className="text-[10px] uppercase font-bold tracking-widest">EXAMPLES</span>
            </div>
            <div className="py-2">
              {Object.keys(EXAMPLES).map(filename => (
                <button
                  key={filename}
                  onClick={() => selectFile(filename)}
                  className={`w-full text-left px-4 py-2 text-xs font-mono transition-colors flex items-center gap-2 ${activeFile === filename ? 'bg-blue-900/30 text-blue-400 border-l-2 border-blue-400' : 'text-gray-400 hover:bg-gray-800'}`}
                >
                  <ChevronRight size={10} />
                  {filename}
                </button>
              ))}
            </div>
          </div>

          {/* Editor */}
          <div className="ide-editor-container">
            <div className="ide-editor-tabs">
              <div className="active-tab">
                {activeFile}
              </div>
            </div>
            <textarea
              value={code}
              onChange={(e) => setCode(e.target.value)}
              spellCheck={false}
              className="ide-textarea"
            />
          </div>

          {/* Terminal */}
          <div className="ide-terminal">
            <div className="p-2 bg-black border-b border-gray-800 flex items-center gap-2 text-gray-400">
              <Terminal size={14} />
              <span className="text-[10px] uppercase font-bold tracking-widest">OUTPUT</span>
            </div>
            <div className="p-4 font-mono text-[11px] h-full overflow-y-auto bg-black text-green-500/80">
              {terminalOutput.map((line, i) => (
                <div key={i} className="mb-1">
                  {line}
                </div>
              ))}
              {isCompiling && <div className="animate-pulse">_</div>}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
