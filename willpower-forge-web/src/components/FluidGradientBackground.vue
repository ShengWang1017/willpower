<script setup>
import { onMounted, onUnmounted, ref } from 'vue';
import * as THREE from 'three';

const canvasRef = ref(null);
let scene, camera, renderer, mesh, animationId;

const vertexShader = `
  varying vec2 vUv;

  void main() {
    vUv = uv;
    gl_Position = projectionMatrix * modelViewMatrix * vec4(position, 1.0);
  }
`;

const fragmentShader = `
  uniform float uTime;
  uniform vec2 uResolution;
  varying vec2 vUv;

  // 噪声函数
  vec3 mod289(vec3 x) { return x - floor(x * (1.0 / 289.0)) * 289.0; }
  vec2 mod289(vec2 x) { return x - floor(x * (1.0 / 289.0)) * 289.0; }
  vec3 permute(vec3 x) { return mod289(((x*34.0)+1.0)*x); }

  float snoise(vec2 v) {
    const vec4 C = vec4(0.211324865405187, 0.366025403784439, -0.577350269189626, 0.024390243902439);
    vec2 i  = floor(v + dot(v, C.yy));
    vec2 x0 = v - i + dot(i, C.xx);
    vec2 i1 = (x0.x > x0.y) ? vec2(1.0, 0.0) : vec2(0.0, 1.0);
    vec4 x12 = x0.xyxy + C.xxzz;
    x12.xy -= i1;
    i = mod289(i);
    vec3 p = permute(permute(i.y + vec3(0.0, i1.y, 1.0)) + i.x + vec3(0.0, i1.x, 1.0));
    vec3 m = max(0.5 - vec3(dot(x0,x0), dot(x12.xy,x12.xy), dot(x12.zw,x12.zw)), 0.0);
    m = m*m;
    m = m*m;
    vec3 x = 2.0 * fract(p * C.www) - 1.0;
    vec3 h = abs(x) - 0.5;
    vec3 ox = floor(x + 0.5);
    vec3 a0 = x - ox;
    m *= 1.79284291400159 - 0.85373472095314 * (a0*a0 + h*h);
    vec3 g;
    g.x  = a0.x  * x0.x  + h.x  * x0.y;
    g.yz = a0.yz * x12.xz + h.yz * x12.yw;
    return 130.0 * dot(m, g);
  }

  // 创建流动的渐变
  vec3 createFluidGradient(vec2 uv, float time) {
    // 多层噪声创造流动效果
    float noise1 = snoise(uv * 2.0 + vec2(time * 0.3, time * 0.2));
    float noise2 = snoise(uv * 3.0 - vec2(time * 0.2, time * 0.4));
    float noise3 = snoise(uv * 4.0 + vec2(sin(time * 0.1), cos(time * 0.15)));

    float combined = (noise1 + noise2 * 0.5 + noise3 * 0.25) * 0.5 + 0.5;

    // 意志力主题色 - 从绿色到蓝色到紫色的渐变
    vec3 color1 = vec3(0.09, 0.64, 0.29);  // 绿色 #16a34a
    vec3 color2 = vec3(0.15, 0.39, 0.92);  // 蓝色 #2563eb
    vec3 color3 = vec3(0.55, 0.36, 0.96);  // 紫色 #8b5cf6
    vec3 color4 = vec3(0.96, 0.62, 0.07);  // 橙色 #f59e0b

    // 基于位置和噪声混合颜色
    float mixFactor = uv.y + combined * 0.3;
    vec3 color = mix(color1, color2, smoothstep(0.0, 0.4, mixFactor));
    color = mix(color, color3, smoothstep(0.4, 0.7, mixFactor));
    color = mix(color, color4, smoothstep(0.7, 1.0, mixFactor));

    // 添加光晕效果
    float glow = pow(combined, 2.0) * 0.3;
    color += glow;

    return color;
  }

  void main() {
    vec2 uv = vUv;

    // 创建流动效果
    vec3 color = createFluidGradient(uv, uTime);

    // 添加脉动效果（象征意志力的律动）
    float pulse = sin(uTime * 2.0) * 0.02 + 0.98;
    color *= pulse;

    // 边缘渐暗
    float vignette = smoothstep(1.5, 0.5, length(uv - 0.5));
    color *= vignette;

    gl_FragColor = vec4(color, 0.25);
  }
`;

const initShader = () => {
  if (!canvasRef.value) return;

  scene = new THREE.Scene();
  camera = new THREE.OrthographicCamera(-1, 1, 1, -1, 0, 1);

  renderer = new THREE.WebGLRenderer({
    canvas: canvasRef.value,
    alpha: true,
    antialias: true
  });
  renderer.setSize(window.innerWidth, window.innerHeight);
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

  const geometry = new THREE.PlaneGeometry(2, 2);
  const material = new THREE.ShaderMaterial({
    vertexShader,
    fragmentShader,
    uniforms: {
      uTime: { value: 0 },
      uResolution: { value: new THREE.Vector2(window.innerWidth, window.innerHeight) }
    },
    transparent: true
  });

  mesh = new THREE.Mesh(geometry, material);
  scene.add(mesh);

  const animate = () => {
    animationId = requestAnimationFrame(animate);
    mesh.material.uniforms.uTime.value += 0.01;
    renderer.render(scene, camera);
  };

  animate();

  const handleResize = () => {
    renderer.setSize(window.innerWidth, window.innerHeight);
    mesh.material.uniforms.uResolution.value.set(window.innerWidth, window.innerHeight);
  };

  window.addEventListener('resize', handleResize);

  return () => {
    window.removeEventListener('resize', handleResize);
  };
};

onMounted(() => {
  initShader();
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
    style="z-index: 0;"
  />
</template>
