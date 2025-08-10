<script lang="ts">
  import {
    errorNotification,
    successNotification,
  } from "../../components/Notification/notificationStore";
  import { link, replace } from "svelte-spa-router";
  import { routes } from "../../routes";
  import {
    login,
    getStoredUsername,
    setStoredUsername,
  } from "./security.service";
  import { getGroups } from "../member/member.service";
  import {getGroupsWithEquipment, getTranslatedGroups, loadTranslatedGroups} from "../../components/groupsStore";
  import { Label, Input, Button, Checkbox } from "flowbite-svelte";

  let username = getStoredUsername();
  let storeUsername = false;
  setTimeout(() => {
    storeUsername = username && username !== "";
  }, 0);
  let loading = false;
  const handleLogin = () => {
    loading = true;
    login(username)
      .then(() => {
        loading = false;
        successNotification("Login erfolgreich.");
        getGroups().then(getGroupsWithEquipment.set);
        loadTranslatedGroups().then(getTranslatedGroups.set);
        if (storeUsername) {
          setStoredUsername(username);
        } else {
          setStoredUsername("");
        }
        replace(routes.MemberOverview.link);
      })
      .catch(() => {
        errorNotification("Fehler beim Login");
        loading = false;
      });
  };
</script>

<h1>Einloggen</h1>

<div class="card">
  <form on:submit|preventDefault={handleLogin}>
    <Label for="username" class="block mb-2">Benutzer:</Label>
    <Input required class="mb-4" id="username" bind:value={username} />
    <Checkbox bind:checked={storeUsername} class="mb-4">
      Benutzername speichern
    </Checkbox>
    <Button color="purple" type="submit" disabled={loading}>Einloggen</Button>
  </form>
</div>

<div>
  <a href={routes.PasswordLogin.link} use:link>Zum Passwordlogin</a>
</div>
<div>
  <a href={routes.Register.link} use:link>Zur Registrierung</a>
</div>
