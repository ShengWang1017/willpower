<script setup>
import { onMounted, onUnmounted, ref, watch } from 'vue';
import * as THREE from 'three';

const props = defineProps({
  recentPerformance: {
    type: Object,
    default: () => ({
      consecutiveDays: 0,
      recentCompletionRate: 0,
      totalCheckIns: 0
    })
  }
});

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
  uniform float uEnergy;
  uniform float uConsecutive;
  uniform float uCompletionRate;
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

  // 根据状态创建不同的背景
  vec3 createDynamicBackground(vec2 uv, float time) {
    // 调整流动速度（连续天数越多，流动越快）
    float flowSpeed = 0.2 + uEnergy * 0.3;

    // 多层噪声
    float noise1 = snoise(uv * 2.0 + vec2(time * flowSpeed, time * flowSpeed * 0.8));
    float noise2 = snoise(uv * 3.0 - vec2(time * flowSpeed * 0.7, time * flowSpeed * 1.2));
    float noise3 = snoise(uv * 4.0 + vec2(sin(time * 0.1), cos(time * 0.15)));

    float combined = (noise1 + noise2 * 0.5 + noise3 * 0.25) * 0.5 + 0.5;

    // 根据完成率决定颜色方案
    vec3 baseColor, accentColor1, accentColor2, accentColor3;

    if (uCompletionRate >= 0.8) {
      // 高完成率 - 金色/紫色/蓝色（荣耀）
      baseColor = vec3(0.96, 0.62, 0.07);      // 金色
      accentColor1 = vec3(0.55, 0.36, 0.96);   // 紫色
      accentColor2 = vec3(0.15, 0.39, 0.92);   // 蓝色
      accentColor3 = vec3(0.09, 0.64, 0.29);   // 绿色
    } else if (uCompletionRate >= 0.6) {
      // 中高完成率 - 绿色/蓝色（稳定）
      baseColor = vec3(0.09, 0.64, 0.29);      // 绿色
      accentColor1 = vec3(0.15, 0.39, 0.92);   // 蓝色
      accentColor2 = vec3(0.20, 0.75, 0.45);   // 浅绿
      accentColor3 = vec3(0.55, 0.36, 0.96);   // 紫色
    } else if (uCompletionRate >= 0.4) {
      // 中等完成率 - 蓝色/青色（努力中）
      baseColor = vec3(0.15, 0.39, 0.92);      // 蓝色
      accentColor1 = vec3(0.20, 0.60, 0.80);   // 青色
      accentColor2 = vec3(0.55, 0.36, 0.96);   // 紫色
      accentColor3 = vec3(0.96, 0.62, 0.07);   // 橙色
    } else {
      // 低完成率 - 灰蓝/淡紫（需要加油）
      baseColor = vec3(0.25, 0.35, 0.55);      // 灰蓝
      accentColor1 = vec3(0.35, 0.45, 0.65);   // 淡蓝
      accentColor2 = vec3(0.45, 0.35, 0.65);   // 淡紫
      accentColor3 = vec3(0.55, 0.45, 0.75);   // 紫色
    }

    // 混合颜色
    float mixFactor = uv.y + combined * 0.3;
    vec3 color = mix(baseColor, accentColor1, smoothstep(0.0, 0.35, mixFactor));
    color = mix(color, accentColor2, smoothstep(0.35, 0.65, mixFactor));
    color = mix(color, accentColor3, smoothstep(0.65, 1.0, mixFactor));

    // 能量光晕（连续天数越多越亮）
    float glow = pow(combined, 2.0) * (0.2 + uEnergy * 0.4);
    color += glow;

    // 动态波纹（基于连续天数）
    if (uConsecutive > 0.0) {
      float ripple = sin(length(uv - 0.5) * 10.0 - time * 2.0) * 0.5 + 0.5;
      ripple *= smoothstep(0.5, 0.0, length(uv - 0.5));
      color += ripple * uConsecutive * vec3(1.0, 0.9, 0.5) * 0.15;
    }

    return color;
  }

  void main() {
    vec2 uv = vUv;

    vec3 color = createDynamicBackground(uv, uTime);

    // 脉动效果（连续打卡越多，脉动越明显）
    float pulse = sin(uTime * 1.5) * (0.02 + uEnergy * 0.03) + (0.98 - uEnergy * 0.03);
    color *= pulse;

    // 边缘渐暗
    float vignette = smoothstep(1.5, 0.5, length(uv - 0.5));
    color *= vignette;

    gl_FragColor = vec4(color, 0.3);
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
      uEnergy: { value: 0 },
      uConsecutive: { value: 0 },
      uCompletionRate: { value: 0 }
    },
    transparent: true
  });

  mesh = new THREE.Mesh(geometry, material);
  scene.add(mesh);

  updateUniforms();

  const animate = () => {
    animationId = requestAnimationFrame(animate);
    mesh.material.uniforms.uTime.value += 0.01;
    renderer.render(scene, camera);
  };

  animate();

  const handleResize = () => {
    renderer.setSize(window.innerWidth, window.innerHeight);
  };

  window.addEventListener('resize', handleResize);

  return () => {
    window.removeEventListener('resize', handleResize);
  };
};

const updateUniforms = () => {
  if (!mesh) return;

  // 能量值：基于连续天数（0-1范围）
  const energy = Math.min(props.recentPerformance.consecutiveDays / 14, 1.0);

  // 连续效果：归一化连续天数
  const consecutive = Math.min(props.recentPerformance.consecutiveDays / 7, 1.0);

  mesh.material.uniforms.uEnergy.value = energy;
  mesh.material.uniforms.uConsecutive.value = consecutive;
  mesh.material.uniforms.uCompletionRate.value = props.recentPerformance.recentCompletionRate;
};

watch(() => props.recentPerformance, () => {
  updateUniforms();
}, { deep: true });

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
