<template>
  <main class="w-full h-full flex flex-col justify-center items-center">
    <!-- OPPONENT CARD -->
    <div></div>

    <!-- BOARD -->
    <div class="board flex flex-col">
      <div class="w-full flex flex-grow" v-for="r in 8" :key="r">
        <div
          v-for="c in 8"
          :key="c"
          class="flex flex-grow"
          :class="(c + r) % 2 === 0 ? 'bg-primary-300' : 'bg-primary-700'"
        ></div>
      </div>
    </div>

    <!-- PLAYER CARD -->
    <div></div>
  </main>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "Game",
  components: {},
});
</script>

<style lang="scss" scoped>
$board-size: clamp(350px, 50vw, 70vh);

.board {
  width: $board-size;
  height: $board-size;
}
</style>
