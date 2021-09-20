<template>
  <div :class="{ white: piece.color === 1, black: piece.color === 0 }">
    <!-- QUEEN -->
    <svg
      v-if="piece.class === 0"
      xmlns="http://www.w3.org/2000/svg"
      viewbox="0 0 45 45"
      preserveAspectRatio="xMinYMin meet"
    >
      <g class="piece-svg" stroke-linejoin="round">
        <path
          stroke-linecap="round"
          stroke-width="1.5"
          d="M8 12a2 2 0 1 1-4 0 2 2 0 1 1 4 0zM24.5 7.5a2 2 0 1 1-4 0 2 2 0 1 1 4 0zM41 12a2 2 0 1 1-4 0 2 2 0 1 1 4 0zM16 8.5a2 2 0 1 1-4 0 2 2 0 1 1 4 0zM33 9a2 2 0 1 1-4 0 2 2 0 1 1 4 0z"
        />
        <path
          fill-rule="evenodd"
          stroke-width="1.5"
          d="M9 26c8.5-1.5 21-1.5 27 0l2-12-7 11V11l-5.5 13.5-3-15-3 15-5.5-14V25L7 14l2 12z"
        />
        <path
          fill-rule="evenodd"
          stroke-width="1.5"
          d="M9 26c0 2 1.5 2 2.5 4 1 1.5 1 1 .5 3.5-1.5 1-1.5 2.5-1.5 2.5-1.5 1.5.5 2.5.5 2.5 6.5 1 16.5 1 23 0 0 0 1.5-1 0-2.5 0 0 .5-1.5-1-2.5-.5-2.5-.5-2 .5-3.5 1-2 2.5-2 2.5-4-8.5-1.5-18.5-1.5-27 0z"
        />
        <path
          fill-rule="evenodd"
          stroke-linecap="round"
          d="M11.5 30c3.5-1 18.5-1 22 0M12 33.5c6-1 15-1 21 0"
        />
      </g>
    </svg>

    <!-- KING -->
    <svg
      v-else-if="piece.class === 1"
      xmlns="http://www.w3.org/2000/svg"
      viewbox="0 0 45 45"
      preserveAspectRatio="xMinYMin meet"
    >
      <g class="piece-svg" fill-rule="evenodd" stroke-width="1.5">
        <path stroke-linecap="round" d="M22.5 11.625V6" />
        <path
          d="M22.5 25s4.5-7.5 3-10.5c0 0-1-2.5-3-2.5s-3 2.5-3 2.5c-1.5 3 3 10.5 3 10.5"
        />
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M11.5 37c5.5 3.5 15.5 3.5 21 0v-7s9-4.5 6-10.5c-4-6.5-13.5-3.5-16 4V27v-3.5c-3.5-7.5-13-10.5-16-4-3 6 5 10 5 10V37z"
        />
        <path
          stroke-linecap="round"
          d="M20 8h5M11.5 29.5C17 27 27 27 32.5 30"
        />
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M11.5 37c5.5-2.5 15.5-2.5 21 0M11.5 33.5c5.5-2 15.5-2 21 0"
        />
      </g>
    </svg>

    <!-- ROOK -->
    <svg
      v-else-if="piece.class === 2"
      xmlns="http://www.w3.org/2000/svg"
      viewbox="0 0 45 45"
      preserveAspectRatio="xMinYMin meet"
      version="1.0"
    >
      <g class="piece-svg" fill-rule="evenodd">
        <path
          stroke-linejoin="round"
          stroke-width="1.5"
          d="M9 39h27v-3H9v3zM12 36v-4h21v4H12zM11 14V9h4v2h5V9h5v2h5V9h4v5"
        />
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="1.5"
          d="m34 14-3 3H14l-3-3"
        />
        <path stroke-width="1.5" d="M31 17v12.5H14V17" />
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="1.5"
          d="m31 29.5 1.5 2.5h-20l1.5-2.5"
        />
        <path stroke-linecap="round" d="M11 14h23" />
      </g>
    </svg>

    <!-- BISHOP -->
    <svg
      v-else-if="piece.class === 3"
      xmlns="http://www.w3.org/2000/svg"
      viewbox="0 0 45 45"
      preserveAspectRatio="xMinYMin meet"
      version="1.0"
    >
      <g class="piece-svg" stroke-width="1.5">
        <path
          fill-rule="evenodd"
          stroke-linejoin="round"
          d="M9 36c3.385-.972 10.115.43 13.5-2 3.385 2.43 10.115 1.028 13.5 2 0 0 1.646.542 3 2-.677.972-1.646.986-3 .5-3.385-.972-10.115.458-13.5-1-3.385 1.458-10.115.028-13.5 1-1.354.486-2.323.472-3-.5 1.354-1.945 3-2 3-2z"
        />
        <path
          fill-rule="evenodd"
          stroke-linejoin="round"
          d="M15 32c2.5 2.5 12.5 2.5 15 0 .5-1.5 0-2 0-2 0-2.5-2.5-4-2.5-4 5.5-1.5 6-11.5-5-15.5-11 4-10.5 14-5 15.5 0 0-2.5 1.5-2.5 4 0 0-.5.5 0 2z"
        />
        <path
          stroke-linejoin="round"
          d="M25 8a2.5 2.5 0 1 1-5 0 2.5 2.5 0 1 1 5 0z"
        />
        <path d="M17.5 26h10M15 30h15" />
        <path stroke-linecap="round" d="M22.5 15.5v5M20 18h5" />
      </g>
    </svg>

    <!-- KNIGHT -->
    <svg
      v-else-if="piece.class === 4"
      xmlns="http://www.w3.org/2000/svg"
      viewbox="0 0 45 45"
      preserveAspectRatio="xMinYMin meet"
      version="1.0"
    >
      <g class="piece-svg" stroke-linecap="round">
        <g fill-rule="evenodd" stroke-width="1.5">
          <path d="M22 10c10.5 1 16.5 8 16 29H15c0-9 10-6.5 8-21" />
          <path
            stroke-linejoin="round"
            d="M24 18c.384 2.911-5.553 7.369-8 9-3 2-2.82 4.343-5 4-1.042-.944 1.413-3.038 0-3-1 0 .187 1.232-1 2-1 0-4.003 1-4-4 0-2 6-12 6-12s1.886-1.902 2-3.5c-.726-.994-.5-2-.5-3 1-1 3 2.5 3 2.5h2s.782-1.992 2.5-3c1 0 1 3 1 3"
          />
        </g>
        <path
          stroke-linejoin="round"
          stroke-width="1.5"
          d="M9.5 25.5a.5.5 0 1 1-1 0 .5.5 0 1 1 1 0zM14.933 15.75a.5 1.5 30 1 1-.866-.5.5 1.5 30 1 1 .866.5z"
        />
        <path fill="none" d="M37 39c1-20-5.5-27.5-12-28.5" />
      </g>
    </svg>

    <!-- PAWN -->
    <svg
      v-else-if="piece.class === 5"
      xmlns="http://www.w3.org/2000/svg"
      viewbox="0 0 45 45"
      preserveAspectRatio="xMinYMin meet"
      version="1.0"
    >
      <path
        class="piece-svg"
        fill="white"
        stroke="black"
        stroke-dashoffset="10"
        stroke-linecap="round"
        stroke-width="1.5"
        d="M22 9a4.001 4.001 0 0 0-3.219 6.375C16.83 16.497 15.5 18.588 15.5 21a6.47 6.47 0 0 0 2.406 5.031C14.907 27.09 10.5 31.578 10.5 39.5h23c0-7.922-4.407-12.41-7.406-13.469A6.47 6.47 0 0 0 28.5 21c0-2.412-1.33-4.503-3.281-5.625A4.001 4.001 0 0 0 22 9z"
      />
    </svg>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

import { Piece } from "@/utils/chess";

export default defineComponent({
  name: "CGamePiece",
  components: {},
  props: {
    piece: {
      type: Object as () => Piece,
      required: true,
    },
  },
});
</script>

<style lang="scss" scoped>
svg {
  width: 45px;
  height: 45px;
}

.white {
  --stroke: var(--primary-900);
  --fill: var(--primary-100);
}

.black {
  --stroke: var(--accent-900);
  --fill: var(--accent-300);
}

.piece-svg {
  stroke: var(--stroke, #000);
  fill: var(--fill, #fff);
}
</style>