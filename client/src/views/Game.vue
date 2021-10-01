<template>
  <main class="w-full h-full flex flex-col justify-center items-center">
    <c-modal v-if="waitingForOpponent" class="flex flex-col text-bg-light">
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

    <c-modal
      v-if="showJoinModal"
      class="
        max-w-lg
        w-full
        flex flex-col
        items-center
        justify-center
        text-bg-light
      "
    >
      <h1 class="font-bold text-3xl mb-4">Join Game</h1>

      <c-input-text
        classes="w-full"
        name="username"
        placeholder="Username"
        label="Username"
        :error="usernameError"
        maxlength="18"
        v-model="username"
        dark
        @keyup.enter="
          () => {
            res = joinGame(code);
            !res ? (showJoinModal = false) : null;
          }
        "
      />
      <c-button
        class="w-full mt-3 h-14"
        @click="
          () => {
            res = joinGame(code);
            !res ? (showJoinModal = false) : null;
          }
        "
      >
        Join Game
      </c-button>
    </c-modal>

    <c-modal v-if="joiningGame" class="flex flex-col text-bg-light">
      <h1 class="font-bold text-3xl">Joining game...</h1>
    </c-modal>

    <c-modal
      v-if="$store.state.ended"
      class="
        flex flex-col
        items-center
        justify-center
        max-w-sm
        w-full
        text-bg-light
      "
    >
      <h1 class="font-bold text-3xl">
        {{
          $store.state.ended[0].toUpperCase() + $store.state.ended.substring(1)
        }}
      </h1>
      <h2
        v-if="$store.state.ended === 'checkmate'"
        class="mt-1 font-bold text-xl"
      >
        {{ $store.state.game.turn === 0 ? "Black" : "White" }} has won.
      </h2>

      <c-button class="mt-4 h-16 w-full" @click="resetGame">
        Leave Game
      </c-button>
    </c-modal>

    <!-- OPPONENT CARD -->
    <c-game-player
      v-if="$store.state.opponent"
      :player="$store.state.opponent"
    />

    <c-game-board :flip="flip" />

    <!-- PLAYER CARD -->
    <c-game-player v-if="$store.state.you" :player="$store.state.you" />
  </main>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  onBeforeMount,
  onBeforeUnmount,
  ref,
  watch,
} from "vue";
import { useStore } from "@/store";
import { useRouter } from "vue-router";

import useComponentEvent from "@/utils/useComponentEvent";
import useGameHandler from "@/utils/useGameHandler";

import CModal from "@/components/shared/Modal/CModal.vue";
import CInputText from "@/components/shared/Input/CInputText.vue";
import CButton from "@/components/shared/Button/CButton.vue";
import CGameBoard from "@/components/app/Game/CGameBoard.vue";
import CGamePlayer from "@/components/app/Game/CGamePlayer.vue";

import { Icon } from "@iconify/vue";
import copyIcon from "@iconify-icons/feather/copy";

export default defineComponent({
  name: "Game",
  components: {
    CModal,
    CInputText,
    CButton,
    Icon,
    CGameBoard,
    CGamePlayer,
  },
  setup() {
    const store = useStore();
    const router = useRouter();

    const code = router.currentRoute.value.query.code?.toString() as string;
    let isOpponent = !!router.currentRoute.value.query.opponent;

    const showJoinModal = ref(false);
    const username = ref("");
    const usernameError = ref("");

    const { joinGame } = useGameHandler(username, usernameError);

    watch(usernameError, (err) => {
      if (err && err.startsWith("Failed")) {
        store.commit("ADD_TOAST", { text: err, duration: 2000 });
        resetGame();
      }
    });

    onBeforeMount(() => {
      if (!store.state.inGame) {
        if (code.length === 6) {
          isOpponent = true;
          showJoinModal.value = true;
        } else {
          router.push({ name: "Home" });
        }
      }
    });

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

    const resetGame = () => {
      store.state.socket.emit("game:leave", store.state.gameCode, isOpponent);

      store.commit("SET_GAME", undefined);
      store.commit("SET_IN_GAME", false);
      store.commit("SET_GAME_CODE", "");
      store.commit("SET_ENDED", "");
      store.commit("SET_PLAYERS", { you: undefined, opponent: undefined });

      router.push({ name: "Home" });
    };

    onBeforeUnmount(() => {
      resetGame();
    });

    useComponentEvent(
      window as any as HTMLElement,
      "beforeunload" as keyof HTMLElementEventMap,
      () => {
        resetGame();
      }
    );

    const flip = computed(() => store.state.color === 0);

    return {
      showJoinModal,
      username,
      usernameError,
      joinGame,

      code,
      copyCode,

      waitingForOpponent,
      joiningGame,

      resetGame,

      flip,

      icons: {
        copy: copyIcon,
      },
    };
  },
});
</script>

<style lang="scss" scoped>
</style>
