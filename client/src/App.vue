<template>
  <c-toast-controller />

  <c-theme-selector />

  <router-view
    v-if="$store.state.isConnected"
    class="bg-bg-light"
  ></router-view>
  <div
    class="w-full h-full bg-bg-light flex justify-center items-center"
    v-else
  >
    <c-spinner-circle class="transform scale-75" />
  </div>

  <div id="modals"></div>
</template>

<script lang="ts">
import { defineComponent, onMounted } from "vue";
import themes from "@/utils/themes";
import { useStore } from "./store";

import CToastController from "@/components/shared/Toast/CToastController.vue";
import CSpinnerCircle from "@/components/shared/Spinner/CSpinnerCircle.vue";
import CThemeSelector from "@/components/shared/CThemeSelector.vue";

export default defineComponent({
  name: "App",
  components: {
    CToastController,
    CSpinnerCircle,
    CThemeSelector,
  },
  setup() {
    const store = useStore();

    onMounted(() => {
      const html = document.querySelector("html");
      html?.classList.add(themes[0]);

      store.commit("SET_THEME", themes[0]);
    });
  },
});
</script>

<style lang="scss">
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  width: 100vw;
  height: 100vh;
}

#app {
  width: 100%;
  height: 100%;
}
</style>
