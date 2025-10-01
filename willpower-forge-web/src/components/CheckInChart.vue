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
  checkIns: {
    type: Array,
    required: true
  }
});

const statusPalette = {
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

const aggregated = computed(() => {
  const result = new Map();
  [...props.checkIns]
    .sort((a, b) => (a.date < b.date ? -1 : a.date > b.date ? 1 : 0))
    .forEach((item) => {
      const dateKey = item.date;
      if (!result.has(dateKey)) {
        result.set(dateKey, { completed: 0, partial: 0, failed: 0 });
      }
      const group = result.get(dateKey);
      if (group[item.status] !== undefined) {
        group[item.status] += 1;
      }
    });

  const labels = Array.from(result.keys());
  const totals = labels.map((label) => result.get(label));

  return {
    labels,
    datasets: ['completed', 'partial', 'failed'].map((status) => ({
      status,
      statusKey: status,
      data: totals.map((entry) => entry[status])
    }))
  };
});

const chartData = computed(() => ({
  labels: aggregated.value.labels,
  datasets: aggregated.value.datasets.map((dataset) => ({
    label: dataset.status.charAt(0).toUpperCase() + dataset.status.slice(1),
    data: dataset.data,
    statusKey: dataset.statusKey,
    borderWidth: 2,
    borderRadius: 8,
    borderColor: statusPalette[dataset.status].border,
    backgroundColor: (context) => {
      const { ctx, chartArea } = context.chart;
      if (!chartArea) {
        return statusPalette[dataset.status].gradientLight;
      }
      const gradient = ctx.createLinearGradient(0, chartArea.top, 0, chartArea.bottom);
      gradient.addColorStop(0, statusPalette[dataset.status].gradientLight);
      gradient.addColorStop(1, statusPalette[dataset.status].gradientDark);
      return gradient;
    },
    hoverBackgroundColor: (context) => {
      const { ctx, chartArea } = context.chart;
      if (!chartArea) {
        return statusPalette[dataset.status].gradientDark;
      }
      const gradient = ctx.createLinearGradient(0, chartArea.top, 0, chartArea.bottom);
      gradient.addColorStop(0, statusPalette[dataset.status].gradientLight);
      gradient.addColorStop(1, statusPalette[dataset.status].gradientDark);
      return gradient;
    },
    hoverBorderColor: statusPalette[dataset.status].hoverBorder,
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
        color: '#1f2937',
        font: {
          family: 'Inter, system-ui, sans-serif',
          size: 12
        }
      }
    },
    y: {
      stacked: true,
      beginAtZero: true,
      grid: {
        color: 'rgba(148, 163, 184, 0.18)'
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
            const statusKey = dataset.statusKey || Object.keys(statusPalette)[i];
            return {
              text: dataset.label,
              fillStyle: statusPalette[statusKey].border,
              strokeStyle: statusPalette[statusKey].border,
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
    duration: 900,
    easing: 'easeOutCubic',
    delay: (context) => context.dataIndex * 80
  }
}));
</script>

<template>
  <div class="h-80 w-full">
    <Bar :data="chartData" :options="chartOptions" />
  </div>
</template>
