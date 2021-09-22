import { Color, Game } from "@/utils/chess";
import { createStore, useStore as vuexUseStore } from "vuex";

export interface Toast {
  text: string;
  duration?: number;
}

export interface Socket {
  io: SocketIOClient.Manager;
  nsp: string;
  id: string;
  connected: boolean;
  disconnected: boolean;
  binaryType: "blob" | "arraybuffer";
  open(): Socket;
  connect(): Socket;
  send(...args: any[]): Socket;
  emit(event: string, ...args: any[]): Socket;
  close(): Socket;
  disconnect(): Socket;
  compress(compress: boolean): Socket;
}

export interface Player {
  username: string;
  time: number;
}

export interface State {
  toastQueue: Toast[];
  theme: string;
  isConnected: boolean;
  color: Color;
  inGame: boolean;
  gameCode: string;
  game: Game;
  socket: Socket;
  you?: Player;
  opponent?: Player;
  ended: "" | "stalemate" | "checkmate" | "timeout" | "disconnect";
}

export default createStore<State>({
  state: {
    toastQueue: [],
    theme: "",
    isConnected: false,
    color: Color.White,
    gameCode: "",
    inGame: false,
    ended: "",

    // @ts-expect-error socket will be defined before app loads
    socket: undefined,
  },
  mutations: {
    ADD_TOAST(state, toast: Toast) {
      state.toastQueue.push(toast);
    },
    REMOVE_TOAST(state) {
      state.toastQueue.shift();
    },

    SET_THEME(state, theme: string) {
      state.theme = theme;
      localStorage.setItem("theme", theme);
    },

    SET_SOCKET(state, socket: Socket) {
      state.socket = socket;
    },

    SET_IS_CONNECTED(state, isConnected: boolean) {
      state.isConnected = isConnected;
    },

    SET_COLOR(state, color: Color) {
      state.color = color;
    },

    SET_GAME_CODE(state, code: string) {
      state.gameCode = code;
    },

    SET_IN_GAME(state, inGame: boolean) {
      state.inGame = inGame;
    },

    SET_ENDED(state, ended: "" | "stalemate" | "checkmate" | "disconnect") {
      state.ended = ended;
    },

    SET_GAME(state, game: Game) {
      state.game = game;
    },

    SET_PLAYERS(state, players: { you: Player; opponent: Player }) {
      state.you = players.you;
      state.opponent = players.opponent;
    },

    PLAY_MOVE(state, m: { file: number; rank: number; dFile: number; dRank: number }) {
      const s = state.game.board[m.file][m.rank];
      const d = state.game.board[m.dFile][m.dRank];

      const piece = s.piece;
      if (!piece) return;

      s.piece = undefined;
      s.containsPiece = false;

      d.piece = piece;
      d.containsPiece = true;
    },
  },
  actions: {},
  modules: {},
});

export const useStore = () => {
  return vuexUseStore<State>();
};
