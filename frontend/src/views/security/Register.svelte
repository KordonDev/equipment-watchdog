<script lang="ts">
  import { register } from "./security.service";

  let username = "";
  const login = (e: SubmitEvent) => {
    e.preventDefault();
    register(username);
  };
</script>

<h1>Registrieren</h1>
<form on:submit={login}>
  <input bind:value={username} />
  <button type="submit">Register</button>
</form>
