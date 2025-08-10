<script lang="ts">
  import {Button, Heading, Input, Label} from "flowbite-svelte";
  import { currentUser } from "../../components/userStore";
  import Navigation from "../../components/Navigation/Navigation.svelte";
  import UserTable from "./UserTable.svelte";
  import {getAllUser, toggleAdminUser, toggleApproveUser} from "./user.service";
  import {
    errorNotification,
    successNotification
  } from "../../components/Notification/notificationStore";
  import {addLogin, changePassword} from "../security/security.service";

  function toggleApproved(username: string) {
    toggleApproveUser(username)
      .then((user) => {
        successNotification(`${user.name} ${user.isApproved? 'bestätigt': 'Bestätigung zurückgenommen'}.`)
      })
      .catch(() => {
        errorNotification(`Fehler beim Speichern.`)
      });
  }
  function toggleAdmin(username: string) {
    toggleAdminUser(username)
      .then((user) => {
        successNotification(`${user.name} hat Adminrechte ${user.isAdmin ? 'bekommen' : 'verloren'}.`)
      })
      .catch(() => {
        errorNotification(`Fehler beim Speichern.`)
      });
  }

  let me = [];
  let disable = true;
  currentUser.subscribe((u) => {
    if (u) {
      me = [u];
      disable = !u.isAdmin;
    } else {
      me = [];
    }
  });

  let usersPromise = getAllUser();

  let password = "";
  let loading = false;

  const handleChangePassword = () => {
    loading = true;
    changePassword(password)
      .then(() => {
        loading = false;
        successNotification("Passwort gespeichert.");
      })
      .catch(() => {
        loading = false;
        errorNotification(`Fehler beim Speichern.`);
      });
  }

  const addLoginMethod = () => {
    addLogin()
      .then(() => {
        successNotification(`Loginmöglichkeit hinzugefügt.`)
      })
      .catch((err) => {
        errorNotification(`Fehler beim Hinzufügen. ${err}`)
      });
  }
</script>

<Navigation />

{#if me && me.length > 0 && me[0].isAdmin}
  <Heading class="mb-2">Du</Heading>
  <UserTable
    users={me}
    onApprove={toggleApproved}
    onAdmin={toggleAdmin}
    {disable}
  />
{/if}

<Heading class="mb-2 mt-6">Password</Heading>
<form on:submit|preventDefault={handleChangePassword}>
  <Label for="password" class="block mb-2">Password:</Label>
  <Input required class="mb-4" id="password" type="password" bind:value={password} />
  <Button color="purple" type="submit" disabled={loading}>Passwort speichern</Button>
</form>

<Heading class="mb-2 mt-6">Loginmöglichkeit hinzufügen</Heading>
<Button color="purple" on:click={addLoginMethod}>
  Hinzufügen
</Button>

{#if me && me.length > 0 && me[0].isAdmin}
  {#await usersPromise then users}
    <Heading class="mb-2 mt-6">Alle</Heading>
    <UserTable
      {users}
      onApprove={toggleApproved}
      onAdmin={toggleAdmin}
      {disable}
    />
  {/await}
{/if}
