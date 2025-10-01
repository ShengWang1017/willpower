<script setup>
import { onMounted, onUnmounted, ref } from 'vue';

const props = defineProps({
  glowColor: {
    type: String,
    default: '#16a34a'
  },
  intensity: {
    type: Number,
    default: 0.5
  }
});

const cardRef = ref(null);

onMounted(() => {
  if (!cardRef.value) return;

  // 添加鼠标移动效果
  const handleMouseMove = (e) => {
    const rect = cardRef.value.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;

    cardRef.value.style.setProperty('--mouse-x', `${x}px`);
    cardRef.value.style.setProperty('--mouse-y', `${y}px`);
  };

  cardRef.value.addEventListener('mousemove', handleMouseMove);

  return () => {
    cardRef.value?.removeEventListener('mousemove', handleMouseMove);
  };
});
</script>

<template>
  <div
    ref="cardRef"
    class="glow-card-wrapper"
    :style="{
      '--glow-color': glowColor,
      '--glow-intensity': intensity
    }"
  >
    <div class="glow-card-content">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.glow-card-wrapper {
  position: relative;
  border-radius: 12px;
  padding: 2px;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.1),
    rgba(255, 255, 255, 0)
  );
  --mouse-x: 50%;
  --mouse-y: 50%;
}

.glow-card-wrapper::before {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: 12px;
  padding: 2px;
  background: radial-gradient(
    600px circle at var(--mouse-x) var(--mouse-y),
    var(--glow-color),
    transparent 40%
  );
  -webkit-mask:
    linear-gradient(#fff 0 0) content-box,
    linear-gradient(#fff 0 0);
  -webkit-mask-composite: xor;
  mask-composite: exclude;
  opacity: 0;
  transition: opacity 0.5s;
  pointer-events: none;
}

.glow-card-wrapper:hover::before {
  opacity: var(--glow-intensity);
}

.glow-card-wrapper::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: 12px;
  background: radial-gradient(
    400px circle at var(--mouse-x) var(--mouse-y),
    rgba(255, 255, 255, 0.1),
    transparent 40%
  );
  opacity: 0;
  transition: opacity 0.5s;
  pointer-events: none;
}

.glow-card-wrapper:hover::after {
  opacity: 1;
}

.glow-card-content {
  position: relative;
  background: white;
  border-radius: 11px;
  overflow: hidden;
  z-index: 1;
}

/* 添加动画关键帧 */
@keyframes glow-pulse {
  0%, 100% {
    opacity: var(--glow-intensity);
  }
  50% {
    opacity: calc(var(--glow-intensity) * 1.5);
  }
}

.glow-card-wrapper.pulsing::before {
  animation: glow-pulse 2s ease-in-out infinite;
}
</style>
