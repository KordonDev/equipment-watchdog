<script lang="ts">
  import { Heading } from "flowbite-svelte";
  import { currentUser } from "../../components/userStore";
  import Navigation from "../../components/Navigation/Navigation.svelte";
  import UserTable from "./UserTable.svelte";
  import { getAllUser } from "./user.service";

  function toogleApproved(id: number) {
    console.log(id);
  }
  function toogleAdmin(id: number) {
    console.log("Admin", id);
  }

  let me = [];
  currentUser.subscribe((u) => {
    if (u) {
      me = [u];
    } else {
      me = [];
    }
  });

  let usersPromise = getAllUser();
</script>

<Navigation />
<Heading>User</Heading>

<UserTable users={me} onApprove={toogleApproved} onAdmin={toogleAdmin} />

{#await usersPromise then users}
  <UserTable {users} onApprove={toogleApproved} onAdmin={toogleAdmin} />
{/await}
