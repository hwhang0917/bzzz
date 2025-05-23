<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from "vue";
import { useWindowSize } from "@vueuse/core";
import type { Emoji } from "../constants";
import Ball from "../elements/Ball";

const GRAVITY = 0.15;
const FRICTION = 0.65;
const EMOJI_SCALE = 1.9;
const VERTICAL_OFFSET = -0.18;
const BALL_LIFETIME_MS = 7_000;
const MAX_ELEMENTS = 20;

const SHADOW_COLOR = "rgba(0, 0, 0, 0.3)";
const SHADOW_BLUR = 6;
const SHADOW_OFFSET_X = 2;
const SHADOW_OFFSET_Y = 2;

const canvasRef = ref<HTMLCanvasElement | null>(null);
const ctx = ref<CanvasRenderingContext2D | null>();
const animationId = ref<number | null>(null);
const emojiBalls = ref<Array<Ball<string>>>([]);
const timestamp = ref<number>(0);

onMounted(() => {
  initCanvas();
  document.addEventListener("visibilitychange", handleVisibilityChange);
});
onUnmounted(() => {
  if (animationId.value) {
    cancelAnimationFrame(animationId.value);
    animationId.value = null;
  }
  document.removeEventListener("visibilitychange", handleVisibilityChange);
});

const { width, height } = useWindowSize();
const isMobile = computed(() => width.value < 800);
const emojiRadius = computed(() => (isMobile.value ? 50 : 100));

const resizeCanvas = () => {
  if (!canvasRef.value) return;
  canvasRef.value.width = width.value;
  canvasRef.value.height = height.value;
};
watch([width, height], resizeCanvas);

const initCanvas = () => {
  if (!canvasRef.value) {
    console.error("Canvas element not found.");
    return;
  }
  resizeCanvas();
  ctx.value = canvasRef.value.getContext("2d");
  if (!ctx.value) return;
  render(0);
};

const handleVisibilityChange = () => {
  if (document.hidden && !!animationId.value) {
    cancelAnimationFrame(animationId.value);
    animationId.value = null;
  } else {
    render(0);
  }
};

const clearCanvas = () => {
  if (!canvasRef.value || !ctx.value) return;
  ctx.value.clearRect(0, 0, canvasRef.value.width, canvasRef.value.height);
};
const drawBall = (ball: Ball<string>, debug = false) => {
  if (!ctx.value) return;

  const fontSize = ball.radius * EMOJI_SCALE;
  ctx.value.font = `${fontSize}px Noto Color Emoji`;

  ctx.value.shadowColor = SHADOW_COLOR;
  ctx.value.shadowBlur = SHADOW_BLUR;
  ctx.value.shadowOffsetX = SHADOW_OFFSET_X;
  ctx.value.shadowOffsetY = SHADOW_OFFSET_Y;

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

  ctx.value.save();

  ctx.value.translate(ball.x, ball.y);
  ctx.value.rotate((ball.rotation * Math.PI) / 180);

  const verticalAdjustment = fontSize * VERTICAL_OFFSET;
  ctx.value.fillText(
    ball.content,
    -metrics.width / 2,
    textHeight / 2 + verticalAdjustment,
  );

  ctx.value.restore();

  ctx.value.shadowColor = "transparent";
  ctx.value.shadowBlur = 0;
  ctx.value.shadowOffsetX = 0;
  ctx.value.shadowOffsetY = 0;

  ctx.value.globalAlpha = 1.0;

  // DEBUG
  if (debug) {
    ctx.value.strokeStyle = "red";
    ctx.value.beginPath();
    ctx.value.arc(ball.x, ball.y, ball.radius, 0, Math.PI * 2);
    ctx.value.stroke();
  }
};
const render = (ts: number) => {
  timestamp.value = ts;
  if (!canvasRef.value || !ctx.value) return;
  clearCanvas();
  emojiBalls.value = emojiBalls.value.filter((ball) => {
    const age = ts - ball.createdAt;
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
      ball.dr = ball.dr * FRICTION;
    } else {
      ball.dy += GRAVITY;
    }

    if (
      ball.x + ball.radius + ball.dx > canvasRef.value!.width ||
      ball.x - ball.radius + ball.dx < 0
    ) {
      ball.dx = -ball.dx;
      ball.dr = ball.dr * FRICTION;
    }

    ball.x += ball.dx;
    ball.y += ball.dy;
    ball.rotation += ball.dr;

    drawBall(ball);
  });
  animationId.value = requestAnimationFrame(render);
};

const dropEmoji = (emoji: Emoji) => {
  if (!canvasRef.value || !ctx.value) return;
  if (emojiBalls.value.length > MAX_ELEMENTS) {
    emojiBalls.value.shift();
  }
  emojiBalls.value.push(
    new Ball<string>({
      radius: emojiRadius.value,
      canvasWidth: canvasRef.value.width,
      canvasHeight: canvasRef.value.height,
      content: emoji.value,
      createdAt: timestamp.value,
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
