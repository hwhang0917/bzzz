import { onMounted, onUnmounted, shallowRef } from "vue";

export const useWindowSize = () => {
  const width = shallowRef(0);
  const height = shallowRef(0);

  const update = () => {
    if (window) {
      width.value = window.innerWidth;
      height.value = window.innerHeight;
    }
  };

  const onResize = () => {
    update();
  };

  update();
  onMounted(() => {
    window.addEventListener("resize", onResize);
  });
  onUnmounted(() => {
    window.removeEventListener("resize", onResize);
  });

  return { width, height };
};
