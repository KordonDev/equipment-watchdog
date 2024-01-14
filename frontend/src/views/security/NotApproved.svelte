<script lang="ts">
  import { Button } from "flowbite-svelte";
  import { logout } from "../security/security.service";

</script>

<h1 class="mb-12">Nicht bestätigt</h1>
<p class="mb-6">
  Dein Account wurde noch nicht bestätigt, melde dich bei einem Administrator,
  damit dieser deinen Account bestätigen kann.
</p>
<Button size="xs" color="red" class="ml-5" on:click={logout}>Ausloggen</Button>
