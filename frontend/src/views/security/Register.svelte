<script lang="ts">
  import { routes } from "../../routes";
  import {link, replace} from "svelte-spa-router";
  import { register } from "./security.service";
  import { createNotification } from "../../components/Notification/notificationStore";

  let username = "";
  const login = (e: SubmitEvent) => {
    e.preventDefault();
    register(username)
      .then((u) => {
        createNotification(
          {
            color: "green",
            text: `Registrierung '${u.name}' erfolgreich. Jetzt direkt einloggen.`,
          },
          5
        );
        replace(routes.Login.link)
      })
      .catch((err) => {
        createNotification({
          color: "red",
          text: `Fehler bei der Registrierung. ${err}`,
        });
      });
  };
</script>

<h1>Registrieren</h1>
<form on:submit={login}>
  <input bind:value={username} />
  <button type="submit">Register</button>
</form>

<a href={routes.Login.link} use:link>Einloggen</a>
