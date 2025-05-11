<script setup lang="ts">
import { onBeforeUnmount, onMounted } from "vue";
import Canvas from "./components/Canvas.vue";

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
</script>

<template>
  <main class="h-screen w-screen">
    <Canvas />
  </main>
</template>
