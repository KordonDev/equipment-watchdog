<script lang="ts">
  import {Button, Checkbox, Heading, Input, Label} from "flowbite-svelte";
  import { currentUser } from "../../components/userStore";
  import Navigation from "../../components/Navigation/Navigation.svelte";
  import UserTable from "./UserTable.svelte";
  import {getAllUser, toggleAdminUser, toggleApproveUser} from "./user.service";
  import {
    createNotification,
    errorNotification,
    successNotification
  } from "../../components/Notification/notificationStore";
  import {changePassword} from "../security/security.service";

  function toggleApproved(username: string) {
    toggleApproveUser(username)
      .then((user) => {
        createNotification(
          {
            color: "green",
            text: `${user.name} ${user.isApproved? 'bestätigt': 'Bestätigung zurückgenommen'}.`,
          },
          5
        );
      })
      .catch(() => {
        createNotification(
          {
            color: "red",
            text: `Fehler beim Speichern.`,
          },
          5
        );
      });
  }
  function toggleAdmin(username: string) {
    toggleAdminUser(username)
      .then((user) => {
        createNotification(
          {
            color: "green",
            text: `${user.name} hat Adminrechte ${user.isAdmin ? 'bekommen' : 'verloren'}.`,
          },
          5
        );
      })
      .catch(() => {
        createNotification(
          {
            color: "red",
            text: `Fehler beim Speichern.`,
          },
          5
        );
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
</script>

<Navigation />
<Heading class="mb-4">User</Heading>

<Heading class="mb-2">Du</Heading>
<UserTable
  users={me}
  onApprove={toggleApproved}
  onAdmin={toggleAdmin}
  {disable}
/>

<Heading class="mb-2 mt-6">Password</Heading>
<form on:submit|preventDefault={handleChangePassword}>
  <Label for="password" class="block mb-2">Password:</Label>
  <Input required class="mb-4" id="password" type="password" bind:value={password} />
  <Button color="purple" type="submit" disabled={loading}>Passwort speichern</Button>
</form>

{#await usersPromise then users}
  <Heading class="mb-2 mt-6">Alle</Heading>
  <UserTable
    {users}
    onApprove={toggleApproved}
    onAdmin={toggleAdmin}
    {disable}
  />
{/await}
