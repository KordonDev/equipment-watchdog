<script lang="ts">
  import {
    createNotification,
    errorNotification,
    successNotification,
  } from "../../components/Notification/notificationStore";
  import { link, replace } from "svelte-spa-router";
  import { routes } from "../../routes";
  import { login } from "./security.service";
  import { getGroups } from "../member/member.service";
  import { getGroupsWithEquipment } from "../../components/groupsStore";

  let username = "";
  const handleLogin = () => {
    login(username)
      .then(() => {
        successNotification("Login erfolgreich.");
        getGroups().then(getGroupsWithEquipment.set);
        replace(routes.MemberOverview.link);
      })
      .catch(() => {
        errorNotification("Fehler beim Login");
      });
  };
</script>

<h1>Einloggen</h1>

<div class="card">
  <form on:submit|preventDefault={handleLogin}>
    <input bind:value={username} autofocus />
    <button type="submit">Einloggen</button>
  </form>
</div>

<a href={routes.Register.link} use:link>Registrieren</a>
