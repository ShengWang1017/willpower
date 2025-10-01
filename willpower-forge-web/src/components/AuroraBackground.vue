<script setup>
import { onMounted, onUnmounted, ref } from 'vue';
import * as THREE from 'three';

const props = defineProps({
  intensity: {
    type: Number,
    default: 0.6
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
  uniform float uIntensity;
  varying vec2 vUv;

  // 改进的噪声函数
  vec3 mod289(vec3 x) { return x - floor(x * (1.0 / 289.0)) * 289.0; }
  vec4 mod289(vec4 x) { return x - floor(x * (1.0 / 289.0)) * 289.0; }
  vec4 permute(vec4 x) { return mod289(((x*34.0)+1.0)*x); }
  vec4 taylorInvSqrt(vec4 r) { return 1.79284291400159 - 0.85373472095314 * r; }

  float snoise(vec3 v) {
    const vec2 C = vec2(1.0/6.0, 1.0/3.0);
    const vec4 D = vec4(0.0, 0.5, 1.0, 2.0);

    vec3 i  = floor(v + dot(v, C.yyy));
    vec3 x0 = v - i + dot(i, C.xxx);

    vec3 g = step(x0.yzx, x0.xyz);
    vec3 l = 1.0 - g;
    vec3 i1 = min(g.xyz, l.zxy);
    vec3 i2 = max(g.xyz, l.zxy);

    vec3 x1 = x0 - i1 + C.xxx;
    vec3 x2 = x0 - i2 + C.yyy;
    vec3 x3 = x0 - D.yyy;

    i = mod289(i);
    vec4 p = permute(permute(permute(
              i.z + vec4(0.0, i1.z, i2.z, 1.0))
            + i.y + vec4(0.0, i1.y, i2.y, 1.0))
            + i.x + vec4(0.0, i1.x, i2.x, 1.0));

    float n_ = 0.142857142857;
    vec3 ns = n_ * D.wyz - D.xzx;

    vec4 j = p - 49.0 * floor(p * ns.z * ns.z);

    vec4 x_ = floor(j * ns.z);
    vec4 y_ = floor(j - 7.0 * x_);

    vec4 x = x_ *ns.x + ns.yyyy;
    vec4 y = y_ *ns.x + ns.yyyy;
    vec4 h = 1.0 - abs(x) - abs(y);

    vec4 b0 = vec4(x.xy, y.xy);
    vec4 b1 = vec4(x.zw, y.zw);

    vec4 s0 = floor(b0)*2.0 + 1.0;
    vec4 s1 = floor(b1)*2.0 + 1.0;
    vec4 sh = -step(h, vec4(0.0));

    vec4 a0 = b0.xzyw + s0.xzyw*sh.xxyy;
    vec4 a1 = b1.xzyw + s1.xzyw*sh.zzww;

    vec3 p0 = vec3(a0.xy,h.x);
    vec3 p1 = vec3(a0.zw,h.y);
    vec3 p2 = vec3(a1.xy,h.z);
    vec3 p3 = vec3(a1.zw,h.w);

    vec4 norm = taylorInvSqrt(vec4(dot(p0,p0), dot(p1,p1), dot(p2, p2), dot(p3,p3)));
    p0 *= norm.x;
    p1 *= norm.y;
    p2 *= norm.z;
    p3 *= norm.w;

    vec4 m = max(0.6 - vec4(dot(x0,x0), dot(x1,x1), dot(x2,x2), dot(x3,x3)), 0.0);
    m = m * m;
    return 42.0 * dot(m*m, vec4(dot(p0,x0), dot(p1,x1), dot(p2,x2), dot(p3,x3)));
  }

  // FBM (Fractional Brownian Motion)
  float fbm(vec3 p) {
    float value = 0.0;
    float amplitude = 0.5;
    float frequency = 1.0;

    for(int i = 0; i < 5; i++) {
      value += amplitude * snoise(p * frequency);
      frequency *= 2.0;
      amplitude *= 0.5;
    }

    return value;
  }

  // 创建极光效果
  vec3 aurora(vec2 uv, float time) {
    // 扭曲 UV 坐标创造流动感
    vec2 warpedUv = uv;
    warpedUv.x += fbm(vec3(uv * 2.0, time * 0.2)) * 0.1;
    warpedUv.y += fbm(vec3(uv * 3.0, time * 0.15)) * 0.05;

    // 多层极光
    float aurora1 = fbm(vec3(warpedUv * 1.5, time * 0.3));
    float aurora2 = fbm(vec3(warpedUv * 2.0 + vec2(5.0), time * 0.25));
    float aurora3 = fbm(vec3(warpedUv * 3.0 + vec2(10.0), time * 0.2));

    // 组合极光层
    float combined = (aurora1 + aurora2 * 0.6 + aurora3 * 0.4) * 0.5 + 0.5;

    // 添加波浪形状
    float wave = sin(uv.y * 8.0 + time) * 0.5 + 0.5;
    combined *= wave;

    // 垂直渐变（极光从上方发光）
    float verticalGradient = smoothstep(0.8, 0.2, uv.y);
    combined *= verticalGradient;

    // 极光颜色 - 绿色到蓝色到紫色
    vec3 color1 = vec3(0.2, 0.8, 0.4);   // 亮绿
    vec3 color2 = vec3(0.3, 0.6, 1.0);   // 蓝色
    vec3 color3 = vec3(0.7, 0.4, 1.0);   // 紫色
    vec3 color4 = vec3(1.0, 0.8, 0.3);   // 金色

    vec3 auroraColor = mix(color1, color2, smoothstep(0.0, 0.3, combined));
    auroraColor = mix(auroraColor, color3, smoothstep(0.3, 0.7, combined));
    auroraColor = mix(auroraColor, color4, smoothstep(0.7, 1.0, combined));

    // 添加闪烁星星
    float stars = snoise(vec3(uv * 50.0, time * 0.1));
    stars = pow(max(0.0, stars), 20.0) * 2.0;

    return auroraColor * pow(combined, 1.5) + vec3(stars);
  }

  void main() {
    vec2 uv = vUv;

    vec3 color = aurora(uv, uTime);

    // 应用强度
    color *= uIntensity;

    // 边缘渐暗
    float vignette = smoothstep(1.2, 0.3, length(uv - 0.5));
    color *= vignette;

    gl_FragColor = vec4(color, smoothstep(0.0, 0.5, length(color)));
  }
`;

const initAurora = () => {
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
      uIntensity: { value: props.intensity }
    },
    transparent: true
  });

  mesh = new THREE.Mesh(geometry, material);
  scene.add(mesh);

  const animate = () => {
    animationId = requestAnimationFrame(animate);
    mesh.material.uniforms.uTime.value += 0.008;
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

onMounted(() => {
  initAurora();
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
