<script lang="ts">
  import { Heading } from "flowbite-svelte";
  import { currentUser } from "../../components/userStore";
  import Navigation from "../../components/Navigation/Navigation.svelte";
  import UserTable from "./UserTable.svelte";
  import { getAllUser, toggleApproveUser } from "./user.service";
  import { createNotification } from "../../components/Notification/notificationStore";

  function toogleApproved(id: number) {
    toggleApproveUser(id)
      .then((user) => {
        createNotification(
          {
            color: "green",
            text: `Toggled ${user.name} approve Status.`,
          },
          5
        );
      })
      .catch(() => {
        createNotification(
          {
            color: "red",
            text: `Fehler beim updaten.`,
          },
          5
        );
      });
  }
  function toogleAdmin(id: number) {
    toggleApproveUser(id)
      .then((user) => {
        createNotification(
          {
            color: "green",
            text: `Toggled ${user.name} admin Status.`,
          },
          5
        );
      })
      .catch(() => {
        createNotification(
          {
            color: "red",
            text: `Fehler beim updaten.`,
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
      disable = u.isAdmin;
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
  onApprove={toogleApproved}
  onAdmin={toogleAdmin}
  {disable}
/>

{#await usersPromise then users}
  <Heading class="mb-2">Alle</Heading>
  <UserTable
    {users}
    onApprove={toogleApproved}
    onAdmin={toogleAdmin}
    {disable}
  />
{/await}
