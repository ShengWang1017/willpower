<script setup>
import { onMounted, onUnmounted, ref, watch } from 'vue';
import * as THREE from 'three';
import gsap from 'gsap';

const props = defineProps({
  progress: {
    type: Number,
    default: 0,
    validator: (value) => value >= 0 && value <= 100
  },
  size: {
    type: Number,
    default: 120
  }
});

const canvasRef = ref(null);
let scene, camera, renderer, mesh, animationId;
let currentProgress = 0;

const vertexShader = `
  varying vec2 vUv;

  void main() {
    vUv = uv;
    gl_Position = projectionMatrix * modelViewMatrix * vec4(position, 1.0);
  }
`;

const fragmentShader = `
  uniform float uProgress;
  uniform float uTime;
  uniform vec2 uResolution;
  varying vec2 vUv;

  #define PI 3.14159265359

  // 创建圆环
  float ring(vec2 uv, float radius, float thickness) {
    float dist = length(uv);
    return smoothstep(radius + thickness, radius, dist) - smoothstep(radius, radius - thickness, dist);
  }

  // 创建能量流动效果
  float energyFlow(vec2 uv, float angle, float time) {
    float wave = sin(angle * 8.0 + time * 3.0) * 0.5 + 0.5;
    return pow(wave, 2.0);
  }

  void main() {
    vec2 uv = (vUv - 0.5) * 2.0;
    uv.y *= -1.0; // 翻转 Y 轴

    // 计算角度
    float angle = atan(uv.y, uv.x);
    float normalizedAngle = (angle + PI) / (2.0 * PI);

    // 背景环
    float bgRing = ring(uv, 0.8, 0.08);
    vec3 bgColor = vec3(0.1, 0.1, 0.15) * bgRing;

    // 进度环
    float progressRing = ring(uv, 0.8, 0.12);
    float progressMask = step(normalizedAngle, uProgress);

    // 渐变色（绿色 -> 蓝色 -> 紫色）
    vec3 color1 = vec3(0.09, 0.64, 0.29);  // 绿色
    vec3 color2 = vec3(0.15, 0.39, 0.92);  // 蓝色
    vec3 color3 = vec3(0.55, 0.36, 0.96);  // 紫色

    vec3 progressColor = mix(color1, color2, smoothstep(0.0, 0.5, uProgress));
    progressColor = mix(progressColor, color3, smoothstep(0.5, 1.0, uProgress));

    // 添加能量流动
    float flow = energyFlow(uv, normalizedAngle * 2.0 * PI, uTime);
    progressColor += flow * 0.3;

    // 发光效果
    float glow = progressRing * 2.0;
    vec3 glowColor = progressColor * glow;

    // 边缘高光
    float edgeDistance = abs(length(uv) - 0.8);
    float edgeGlow = smoothstep(0.15, 0.0, edgeDistance) * 0.5;

    // 组合颜色
    vec3 finalColor = bgColor;
    finalColor += progressColor * progressRing * progressMask;
    finalColor += glowColor * progressMask * 0.3;
    finalColor += progressColor * edgeGlow * progressMask;

    // 闪烁效果（当接近完成时）
    if (uProgress > 0.9) {
      float sparkle = sin(uTime * 5.0) * 0.3 + 0.7;
      finalColor *= sparkle;
    }

    float alpha = max(bgRing, progressRing * progressMask);
    gl_FragColor = vec4(finalColor, alpha);
  }
`;

const initRing = () => {
  if (!canvasRef.value) return;

  scene = new THREE.Scene();
  camera = new THREE.OrthographicCamera(-1, 1, 1, -1, 0, 1);

  renderer = new THREE.WebGLRenderer({
    canvas: canvasRef.value,
    alpha: true,
    antialias: true
  });
  renderer.setSize(props.size, props.size);
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

  const geometry = new THREE.PlaneGeometry(2, 2);
  const material = new THREE.ShaderMaterial({
    vertexShader,
    fragmentShader,
    uniforms: {
      uProgress: { value: 0 },
      uTime: { value: 0 },
      uResolution: { value: new THREE.Vector2(props.size, props.size) }
    },
    transparent: true
  });

  mesh = new THREE.Mesh(geometry, material);
  scene.add(mesh);

  const animate = () => {
    animationId = requestAnimationFrame(animate);
    mesh.material.uniforms.uTime.value += 0.016;
    renderer.render(scene, camera);
  };

  animate();
};

watch(() => props.progress, (newProgress) => {
  if (!mesh) return;

  gsap.to(mesh.material.uniforms.uProgress, {
    value: newProgress / 100,
    duration: 1.5,
    ease: 'power2.out'
  });
});

onMounted(() => {
  initRing();
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
  <div class="relative inline-block" :style="{ width: size + 'px', height: size + 'px' }">
    <canvas ref="canvasRef" />
    <div class="absolute inset-0 flex items-center justify-center">
      <div class="text-center">
        <div class="text-2xl font-bold text-gray-900">{{ Math.round(progress) }}%</div>
        <div class="text-xs text-gray-500">完成度</div>
      </div>
    </div>
  </div>
</template>
