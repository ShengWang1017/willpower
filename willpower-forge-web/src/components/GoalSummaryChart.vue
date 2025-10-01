<script setup>
import { computed } from 'vue';
import { Bar } from 'vue-chartjs';
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale
} from 'chart.js';

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale);

const props = defineProps({
  summaries: {
    type: Array,
    default: () => []
  }
});

const palette = {
  completed: {
    border: '#16a34a',
    gradientLight: '#bbf7d0',
    gradientDark: '#16a34a',
    hoverBorder: '#15803d'
  },
  partial: {
    border: '#f59e0b',
    gradientLight: '#fde68a',
    gradientDark: '#f59e0b',
    hoverBorder: '#d97706'
  },
  failed: {
    border: '#ef4444',
    gradientLight: '#fecaca',
    gradientDark: '#ef4444',
    hoverBorder: '#dc2626'
  }
};

const labels = computed(() => props.summaries.map((item) => {
  const title = item.title || `Goal #${item.goal_id}`;
  // 截断过长的标题，保留前10个字符
  return title.length > 10 ? title.substring(0, 10) + '...' : title;
}));

const datasets = computed(() => {
  const statuses = ['completed', 'partial', 'failed'];
  return statuses.map((status) => ({
    status,
    data: props.summaries.map((item) => item[status] ?? 0),
    statusKey: status
  }));
});

const chartData = computed(() => ({
  labels: labels.value,
  datasets: datasets.value.map((dataset) => ({
    label: dataset.status.charAt(0).toUpperCase() + dataset.status.slice(1),
    data: dataset.data,
    statusKey: dataset.statusKey,
    borderWidth: 2,
    borderRadius: 8,
    borderColor: palette[dataset.status].border,
    backgroundColor: (context) => {
      const { ctx, chartArea } = context.chart;
      if (!chartArea) {
        return palette[dataset.status].gradientLight;
      }
      const gradient = ctx.createLinearGradient(0, chartArea.top, 0, chartArea.bottom);
      gradient.addColorStop(0, palette[dataset.status].gradientLight);
      gradient.addColorStop(1, palette[dataset.status].gradientDark);
      return gradient;
    },
    hoverBackgroundColor: (context) => {
      const { ctx, chartArea } = context.chart;
      if (!chartArea) {
        return palette[dataset.status].gradientDark;
      }
      const gradient = ctx.createLinearGradient(0, chartArea.top, 0, chartArea.bottom);
      gradient.addColorStop(0, palette[dataset.status].gradientLight);
      gradient.addColorStop(1, palette[dataset.status].gradientDark);
      return gradient;
    },
    hoverBorderColor: palette[dataset.status].hoverBorder,
    hoverBorderWidth: 3,
    maxBarThickness: 40
  }))
}));

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    x: {
      stacked: true,
      grid: {
        display: false
      },
      ticks: {
        color: '#475569',
        font: {
          family: 'Inter, system-ui, sans-serif',
          size: 12,
          weight: '500'
        },
        maxRotation: 0,
        minRotation: 0,
        autoSkip: false
      }
    },
    y: {
      stacked: true,
      beginAtZero: true,
      grid: {
        color: 'rgba(148, 163, 184, 0.2)'
      },
      ticks: {
        precision: 0,
        color: '#475569',
        font: {
          family: 'Inter, system-ui, sans-serif'
        }
      }
    }
  },
  plugins: {
    legend: {
      align: 'start',
      position: 'top',
      labels: {
        usePointStyle: true,
        pointStyle: 'circle',
        padding: 16,
        color: '#1f2937',
        font: {
          family: 'Inter, system-ui, sans-serif',
          size: 13,
          weight: '500'
        },
        boxWidth: 12,
        boxHeight: 12,
        generateLabels: (chart) => {
          const datasets = chart.data.datasets;
          return datasets.map((dataset, i) => {
            const statusKey = dataset.statusKey || Object.keys(palette)[i];
            return {
              text: dataset.label,
              fillStyle: palette[statusKey].border,
              strokeStyle: palette[statusKey].border,
              lineWidth: 0,
              hidden: !chart.isDatasetVisible(i),
              index: i
            };
          });
        }
      }
    },
    tooltip: {
      backgroundColor: '#0f172a',
      titleColor: '#e2e8f0',
      bodyColor: '#cbd5f5',
      cornerRadius: 12,
      padding: 12,
      callbacks: {
        label: (context) => {
          const value = context.parsed.y || 0;
          return `${context.dataset.label}: ${value}`;
        }
      }
    },
    title: {
      display: false
    }
  },
  animation: {
    duration: 1200,
    easing: 'easeOutCubic',
    delay: (context) => context.dataIndex * 100,
    onProgress: (animation) => {
      const progress = animation.currentStep / animation.numSteps;
      if (animation.chart) {
        animation.chart.options.animation.y = {
          from: (ctx) => ctx.chart.height
        };
      }
    }
  },
  interaction: {
    intersect: false,
    mode: 'index'
  }
}));
</script>

<template>
  <div class="h-80 w-full">
    <Bar :data="chartData" :options="chartOptions" />
  </div>
</template>
