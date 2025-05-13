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
import { useWindowSize } from "@vueuse/core";
import { useSound } from "@vueuse/sound";
import bubbleSfx from "../assets/bubble.wav";
import { useHoverState } from "../composables/useHoverState";

interface Position {
  x: number;
  y: number;
}
interface CanvasProps {
  color?: string;
}
const { color = "#3c40c6" } = defineProps<CanvasProps>();

const isMobile = computed(() => width.value < 800);
const rippleSize = computed(() => (isMobile.value ? 150 : 100));
const rippleIncrement = computed(() => (isMobile.value ? 10 : 5));
const rippleColor = computed(() =>
  new Color(color).fade(0.2).string(),
);
const rippleInterval = 30;

const canvasRef = ref<HTMLCanvasElement | null>(null);
const ctx = ref<CanvasRenderingContext2D | null>();
const animationId = ref<number>();
const isRippling = ref<boolean>(false);
const position = reactive<Position>({ x: 0, y: 0 });
const currentRippleSize = ref<number>(0);
const vibrationInterval = ref<number>();

const { play, stop } = useSound(bubbleSfx, { volume: 0.2 });
const { isHovering: isEmojiMenuHovering } = useHoverState();
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
  console.log(rippleColor.value);
  canvasRef.value.width = width.value;
  canvasRef.value.height = height.value;
  ctx.value = canvasRef.value.getContext("2d");
  if (!ctx.value) return;
  ctx.value.lineCap = "round";
  ctx.value.lineJoin = "round";
  render(0);
};
const updatePosition = (clientX: number, clientY: number) => {
  if (!canvasRef.value) return;
  const rect = canvasRef.value.getBoundingClientRect();
  position.x = clientX - rect.left;
  position.y = clientY - rect.top;
};
const startVibrating = () => {
  if (!navigator.vibrate) return;
  if (vibrationInterval.value) clearInterval(vibrationInterval.value);
  vibrationInterval.value = setInterval(() => {
    navigator.vibrate(rippleInterval);
  }, rippleInterval);
};
const stopVibrating = () => {
  if (!vibrationInterval.value) return;
  clearInterval(vibrationInterval.value);
};
const startRippling = (e: unknown) => {
  if (e instanceof TouchEvent) {
    startVibrating();
    updatePosition(e.touches[0]?.clientX || 0, e.touches[0]?.clientY || 0);
  }
  if (!isRippling.value) play();
  isRippling.value = true;
  currentRippleSize.value = 0;
  lastRippleTime = 0;
};
const stopRippling = (e: unknown) => {
  if (e instanceof TouchEvent) {
    stopVibrating();
  }
  clearCanvas();
  isRippling.value = false;
  currentRippleSize.value = 0;
  stop();
};
const changePosition = (e: MouseEvent | TouchEvent) => {
  clearCanvas();
  if (!canvasRef.value) {
    position.x = 0;
    position.y = 0;
    return;
  }

  let clientX, clientY;
  if (e instanceof MouseEvent) {
    clientX = e.clientX;
    clientY = e.clientY;
  } else {
    clientX = e.touches[0]?.clientX || 0;
    clientY = e.touches[0]?.clientY || 0;
  }

  updatePosition(clientX, clientY);
};
const drawRipple = (size: number) => {
  if (!ctx.value) return;

  ctx.value.beginPath();
  ctx.value.arc(position.x, position.y, size / 2, 0, Math.PI * 2);
  ctx.value.fillStyle = rippleColor.value;
  ctx.value.fill();
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
      currentRippleSize.value += rippleIncrement.value;
      if (currentRippleSize.value >= rippleSize.value) {
        currentRippleSize.value = 0;
      }
      lastRippleTime = timestamp;
    }
    drawRipple(currentRippleSize.value);
  }
  animationId.value = requestAnimationFrame(render);
};

watch(() => width.value, initCanvas);
watch(() => height.value, initCanvas);
watch(() => isEmojiMenuHovering.value, stopRippling);
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
  position: fixed;
  width: 100%;
  height: 100%;
  background-color: transparent;
}
</style>
