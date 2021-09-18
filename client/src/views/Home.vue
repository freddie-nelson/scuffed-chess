<template>
  <main class="w-full h-full flex flex-col justify-center items-center">
    <c-gradient-heading noscale :size="7">Scuffed Chess</c-gradient-heading>

    <div class="mt-8 max-w-md w-full">
      <c-input-text
        name="gameCode"
        placeholder="Game Code"
        label="Game Code"
        :error="joinCodeError"
        maxlength="6"
        v-model="joinCode"
      />
      <c-button class="w-full mt-3 h-14" @click="tryShowJoinModal">
        Join Game
      </c-button>

      <c-button-outline
        class="w-full mt-6 h-16"
        @click="showCreateModal = true"
      >
        Create Game
      </c-button-outline>
    </div>

    <c-modal
      v-if="showJoinModal"
      closeable
      @close="
        showJoinModal = false;
        usernameError = '';
      "
      class="max-w-lg w-full flex flex-col items-center justify-center"
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
      />
      <c-button class="w-full mt-3 h-14" @click="joinGame">
        Join Game
      </c-button>
    </c-modal>

    <c-modal
      v-if="showCreateModal"
      closeable
      @close="
        showCreateModal = false;
        usernameError = '';
      "
      class="max-w-lg w-full flex flex-col items-center justify-center"
    >
      <h1 class="font-bold text-3xl mb-4">Create Game</h1>

      <c-input-text
        classes="w-full"
        name="username"
        placeholder="Username"
        label="Username"
        :error="usernameError"
        maxlength="18"
        v-model="username"
        dark
      />
      <c-button class="w-full mt-3 h-14" @click="createGame">
        Create Game
      </c-button>
    </c-modal>
  </main>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";

import CButton from "../components/shared/Button/CButton.vue";
import CButtonOutline from "../components/shared/Button/CButtonOutline.vue";
import CGradientHeading from "../components/shared/Heading/CGradientHeading.vue";
import CInputText from "../components/shared/Input/CInputText.vue";
import CModal from "../components/shared/Modal/CModal.vue";

export default defineComponent({
  name: "Home",
  components: {
    CGradientHeading,
    CInputText,
    CButton,
    CButtonOutline,
    CModal,
  },
  setup() {
    const username = ref("");
    const usernameError = ref("");

    // join
    const joinCode = ref("");
    const joinCodeError = ref("");

    const showJoinModal = ref(false);
    const tryShowJoinModal = () => {
      usernameError.value = "";

      if (joinCode.value.length === 6) {
        showJoinModal.value = true;
        joinCodeError.value = "";
      } else {
        joinCodeError.value = "Invalid game code.";
      }
    };

    const joinGame = () => {
      if (username.value.length > 18)
        return (usernameError.value = "Username exceeds 18 characters.");
      else if (!username.value)
        return (usernameError.value = "Username is empty.");
    };

    // create
    const showCreateModal = ref(false);

    const createGame = () => {
      if (username.value.length > 18)
        return (usernameError.value = "Username exceeds 18 characters.");
      else if (!username.value)
        return (usernameError.value = "Username is empty.");
    };

    return {
      username,
      usernameError,

      joinCode,
      joinCodeError,
      showJoinModal,
      tryShowJoinModal,

      joinGame,

      showCreateModal,

      createGame,
    };
  },
});
</script>