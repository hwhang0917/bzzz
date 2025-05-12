<script setup lang="ts">
import Color from "color";
import {
  computed,
  onBeforeUnmount,
  onMounted,
  reactive,
  ref,
  watch,
} from "vue";
import { useWindowSize } from "../composables/useWindowSize";
import clickSfx from "../assets/click.wav";
import { useSound } from "@vueuse/sound";

interface Position {
  x: number;
  y: number;
}
interface CanvasProps {
  rippleColor?: string;
}
const { rippleColor = "#3c40c6" } = defineProps<CanvasProps>();

const rippleSize = 500;
const rippleLevel = 25;
const rippleInterval = 30;

const canvasRef = ref<HTMLCanvasElement | null>(null);
const ctx = ref<CanvasRenderingContext2D | null>();
const animationId = ref<number>();
const isRippling = ref<boolean>(false);
const position = reactive<Position>({ x: 0, y: 0 });
const currentRippleLevel = ref<number>(0);

const { play } = useSound(clickSfx, { volume: 0.2 });

const sizes = computed(() => {
  const increment = rippleSize / (rippleLevel - 1);
  return Array.from({ length: rippleLevel }, (_, i) => {
    return 1 + increment * i;
  });
});
const colors = computed(() => {
  const hsl = Color(rippleColor).hsl();
  return Array.from({ length: rippleLevel }, (_, i) => {
    return hsl.lighten(0.1 * (i + 1)).hex();
  });
});

const { width, height } = useWindowSize();

onMounted(() => {
  initCanvas();
});
onBeforeUnmount(() => {
  if (animationId.value) cancelAnimationFrame(animationId.value);
});

const initCanvas = () => {
  if (!canvasRef.value) {
    console.error("Canvas element not found.");
    return;
  }
  canvasRef.value.width = width.value;
  canvasRef.value.height = height.value;
  ctx.value = canvasRef.value.getContext("2d");
  if (!ctx.value) return;
  ctx.value.lineCap = "round";
  ctx.value.lineJoin = "round";
  render(0);
};
watch(() => width.value, initCanvas);
watch(() => height.value, initCanvas);

const startRippling = () => {
  isRippling.value = true;
  currentRippleLevel.value = 0;
  lastRippleTime = 0;
  play();
};
const stopRippling = () => {
  isRippling.value = false;
  currentRippleLevel.value = 0;
};
const changePosition = (e: MouseEvent | TouchEvent) => {
  clearCanvas();
  if (!canvasRef.value) {
    position.x = 0;
    position.y = 0;
    return;
  }

  const rect = canvasRef.value.getBoundingClientRect();
  let clientX, clientY;
  if (e instanceof MouseEvent) {
    clientX = e.clientX;
    clientY = e.clientY;
  } else {
    clientX = e.touches[0]?.clientX || 0;
    clientY = e.touches[0]?.clientY || 0;
  }

  position.x = clientX - rect.left;
  position.y = clientY - rect.top;
};
const drawRipple = () => {
  if (!ctx.value) return;

  const level = Math.min(currentRippleLevel.value, rippleLevel - 1);
  for (let i = level; i >= 0; i--) {
    ctx.value.beginPath();
    ctx.value.arc(position.x, position.y, sizes.value[i] / 2, 0, Math.PI * 2);
    ctx.value.fillStyle = colors.value[i];
    ctx.value.fill();
  }
};
const clearCanvas = () => {
  if (!canvasRef.value || !ctx.value) return;
  ctx.value.clearRect(0, 0, canvasRef.value.width, canvasRef.value.height);
};
let lastRippleTime = 0;
const render = (timestamp: DOMHighResTimeStamp) => {
  clearCanvas();

  if (ctx.value && isRippling.value) {
    if (timestamp - lastRippleTime > rippleInterval) {
      currentRippleLevel.value++;
      if (currentRippleLevel.value >= rippleLevel) {
        currentRippleLevel.value = 0;
      }
      lastRippleTime = timestamp;
    }
    drawRipple();
  }
  animationId.value = requestAnimationFrame(render);
};
</script>

<template>
  <canvas
    ref="canvasRef"
    class="canvas"
    @mousedown="startRippling"
    @mouseup="stopRippling"
    @mousemove="changePosition"
    @touchstart="startRippling"
    @touchend="stopRippling"
    @touchmove="changePosition"
  ></canvas>
</template>

<style scoped>
.canvas {
  width: 100%;
  height: 100%;
}
</style>
