<script lang="ts">
  import { routes } from "../../routes";
  import {link, replace} from "svelte-spa-router";
  import { register } from "./security.service";
  import { createNotification } from "../../components/Notification/notificationStore";
  import { Label, Input, Button } from "flowbite-svelte";

  let username = "";
  let loading = false;
  const onRegister = (e: SubmitEvent) => {
    e.preventDefault();
    loading = true;
    register(username)
      .then((u) => {
        loading = false;
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
        loading = false;
        createNotification({
          color: "red",
          text: `Fehler bei der Registrierung. ${err}`,
        });
      });
  };
</script>

<h1>Registrieren</h1>
<div class="card">
    <form on:submit={onRegister}>
        <Label for="usernname" class="block mb-2">Benutzer:</Label>
        <Input required class="mb-4" id="username" bind:value={username} />
        <Button color="purple" type="submit" disabled={loading}>
          Registrieren
        </Button>
    </form>
</div>

<a href={routes.Login.link} use:link>Zum Einloggen</a>
