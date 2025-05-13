<script setup lang="ts">
import { onBeforeUnmount, onMounted } from "vue";
import BuzzCanvas from "./components/BuzzCanvas.vue";
import EmojiCanvas from "./components/EmojiCanvas.vue";
import EmojiMenu from "./components/EmojiMenu.vue";
import { provideHoverState } from "./composables/useHoverState";
import type { Emoji } from "./constants";

const onContextMenu = (e: Event) => e.preventDefault();

onMounted(() => {
  if (typeof window !== "undefined") {
    window.addEventListener("contextmenu", onContextMenu);
  }
});
onBeforeUnmount(() => {
  if (typeof window !== "undefined") {
    window.removeEventListener("contextmenu", onContextMenu);
  }
});

const dispatchEmoji = (emoji: Emoji) => {
  console.log(emoji.value);
};

provideHoverState();
</script>

<template>
  <main class="h-screen w-screen">
    <BuzzCanvas />
    <EmojiCanvas />
    <EmojiMenu @dispatch-emoji="dispatchEmoji" />
  </main>
</template>
