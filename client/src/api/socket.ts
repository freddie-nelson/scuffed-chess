import store from "@/store";
import Chess from "@/utils/chess";
import io from "socket.io-client";

const SERVER_URL =
  process.env.NODE_ENV === "development"
    ? `http://${window.location.hostname}:8000`
    : "https://scuffed-chess.herokuapp.com/";

export default async function () {
  const socket = io(SERVER_URL);
  store.commit("SET_SOCKET", socket);

  // track connection status
  socket.on("connect", () => store.commit("SET_IS_CONNECTED", true));
  socket.on("disconnect", () => store.commit("SET_IS_CONNECTED", false));

  // game events
  socket.on("game:fen", (fen: string) => {
    try {
      store.commit("SET_GAME", Chess.fromFENString(fen));
    } catch (error) {
      store.commit("ADD_TOAST", { text: "Error while parsing game data.", duration: 2000 });
    }
  });

  socket.on("game:players", (json: string) => {
    try {
      store.commit("SET_PLAYERS", JSON.parse(json));
    } catch (error) {
      console.log(error);
      store.commit("ADD_TOAST", { text: "Error while parsing players data.", duration: 2000 });
    }
  });

  socket.on("game:end-state", (ended: string) => {
    try {
      store.commit("SET_ENDED", ended);
    } catch (error) {
      console.log(error);
      store.commit("ADD_TOAST", { text: "Error while parsing end state data.", duration: 2000 });
    }
  });
}
