<template>
  <div class="player-card h-16 py-3 flex items-center text-t-main">
    <p class="text-lg font-semibold font-mono">{{ player.username }}</p>
  </div>
</template>

<script lang="ts">
import { Player } from "@/store";
import { defineComponent } from "vue";

export default defineComponent({
  name: "CGamePlayer",
  components: {},
  props: {
    player: {
      type: Object as () => Player,
      required: true,
    },
  },
});
</script>

<style lang="scss" scoped>
$board-size: clamp(380px, 50vw, 70vh);

.player-card {
  width: $board-size;
}
</style>