import React, { useRef } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { OrbitControls, Line, Sphere, Text } from '@react-three/drei';
import * as THREE from 'three';
import { HardwareState } from '../App';

// A single Unbraided NAND Loop
const RingOscillator = ({ position, color, isActive, speed = 1 }: { position: [number, number, number], color: string, isActive: boolean, speed?: number }) => {
  const groupRef = useRef<THREE.Group>(null);
  
  useFrame((state, delta) => {
    if (groupRef.current && isActive) {
      groupRef.current.rotation.x += delta * speed;
      groupRef.current.rotation.y += delta * speed * 0.5;
    }
  });

  return (
    <group ref={groupRef} position={position}>
      {/* The physical ring */}
      <mesh>
        <torusGeometry args={[1, 0.1, 16, 100]} />
        <meshStandardMaterial color={isActive ? color : '#444c5e'} wireframe={true} emissive={isActive ? color : '#000'} emissiveIntensity={0.5} />
      </mesh>
      {/* The electron state traversing */}
      {isActive && (
        <mesh position={[1, 0, 0]}>
          <sphereGeometry args={[0.2, 16, 16]} />
          <meshBasicMaterial color="#FFF" />
        </mesh>
      )}
    </group>
  );
};

const InnerKnot = ({ state }: { state: HardwareState }) => {
  const knotRef = useRef<THREE.Group>(null);

  // The SPHY Phase Field causes the entire macro-cell to pulsate/rotate
  useFrame(() => {
    if (knotRef.current) {
      knotRef.current.rotation.y = state.phase_field;
      // Pulse scale based on thermal load mapping (simulated expansion)
      const scale = 1.0 + (state.thermal_load - 35) * 0.05;
      knotRef.current.scale.set(scale, scale, scale);
    }
  });

  const q0Color = state.q0 ? '#00FF00' : '#FF0000';
  const q1Color = state.q1 ? '#0088FF' : '#AA0000';

  return (
    <group ref={knotRef}>
      {/* Loop 0 (Polar) */}
      <RingOscillator position={[-1.5, 0, 0]} color={q0Color} isActive={state.q0} speed={2} />
      <Text position={[-1.5, -1.8, 0]} fontSize={0.3} color="white">Q0 (Polar)</Text>
      
      {/* Loop 1 (Azimuthal) */}
      <RingOscillator position={[1.5, 0, 0]} color={q1Color} isActive={state.q1} speed={1.5} />
      <Text position={[1.5, -1.8, 0]} fontSize={0.3} color="white">Q1 (Azimuthal)</Text>

      {/* The 4-NAND Braiding Operator (XOR Bridge) */}
      <Line
        points={[[-0.5, 0, 0], [0.5, 0, 0]]}
        color={(state.q0 || state.q1) ? "#FFFF00" : "#5a667d"}
        lineWidth={3}
      />
      <Sphere position={[0, 0, 0]} args={[0.3, 16, 16]}>
        <meshStandardMaterial color={(state.q0 !== state.q1) ? "#FFFF00" : "#5a667d"} emissive={(state.q0 !== state.q1) ? "#FFFF00" : "#000"} emissiveIntensity={0.8} />
      </Sphere>
      <Text position={[0, -0.6, 0]} fontSize={0.2} color="#AAA">4-NAND Bridge</Text>
    </group>
  );
};

export const LogicBlock = ({ state }: { state: HardwareState }) => {
  return (
    <div className="canvas-container">
      <Canvas camera={{ position: [0, 0, 8], fov: 45 }}>
        <color attach="background" args={['#1c2841']} />
        <ambientLight intensity={0.6} />
        <pointLight position={[10, 10, 10]} intensity={2.5} />
        <InnerKnot state={state} />
        <OrbitControls enableZoom={true} />
      </Canvas>
    </div>
  );
};
