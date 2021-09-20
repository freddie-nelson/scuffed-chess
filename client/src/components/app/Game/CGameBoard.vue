<template>
  <!-- BOARD -->
  <div class="board flex bg-primary-300">
    <div
      class="h-full flex flex-col flex-grow"
      v-for="(file, c) in board"
      :key="c"
    >
      <div
        v-for="(cell, r) in file"
        :key="r"
        class="flex flex-grow relative justify-center items-center"
        :class="(c + r) % 2 === 0 ? 'bg-primary-300' : 'bg-primary-700'"
      >
        <c-game-piece
          v-if="board && board[c] && board[c][r] && board[c][r].containsPiece"
          draggable
          class="
            absolute
            flex
            cursor-pointer
            text-primary-500
            transform
            scale-110
          "
          :piece="board[c][r].piece"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent } from "vue";
import { useStore } from "@/store";
import CGamePiece from "./CGamePiece.vue";

export default defineComponent({
  name: "CGameBoard",
  components: { CGamePiece },
  setup() {
    const store = useStore();

    const board = computed(() => {
      return store.state.game?.board;
    });

    return {
      board,
    };
  },
});
</script>

<style lang="scss" scoped>
$board-size: clamp(380px, 50vw, 70vh);

.board {
  width: $board-size;
  height: $board-size;
}
</style>