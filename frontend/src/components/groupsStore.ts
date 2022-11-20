import { writable } from "svelte/store";

export interface Group {
  name: string;
  value: string;
}

export const groupsStore = writable([
  { value: "monday", name: "Montagsgruppe" },
  { value: "friday", name: "Freitag" },
  { value: "mini", name: "Minigruppe" },
]);

// TODO: fetch groups when logged in. Currently loginState is not persisted
