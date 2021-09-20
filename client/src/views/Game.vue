<template>
  <main class="w-full h-full flex flex-col justify-center items-center">
    <c-modal v-if="waitingForOpponent" class="flex flex-col">
      <h1 class="font-bold text-3xl">Waiting for opponent...</h1>
      <div class="flex mt-5 mx-auto">
        <c-input-text
          dark
          classes="w-28 "
          class="text-center font-mono text-lg cursor-pointer"
          disabled
          name="code"
          v-model="code"
        />
        <c-button class="ml-2" @click="copyCode">
          <Icon class="w-6 h-6" :icon="icons.copy" />
        </c-button>
      </div>
    </c-modal>

    <c-modal v-if="joiningGame" class="flex flex-col">
      <h1 class="font-bold text-3xl">Joining game...</h1>
    </c-modal>

    <!-- OPPONENT CARD -->
    <div></div>

    <c-game-board />

    <!-- PLAYER CARD -->
    <div></div>
  </main>
</template>

<script lang="ts">
import { computed, defineComponent, onBeforeMount } from "vue";
import { useStore } from "@/store";
import { useRouter } from "vue-router";

import CModal from "@/components/shared/Modal/CModal.vue";
import CInputText from "@/components/shared/Input/CInputText.vue";
import CButton from "@/components/shared/Button/CButton.vue";

import { Icon } from "@iconify/vue";
import copyIcon from "@iconify-icons/feather/copy";
import CGameBoard from "@/components/app/Game/CGameBoard.vue";

export default defineComponent({
  name: "Game",
  components: {
    CModal,
    CInputText,
    CButton,
    Icon,
    CGameBoard,
  },
  setup() {
    const store = useStore();
    const router = useRouter();

    onBeforeMount(() => {
      if (!store.state.inGame) {
        router.push({ name: "Home" });
      }
    });

    const code = router.currentRoute.value.query.code?.toString() as string;
    const isOpponent = router.currentRoute.value.query.opponent;

    const copyCode = () => {
      navigator.clipboard
        .writeText(code)
        .then(() => {
          store.commit("ADD_TOAST", {
            text: "Copied to clipboard!",
            duration: 1500,
          });
        })
        .catch(() => {
          store.commit("ADD_TOAST", {
            text: "Error while copying to clipboard.",
            duration: 1500,
          });
        });
    };

    const waitingForOpponent = computed(() => {
      return store.state.inGame && !store.state.game && !isOpponent;
    });

    const joiningGame = computed(() => {
      return store.state.inGame && !store.state.game && isOpponent;
    });

    return {
      code,
      copyCode,

      waitingForOpponent,
      joiningGame,

      icons: {
        copy: copyIcon,
      },
    };
  },
});
</script>

<style lang="scss" scoped>
</style>
