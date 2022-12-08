<script lang="ts">
  import { Heading } from "flowbite-svelte";
  import { currentUser } from "../../components/userStore";
  import Navigation from "../../components/Navigation/Navigation.svelte";
  import UserTable from "./UserTable.svelte";
  import {getAllUser, toggleAdminUser, toggleApproveUser} from "./user.service";
  import { createNotification } from "../../components/Notification/notificationStore";

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

{#await usersPromise then users}
  <Heading class="mb-2">Alle</Heading>
  <UserTable
    {users}
    onApprove={toggleApproved}
    onAdmin={toggleAdmin}
    {disable}
  />
{/await}
