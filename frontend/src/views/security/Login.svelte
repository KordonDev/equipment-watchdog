<script lang="ts">
  import { login } from "./security.service";

  let username = "";
  const handleLogin = (e) => {
    e.preventDefault();
    login(username);
  };
</script>

<h1>Einloggen</h1>

<div class="card">
  <form on:submit={handleLogin}>
    <input bind:value={username} />
    <button type="submit">Einloggen</button>
  </form>
</div>
