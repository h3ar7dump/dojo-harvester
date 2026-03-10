import React, { useRef } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { useTelemetryStore } from '@/stores/telemetryStore';
import * as THREE from 'three';

// A mock robot visualization. In a real app we'd load URDF via urdf-loader
const MockRobot: React.FC = () => {
  const meshRef = useRef<THREE.Mesh>(null);
  const pose = useTelemetryStore(state => state.pose);

  useFrame(() => {
    if (meshRef.current) {
      if (pose) {
        meshRef.current.position.set(pose.x, pose.y, pose.z);
        meshRef.current.quaternion.set(pose.qx, pose.qy, pose.qz, pose.qw);
      } else {
        // Rotate slowly if no pose data
        meshRef.current.rotation.y += 0.01;
      }
    }
  });

  return (
    <mesh ref={meshRef}>
      <boxGeometry args={[1, 2, 1]} />
      <meshStandardMaterial color="orange" />
    </mesh>
  );
};

export const RobotVisualization: React.FC = () => {
  return (
    <div className="w-full h-[400px] bg-slate-900 rounded-lg overflow-hidden border">
      <Canvas camera={{ position: [5, 5, 5] }}>
        <ambientLight intensity={0.5} />
        <pointLight position={[10, 10, 10]} intensity={1} />
        <MockRobot />
        <gridHelper args={[10, 10]} />
        <axesHelper args={[10]} />
      </Canvas>
    </div>
  );
};
