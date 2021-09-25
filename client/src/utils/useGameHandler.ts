import { Ref } from "vue";
import { Color } from "./chess";
import store from "@/store";
import router from "@/router";

export default function (username: Ref<string>, usernameError: Ref<string>) {
  const joinGame = (code: string) => {
    if (username.value.length > 18) return (usernameError.value = "Username exceeds 18 characters.");
    else if (!username.value) return (usernameError.value = "Username is empty.");

    store.state.socket.emit("game:join", username.value, code, (code: string) => {
      if (code) {
        router.push({
          name: "Game",
          query: { code: code, opponent: "true" },
        });
        store.commit("SET_IN_GAME", true);
        store.commit("SET_GAME_CODE", code);
        store.commit("SET_COLOR", Color.Black);
      } else {
        usernameError.value = "Failed to join game.";
      }
    });
  };

  const createGame = () => {
    if (username.value.length > 18) return (usernameError.value = "Username exceeds 18 characters.");
    else if (!username.value) return (usernameError.value = "Username is empty.");

    store.state.socket.emit("game:create", username.value, (code: string) => {
      if (code) {
        router.push({ name: "Game", query: { code: code } });
        store.commit("SET_GAME_CODE", code);
        store.commit("SET_IN_GAME", true);
        store.commit("SET_COLOR", Color.White);
      } else {
        usernameError.value = "Failed to create game.";
      }
    });
  };

  return {
    joinGame,
    createGame,
  };
}
