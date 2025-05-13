import { inject, provide, ref, type Ref } from "vue";

const HOVER_KEY = "hover-state";

interface HoverState {
  isHovering: Ref<boolean, boolean>;
  setHovering: (value: boolean) => boolean;
}

export const provideHoverState = () => {
  const isHovering = ref<boolean>(false);
  const setHovering = (value: boolean) => (isHovering.value = value);
  provide<HoverState>(HOVER_KEY, {
    isHovering,
    setHovering,
  });
};

export const useHoverState = (): HoverState => {
  const state = inject<HoverState>(HOVER_KEY);
  if (!state) {
    throw new Error("Hover State is not initialized");
  }
  return state;
};
