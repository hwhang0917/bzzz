<script setup lang="ts">
import { useSound } from "@vueuse/sound";
import { useHoverState } from "../composables/useHoverState";
import { EMOJIS } from "../constants";
import clickSfx from "../assets/click.wav";
import { useWebSocket } from "@vueuse/core";

const { setHovering } = useHoverState();
const { play, stop } = useSound(clickSfx, { volume: 0.2 });
const emit = defineEmits(["dispatchEmoji"]);
const ws = useWebSocket("/ws", {
  onMessage(_, event) {
    const idx = Number(event.data);
    emit("dispatchEmoji", EMOJIS[idx]);
  },
});

const onMouseEnter = () => {
  play({ playbackRate: 1 });
  setHovering(true);
};
const onMouseLeave = () => {
  stop();
  setHovering(false);
};
const click = (idx?: number) => {
  if (!!navigator.vibrate) {
    navigator.vibrate(50);
  }
  play({ playbackRate: 1.5 });
  if (idx) {
    emit("dispatchEmoji", EMOJIS[idx]);
    ws.send(idx.toString());
  }
};
</script>

<template>
  <ul
    class="fixed justify-center rounded-4xl border-2 border-stone-400 bg-neutral-100 p-2 shadow-lg select-none max-sm:top-1/2 max-sm:right-5 max-sm:-translate-y-1/2 md:bottom-10 md:left-1/2 md:flex md:-translate-x-1/2 md:border-4"
  >
    <li v-for="(emoji, idx) in EMOJIS" class="font-emoji text-3xl md:text-5xl">
      <button
        class="my-2 cursor-pointer duration-100 active:scale-150 md:my-0 md:hover:scale-110 md:active:scale-125"
        :aria-label="`Send ${emoji.value} emoji to others (${emoji.descriptor}).`"
        :title="emoji.descriptor"
        @mouseenter="onMouseEnter"
        @mouseleave="onMouseLeave"
        @click="() => click(idx)"
      >
        {{ emoji.value }}
      </button>
    </li>
    <li class="font-mono text-5xl text-stone-400 max-sm:hidden">|</li>
    <li
      class="w-12 font-mono duration-100 active:scale-150 max-sm:hidden md:text-5xl md:hover:scale-110 md:active:scale-125"
      @mouseenter="onMouseEnter"
      @mouseleave="onMouseLeave"
      @click="() => click()"
    >
      <a
        href="https://github.com/hwhang0917/bzzz"
        target="_blank"
        aria-label="Go to Github link"
        title="Github Link"
      >
        <img src="/github-mark.svg" class="w-full" />
      </a>
    </li>
  </ul>
</template>
