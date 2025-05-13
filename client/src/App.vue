<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from "vue";
import BuzzCanvas from "./components/BuzzCanvas.vue";
import EmojiCanvas from "./components/EmojiCanvas.vue";
import EmojiMenu from "./components/EmojiMenu.vue";
import Footer from "./components/Footer.vue";
import { provideHoverState } from "./composables/useHoverState";
import type { Emoji } from "./constants";

const emojiCanvasRef = ref<InstanceType<typeof EmojiCanvas> | null>(null);

onMounted(() => {
  window.addEventListener("contextmenu", onContextMenu);
});
onBeforeUnmount(() => {
  window.removeEventListener("contextmenu", onContextMenu);
});

const onContextMenu = (e: Event) => e.preventDefault();
const dispatchEmoji = (emoji: Emoji) => {
  if (!emojiCanvasRef.value) return;
  emojiCanvasRef.value.dropEmoji(emoji);
};

provideHoverState();
</script>

<template>
  <main class="h-screen w-screen">
    <BuzzCanvas />
    <EmojiCanvas ref="emojiCanvasRef" />
    <EmojiMenu @dispatch-emoji="dispatchEmoji" />
  </main>
  <Footer />
</template>
