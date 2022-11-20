<script lang="ts">
  import { Alert } from "flowbite-svelte";
  import { notificationsStore, removeNotification } from "./notificationStore";
  import type { Notification } from "src/components/notificationStore";

  let notifications: Notification[];

  notificationsStore.subscribe((value) => {
    notifications = value;
  });
</script>

<div class="alert-container">
  {#each notifications as notification}
    <Alert
      class="mb-4 top-alert"
      color={notification.color}
      dismissable
      style="display: block;"
      on:close={() => removeNotification(notification.id)}
    >
      {notification.text}
    </Alert>
  {/each}
</div>

<style>
  .alert-container {
    position: absolute;
    top: 20px;
    right: 20px;
  }
</style>
