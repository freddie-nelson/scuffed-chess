<template>
  <!-- BOARD -->
  <div ref="boardElement" class="board flex relative" :class="{ flip }">
    <div
      class="file h-full flex flex-grow"
      :class="flip ? 'flex-col-reverse' : 'flex-col'"
      v-for="(file, c) in board"
      :key="c"
    >
      <div
        v-for="(cell, r) in file"
        :key="r"
        class="spot flex flex-grow relative justify-center items-center"
        :class="{
          'bg-primary-300': (c + r) % 2 === 0,
          'bg-primary-700': (c + r) % 2 !== 0,
          highlight: highlightedSpot.file === c && highlightedSpot.rank === r,
          draggingHome: draggingHome.file === c && draggingHome.rank === r,
          valid: board && board[c][r].valid && draggingHome.file !== -1,
        }"
        :style="{ cursor: board && board[c][r].containsPiece ? 'grab' : null }"
        @mousedown="dragStart(c, r)"
      >
        <c-game-piece
          v-if="board && board[c] && board[c][r] && board[c][r].containsPiece"
          class="piece absolute flex text-primary-500 transform scale-110"
          :piece="board[c][r].piece"
        />
      </div>
    </div>

    <div
      v-if="promotionLoc"
      class="
        absolute
        bg-bg-dark
        rounded-lg
        w-1/8
        py-4
        my-4
        transform
        scale-125
        flex flex-col
        justify-center
        items-center
      "
      :style="{
        top:
          promotionLoc.dRank === 0 && !flip
            ? '0'
            : promotionLoc.dRank === 7
            ? '0'
            : null,
        bottom:
          promotionLoc.dRank === 7 && !flip
            ? '0'
            : promotionLoc.dRank === 0
            ? '0'
            : null,
        left: `${promotionLoc.file * (100 / 8)}%`,
      }"
    >
      <c-game-piece
        class="cursor-pointer"
        :piece="{ color: $store.state.color, class: 0 }"
        @click="promotePawn(0)"
      />
      <c-game-piece
        class="mt-3 cursor-pointer"
        :piece="{ color: $store.state.color, class: 2 }"
        @click="promotePawn(2)"
      />
      <c-game-piece
        class="mt-3 cursor-pointer"
        :piece="{ color: $store.state.color, class: 3 }"
        @click="promotePawn(3)"
      />
      <c-game-piece
        class="mt-3 cursor-pointer"
        :piece="{ color: $store.state.color, class: 4 }"
        @click="promotePawn(4)"
      />
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, Ref, ref } from "vue";
import { useStore } from "@/store";
import useComponentEvent from "@/utils/useComponentEvent";

import CGamePiece from "./CGamePiece.vue";
import { Class, Color } from "@/utils/chess";

export default defineComponent({
  name: "CGameBoard",
  components: { CGamePiece },
  props: {
    flip: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const store = useStore();
    const socket = store.state.socket;

    const validMoves: Ref<{ file: number; rank: number }[]> = ref([]);

    const boardElement = ref(document.createElement("div"));
    const board = computed(() => {
      if (validMoves.value.length === 0) {
        return store.state.game?.board;
      } else {
        return store.state.game?.board.map((file) => {
          return file.map((s) => {
            return {
              ...s,
              valid:
                validMoves.value.findIndex(
                  (m) => m.file === s.file && m.rank === s.rank
                ) !== -1,
            };
          });
        });
      }
    });

    let isDragging = false;
    let draggingElement: HTMLDivElement;
    let draggingElementSpot: HTMLDivElement;
    const draggingPiece = {
      file: 0,
      rank: 0,
    };

    const highlightedSpot = ref({ file: -1, rank: -1 });
    const draggingHome = ref({ file: -1, rank: -1 });

    const dragStart = (file: number, rank: number) => {
      if (
        !board.value[file][rank].containsPiece ||
        board.value[file][rank].piece?.color !== store.state.color ||
        promotionLoc.value
      )
        return;

      socket.emit(
        "game:valid-moves",
        store.state.gameCode,
        file,
        rank,
        (json: string) => {
          validMoves.value = JSON.parse(json);
        }
      );

      isDragging = true;
      draggingPiece.file = file;
      draggingPiece.rank = rank;

      highlightedSpot.value = { file, rank };

      draggingElementSpot = boardElement.value.children[draggingPiece.file]
        .children[draggingPiece.rank] as HTMLDivElement;
      draggingElement = draggingElementSpot.querySelector(
        ".piece"
      ) as HTMLDivElement;

      draggingHome.value = {
        file,
        rank,
      };
    };

    const dragEnd = (dFile: number, dRank: number) => {
      if (!isDragging) return;

      isDragging = false;

      draggingHome.value = {
        file: -1,
        rank: -1,
      };

      draggingElement.style.zIndex = "";
      draggingElement.style.transform = "";

      highlightedSpot.value = { file: -1, rank: -1 };

      if (
        (draggingPiece.file == dFile && draggingPiece.rank == dRank) ||
        store.state.game.turn !== store.state.color ||
        validMoves.value.findIndex(
          (m) => m.file == dFile && m.rank == dRank
        ) === -1
      ) {
        validMoves.value.length = 0;
        return;
      } else {
        const lastRank = store.state.color === Color.White ? 0 : 7;
        if (
          dRank == lastRank &&
          board.value[draggingPiece.file][draggingPiece.rank].piece?.class ==
            Class.Pawn
        ) {
          promotionLoc.value = {
            file: draggingPiece.file,
            rank: draggingPiece.rank,
            dFile,
            dRank,
          };

          validMoves.value.length = 0;
          return;
        }

        makeMove(draggingPiece.file, draggingPiece.rank, dFile, dRank);
      }

      validMoves.value.length = 0;
    };

    const mouseToFileRank = (
      mouseX: number,
      mouseY: number
    ): { file: number; rank: number } => {
      if (!draggingElementSpot) return { file: -1, rank: -1 };

      const boardRect = boardElement.value.getBoundingClientRect();
      const spotRect = draggingElementSpot.getBoundingClientRect();

      let boardX = mouseX - boardRect.x;
      let boardY = mouseY - boardRect.y;

      // constrain values to inside board bounding rect
      if (boardX < 0) boardX = 0;
      else if (boardX > boardRect.width) boardX = boardRect.width - 1;

      if (boardY < 0) boardY = 0;
      else if (boardY > boardRect.height) boardY = boardRect.height - 1;

      const file = Math.floor(boardX / spotRect.width);
      const rank = Math.floor(boardY / spotRect.height);

      if (props.flip) {
        return {
          file: file,
          rank: 7 - rank,
        };
      } else {
        return {
          file,
          rank,
        };
      }
    };

    useComponentEvent(document.body, "mouseup", (event) => {
      const e = event as MouseEvent;

      const { file, rank } = mouseToFileRank(e.clientX, e.clientY);
      dragEnd(file, rank);
    });

    useComponentEvent(document.body, "mousemove", (event) => {
      const e = event as MouseEvent;

      if (isDragging) {
        const boardRect = boardElement.value.getBoundingClientRect();
        const spotRect = draggingElementSpot.getBoundingClientRect();

        const centerX = spotRect.x + spotRect.width / 2;
        const centerY = spotRect.y + spotRect.height / 2;

        let offX = e.clientX - spotRect.x - spotRect.width / 2;
        let offY = e.clientY - spotRect.y - spotRect.height / 2;

        // constrain piece to inside board
        if (offX < boardRect.x - centerX) offX = boardRect.x - centerX;
        else if (offX > boardRect.x + boardRect.width - centerX)
          offX = boardRect.x + boardRect.width - centerX;

        if (offY < boardRect.y - centerY) offY = boardRect.y - centerY;
        else if (offY > boardRect.y + boardRect.height - centerY)
          offY = boardRect.y + boardRect.height - centerY;

        draggingElement.style.zIndex = "10";
        draggingElement.style.transform = `translate(${offX}px, ${offY}px) scale(1.1)`;

        highlightedSpot.value = mouseToFileRank(e.clientX, e.clientY);
      }
    });

    const makeMove = (
      file: number,
      rank: number,
      dFile: number,
      dRank: number,
      promotion: Class = -1
    ) => {
      if (
        store.state.gameCode &&
        file >= 0 &&
        file < 8 &&
        rank >= 0 &&
        rank < 8 &&
        dFile >= 0 &&
        dFile < 8 &&
        dRank >= 0 &&
        dRank < 8
      ) {
        socket.emit(
          "game:move",
          store.state.gameCode,
          file,
          rank,
          dFile,
          dRank,
          promotion,
          (res: boolean) => {
            console.log(res);
          }
        );
      }
    };

    const promotionLoc =
      ref<{ file: number; rank: number; dFile: number; dRank: number } | null>(
        null
      );

    const promotePawn = (c: Class) => {
      if (!promotionLoc.value) return;

      makeMove(
        promotionLoc.value.file,
        promotionLoc.value.rank,
        promotionLoc.value.dFile,
        promotionLoc.value.dRank,
        c
      );

      promotionLoc.value = null;
    };

    return {
      boardElement,
      board,

      dragStart,
      dragEnd,

      highlightedSpot,
      draggingHome,

      promotionLoc,
      promotePawn,
    };
  },
});
</script>

<style lang="scss" scoped>
$board-size: clamp(380px, 50vw, 70vh);

.board {
  width: $board-size;
  height: $board-size;

  --rounded: 0.35rem;

  .file:nth-of-type(1) .spot:nth-of-type(1) {
    border-top-left-radius: var(--rounded);
  }

  .file:nth-of-type(8) .spot:nth-of-type(1) {
    border-top-right-radius: var(--rounded);
  }

  .file:nth-of-type(1) .spot:nth-of-type(8) {
    border-bottom-left-radius: var(--rounded);
  }

  .file:nth-of-type(8) .spot:nth-of-type(8) {
    border-bottom-right-radius: var(--rounded);
  }

  &.flip {
    .file:nth-of-type(1) .spot:nth-of-type(1) {
      border-radius: 0;
      border-bottom-left-radius: var(--rounded);
    }

    .file:nth-of-type(8) .spot:nth-of-type(1) {
      border-radius: 0;
      border-bottom-right-radius: var(--rounded);
    }

    .file:nth-of-type(1) .spot:nth-of-type(8) {
      border-radius: 0;
      border-top-left-radius: var(--rounded);
    }

    .file:nth-of-type(8) .spot:nth-of-type(8) {
      border-radius: 0;
      border-top-right-radius: var(--rounded);
    }
  }
}

.highlight {
  box-shadow: inset 0 0 0 4px var(--accent-100);
}

.draggingHome {
  background-color: var(--accent-400);
}

.valid {
  &::before {
    content: "";
    position: absolute;
    z-index: 5;
    background-color: var(--b-dark-dark);
    width: 35%;
    height: 35%;
    opacity: 0.5;
    border-radius: 50%;
  }
}
</style>