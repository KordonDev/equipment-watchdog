<script lang="ts">
  import { createNotification } from "../../components/Notification/notificationStore";
  import { link, replace } from "svelte-spa-router";
  import { routes } from "../../routes";
  import { login } from "./security.service";

  let username = "";
  const handleLogin = () => {
    login(username)
      .then((u) => {
        createNotification(
          {
            color: "green",
            text: `Login erfolgreich.`,
          },
          5
        );
        replace(routes.MemberOverview.link);
      })
      .catch((err) => {
        createNotification({
          color: "red",
          text: `Fehler beim Login`,
        });
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
