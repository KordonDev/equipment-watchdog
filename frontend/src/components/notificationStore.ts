import { writable } from "svelte/store";

export interface Notification {
  color: string;
  text: string;
  id: string;
}

export const notificationsStore = writable([]);

export const createNotification = (
  n: Omit<Notification, "id">,
  hideAfterSeconds?: number
) => {
  const id = randomId();
  notificationsStore.update((nS) => [...nS, { ...n, id }]);
  if (hideAfterSeconds) {
    setTimeout(() => removeNotification(id), hideAfterSeconds * 1000);
  }
};

export const removeNotification = (id: string) => {
  notificationsStore.update((nS) => nS.filter((n) => n.id !== id));
};

function randomId() {
  return (Math.random() + 1).toString(36).substring(7);
}
