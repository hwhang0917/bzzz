<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useWindowSize } from "@vueuse/core";
import type { Emoji } from "../constants";

const canvasRef = ref<HTMLCanvasElement | null>(null);
const ctx = ref<CanvasRenderingContext2D | null>();

onMounted(() => {
  initCanvas();
});

const { width, height } = useWindowSize();

const initCanvas = () => {
  if (!canvasRef.value) {
    console.error("Canvas element not found.");
    return;
  }
  canvasRef.value.width = width.value;
  canvasRef.value.height = height.value;
  ctx.value = canvasRef.value.getContext("2d");
  if (!ctx.value) return;
};

const dropEmoji = (emoji: Emoji) => {
  if (!ctx.value) return;
};
</script>
<template>
  <canvas ref="canvasRef" class="canvas"></canvas>
</template>

<style scoped>
.canvas {
  width: 100%;
  height: 100%;
}
</style>
