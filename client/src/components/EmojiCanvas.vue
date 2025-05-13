<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useWindowSize } from "@vueuse/core";
import type { Emoji } from "../constants";
import Ball from "../elements/Ball";

const GRAVITY = 0.05;
const FRICTION = 0.65;
const EMOJI_SCALE = 1.9;
const VERTICAL_OFFSET = -0.18;
const BALL_LIFETIME_MS = 5000;

const canvasRef = ref<HTMLCanvasElement | null>(null);
const ctx = ref<CanvasRenderingContext2D | null>();
const animationId = ref<number>();
const emojiBalls = ref<Array<Ball<string>>>([]);

onMounted(() => {
  initCanvas();
});

const { width, height } = useWindowSize();
const isMobile = computed(() => width.value < 800);
const maxElements = computed(() => (isMobile.value ? 50 : 300));
const emojiRadius = computed(() => (isMobile.value ? 20 : 100));

const initCanvas = () => {
  if (!canvasRef.value) {
    console.error("Canvas element not found.");
    return;
  }
  canvasRef.value.width = width.value;
  canvasRef.value.height = height.value;
  ctx.value = canvasRef.value.getContext("2d");
  if (!ctx.value) return;
  render();
};

const clearCanvas = () => {
  if (!canvasRef.value || !ctx.value) return;
  ctx.value.clearRect(0, 0, canvasRef.value.width, canvasRef.value.height);
};
const drawBall = (ball: Ball<string>, debug = false) => {
  if (!ctx.value) return;

  const fontSize = ball.radius * EMOJI_SCALE;
  ctx.value.font = `${fontSize}px Noto Color Emoji`;

  ctx.value.globalAlpha = ball.opacity;

  const metrics = ctx.value.measureText(ball.content);
  let textHeight;
  if (
    metrics.actualBoundingBoxAscent !== undefined &&
    metrics.actualBoundingBoxDescent !== undefined
  ) {
    textHeight =
      metrics.actualBoundingBoxAscent + metrics.actualBoundingBoxDescent;
  } else {
    textHeight = fontSize * 0.7;
  }
  const verticalAdjustment = fontSize * VERTICAL_OFFSET;
  ctx.value.fillText(
    ball.content,
    ball.x - metrics.width / 2,
    ball.y + textHeight / 2 + verticalAdjustment,
  );

  ctx.value.globalAlpha = 1.0;

  // DEBUG
  if (debug) {
    ctx.value.strokeStyle = "red";
    ctx.value.beginPath();
    ctx.value.arc(ball.x, ball.y, ball.radius, 0, Math.PI * 2);
    ctx.value.stroke();
  }
};
const render = () => {
  if (!canvasRef.value || !ctx.value) return;
  clearCanvas();
  const currentTime = Date.now();
  emojiBalls.value = emojiBalls.value.filter((ball) => {
    const age = currentTime - ball.createdAt;
    if (age > BALL_LIFETIME_MS * 0.75) {
      const remainingLifePercentage =
        1 - (age - BALL_LIFETIME_MS * 0.75) / (BALL_LIFETIME_MS * 0.25);
      ball.opacity = Math.max(0, remainingLifePercentage);
    }
    return age < BALL_LIFETIME_MS;
  });
  emojiBalls.value.forEach((ball) => {
    if (ball.y + ball.radius + ball.dy > canvasRef.value!.height) {
      ball.dy = -ball.dy * FRICTION;
    } else {
      ball.dy += GRAVITY;
    }

    if (
      ball.x + ball.radius + ball.dx > canvasRef.value!.width ||
      ball.x - ball.radius + ball.dx < 0
    ) {
      ball.dx = -ball.dx;
    }

    ball.x += ball.dx;
    ball.y += ball.dy;
    drawBall(ball);
  });
  animationId.value = requestAnimationFrame(render);
};

const dropEmoji = (emoji: Emoji) => {
  if (!canvasRef.value || !ctx.value) return;
  if (emojiBalls.value.length > maxElements.value) {
    emojiBalls.value.shift();
  }
  emojiBalls.value.push(
    new Ball<string>({
      radius: emojiRadius.value,
      canvasWidth: canvasRef.value.width,
      canvasHeight: canvasRef.value.height,
      content: emoji.value,
    }),
  );
};
defineExpose({ dropEmoji });
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
