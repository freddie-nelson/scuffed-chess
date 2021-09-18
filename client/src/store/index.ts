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

export interface State {
  toastQueue: Toast[];
  isConnected: boolean;
  socket: Socket;
}

export default createStore<State>({
  state: {
    toastQueue: [],
    isConnected: false,

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

    SET_SOCKET(state, socket: Socket) {
      state.socket = socket;
    },

    SET_IS_CONNECTED(state, isConnected: boolean) {
      state.isConnected = isConnected;
    },
  },
  actions: {},
  modules: {},
});

export const useStore = () => {
  return vuexUseStore<State>();
};
