<template>
  <button
    class="
      absolute
      bottom-5
      right-5
      text-t-main
      h-6
      flex
      items-center
      transform
      scale-110
      opacity-40
      hover:opacity-75
      transition-opacity
      duration-300
      outline-none
    "
    @click="showModal = true"
  >
    <p class="mr-2 text-sm font-mono font-medium">{{ $store.state.theme }}</p>
    <Icon class="w-6 h-6" :icon="icons.theme" />
  </button>

  <c-modal
    v-if="showModal"
    closeable
    @close="showModal = false"
    class="max-w-2xl w-10/12 h-5/6"
  >
    <div
      class="flex flex-col overflow-y-scroll h-full scroll pr-2 text-bg-light"
    >
      <button
        v-for="theme in themes"
        :key="theme"
        class="
          flex
          py-1.5
          px-2
          my-1
          rounded-sm
          text-sm
          font-mono font-medium
          opacity-70
          hover:opacity-100
          hover:text-bg-dark
          hover:bg-input-blur-dark
          transition-all
          duration-300
          outline-none
        "
        :class="
          $store.state.theme === theme ? 'opacity-100 text-primary-600' : ''
        "
        @click="changeTheme(theme)"
      >
        {{ theme }}
      </button>
    </div>
  </c-modal>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import themes from "@/utils/themes";

import CModal from "./Modal/CModal.vue";

import { Icon } from "@iconify/vue";
import themeIcon from "@iconify-icons/feather/layout";
import { useStore } from "@/store";

export default defineComponent({
  name: "CThemeSelector",
  components: {
    CModal,
    Icon,
  },
  setup() {
    const store = useStore();

    const showModal = ref(false);

    const changeTheme = (theme: string) => {
      if (!themes.includes(theme)) return;

      const html = document.querySelector("html");
      if (html) {
        html.classList.remove(store.state.theme);
        html.classList.add(theme);

        store.commit("SET_THEME", theme);
      }
    };

    return {
      showModal,

      themes,
      changeTheme,

      icons: {
        theme: themeIcon,
      },
    };
  },
});
</script>

<style lang="scss">
.scroll {
  &::-webkit-scrollbar {
    background: transparent;
  }

  &::-webkit-scrollbar-thumb {
    background: var(--bg-light);
    border-radius: 2rem;
    opacity: 0.6;
  }
}

:root.light {
  --primary-100: theme("colors.primaryColor.100");
  --primary-200: theme("colors.primaryColor.200");
  --primary-300: theme("colors.primaryColor.300");
  --primary-400: theme("colors.primaryColor.400");
  --primary-500: theme("colors.primaryColor.500");
  --primary-600: theme("colors.primaryColor.600");
  --primary-700: theme("colors.primaryColor.700");
  --primary-800: theme("colors.primaryColor.800");
  --primary-900: theme("colors.primaryColor.900");
  --accent-100: theme("colors.accentColor.100");
  --accent-200: theme("colors.accentColor.200");
  --accent-300: theme("colors.accentColor.300");
  --accent-400: theme("colors.accentColor.400");
  --accent-500: theme("colors.accentColor.500");
  --accent-600: theme("colors.accentColor.600");
  --accent-700: theme("colors.accentColor.700");
  --accent-800: theme("colors.accentColor.800");
  --accent-900: theme("colors.accentColor.900");
  --bg-dark: theme("colors.coolGray.800");
  --bg-light: theme("colors.coolGray.50");
  --t-main: theme("colors.coolGray.900");
  --t-sub: theme("colors.coolGray.400");
  --b-light: theme("colors.coolGray.400");
  --b-dark: theme("colors.coolGray.600");
  --b-highlight: var(--primary-400);
  --b-light-dark: theme("colors.coolGray.500");
  --b-dark-dark: theme("colors.coolGray.700");
  --b-highlight-dark: var(--primary-500);
  --input-focus: theme("colors.coolGray.700");
  --input-blur: theme("colors.coolGray.500");
  --input-light: theme("colors.coolGray.300");
  --input-focus-dark: theme("colors.coolGray.200");
  --input-blur-dark: theme("colors.coolGray.400");
  --input-light-dark: theme("colors.coolGray.600");
}

:root.dark {
  --primary-100: theme("colors.primaryColorDark.100");
  --primary-200: theme("colors.primaryColorDark.200");
  --primary-300: theme("colors.primaryColorDark.300");
  --primary-400: theme("colors.primaryColorDark.400");
  --primary-500: theme("colors.primaryColorDark.500");
  --primary-600: theme("colors.primaryColorDark.600");
  --primary-700: theme("colors.primaryColorDark.700");
  --primary-800: theme("colors.primaryColorDark.800");
  --primary-900: theme("colors.primaryColorDark.900");
  --accent-100: theme("colors.accentColorDark.100");
  --accent-200: theme("colors.accentColorDark.200");
  --accent-300: theme("colors.accentColorDark.300");
  --accent-400: theme("colors.accentColorDark.400");
  --accent-500: theme("colors.accentColorDark.500");
  --accent-600: theme("colors.accentColorDark.600");
  --accent-700: theme("colors.accentColorDark.700");
  --accent-800: theme("colors.accentColorDark.800");
  --accent-900: theme("colors.accentColorDark.900");
  --bg-dark: theme("colors.warmGray.50");
  --bg-light: theme("colors.warmGray.900");
  --t-main: theme("colors.warmGray.200");
  --t-sub: theme("colors.warmGray.500");
  --b-light: theme("colors.warmGray.500");
  --b-dark: theme("colors.warmGray.300");
  --b-highlight: var(--primary-500);
  --b-light-dark: theme("colors.warmGray.400");
  --b-dark-dark: theme("colors.warmGray.600");
  --b-highlight-dark: var(--primary-400);
  --input-focus: theme("colors.warmGray.200");
  --input-blur: theme("colors.warmGray.400");
  --input-light: theme("colors.warmGray.800");
  --input-focus-dark: theme("colors.warmGray.700");
  --input-blur-dark: theme("colors.warmGray.500");
  --input-light-dark: theme("colors.warmGray.100");
}
</style>