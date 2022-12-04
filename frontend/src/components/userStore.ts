import { getMe, type User } from "../views/user/user.service";
import { readable } from "svelte/store";

export const currentUser = readable<User>(undefined, function start(set) {
  getMe().then((user) => set(user));
});
