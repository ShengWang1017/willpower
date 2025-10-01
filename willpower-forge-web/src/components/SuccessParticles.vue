<script setup>
import { onMounted, onUnmounted, ref, watch } from 'vue';
import * as THREE from 'three';

const props = defineProps({
  trigger: {
    type: Boolean,
    default: false
  },
  color: {
    type: String,
    default: '#16a34a'
  }
});

const canvasRef = ref(null);
let scene, camera, renderer, particles, animationId;

const vertexShader = `
  attribute float size;
  attribute float alpha;
  varying float vAlpha;

  void main() {
    vAlpha = alpha;
    vec4 mvPosition = modelViewMatrix * vec4(position, 1.0);
    gl_PointSize = size * (300.0 / -mvPosition.z);
    gl_Position = projectionMatrix * mvPosition;
  }
`;

const fragmentShader = `
  uniform vec3 uColor;
  varying float vAlpha;

  void main() {
    // 创建圆形粒子
    vec2 coord = gl_PointCoord - vec2(0.5);
    float dist = length(coord);
    float alpha = 1.0 - smoothstep(0.0, 0.5, dist);

    // 添加闪烁效果
    float sparkle = pow(alpha, 3.0);
    vec3 color = uColor + vec3(sparkle * 0.5);

    gl_FragColor = vec4(color, alpha * vAlpha);
  }
`;

const initParticles = () => {
  if (!canvasRef.value) return;

  scene = new THREE.Scene();
  camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
  camera.position.z = 30;

  renderer = new THREE.WebGLRenderer({
    canvas: canvasRef.value,
    alpha: true
  });
  renderer.setSize(window.innerWidth, window.innerHeight);
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

  const handleResize = () => {
    camera.aspect = window.innerWidth / window.innerHeight;
    camera.updateProjectionMatrix();
    renderer.setSize(window.innerWidth, window.innerHeight);
  };

  window.addEventListener('resize', handleResize);

  return () => {
    window.removeEventListener('resize', handleResize);
  };
};

const createExplosion = () => {
  if (!scene) return;

  // 清除旧粒子
  if (particles) {
    scene.remove(particles);
  }

  const particleCount = 200;
  const positions = new Float32Array(particleCount * 3);
  const velocities = [];
  const sizes = new Float32Array(particleCount);
  const alphas = new Float32Array(particleCount);

  for (let i = 0; i < particleCount; i++) {
    const i3 = i * 3;

    // 从中心爆炸
    positions[i3] = 0;
    positions[i3 + 1] = 0;
    positions[i3 + 2] = 0;

    // 随机速度（球形爆炸）
    const theta = Math.random() * Math.PI * 2;
    const phi = Math.random() * Math.PI;
    const speed = Math.random() * 15 + 10;

    velocities.push({
      x: Math.sin(phi) * Math.cos(theta) * speed,
      y: Math.sin(phi) * Math.sin(theta) * speed,
      z: Math.cos(phi) * speed
    });

    sizes[i] = Math.random() * 5 + 2;
    alphas[i] = 1.0;
  }

  const geometry = new THREE.BufferGeometry();
  geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));
  geometry.setAttribute('size', new THREE.BufferAttribute(sizes, 1));
  geometry.setAttribute('alpha', new THREE.BufferAttribute(alphas, 1));

  const color = new THREE.Color(props.color);
  const material = new THREE.ShaderMaterial({
    vertexShader,
    fragmentShader,
    uniforms: {
      uColor: { value: new THREE.Vector3(color.r, color.g, color.b) }
    },
    transparent: true,
    depthWrite: false,
    blending: THREE.AdditiveBlending
  });

  particles = new THREE.Points(geometry, material);
  scene.add(particles);

  // 动画
  let frame = 0;
  const maxFrames = 120;

  const animate = () => {
    if (frame >= maxFrames) {
      cancelAnimationFrame(animationId);
      scene.remove(particles);
      return;
    }

    frame++;
    const positions = particles.geometry.attributes.position.array;
    const alphas = particles.geometry.attributes.alpha.array;

    for (let i = 0; i < particleCount; i++) {
      const i3 = i * 3;
      const velocity = velocities[i];

      // 更新位置
      positions[i3] += velocity.x * 0.1;
      positions[i3 + 1] += velocity.y * 0.1 - 0.3; // 添加重力
      positions[i3 + 2] += velocity.z * 0.1;

      // 速度衰减
      velocity.x *= 0.98;
      velocity.y *= 0.98;
      velocity.z *= 0.98;

      // 透明度衰减
      alphas[i] = 1.0 - (frame / maxFrames);
    }

    particles.geometry.attributes.position.needsUpdate = true;
    particles.geometry.attributes.alpha.needsUpdate = true;

    renderer.render(scene, camera);
    animationId = requestAnimationFrame(animate);
  };

  animate();
};

watch(() => props.trigger, (newVal) => {
  if (newVal) {
    createExplosion();
  }
});

onMounted(() => {
  initParticles();
});

onUnmounted(() => {
  if (animationId) {
    cancelAnimationFrame(animationId);
  }
  if (renderer) {
    renderer.dispose();
  }
});
</script>

<template>
  <canvas
    ref="canvasRef"
    class="fixed inset-0 pointer-events-none"
    style="z-index: 9999;"
  />
</template>
