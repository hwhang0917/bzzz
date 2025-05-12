<script setup lang="ts">
const EMOJIS = ["😃", "😁", "🤣", "😜", "😱", "😒"];
</script>

<template>
  <ul class="fixed bottom-10 mx-auto flex justify-center text-5xl select-none">
    <li v-for="emoji in EMOJIS" class="font-emoji cursor-pointer">
      {{ emoji }}
    </li>
  </ul>
</template>
