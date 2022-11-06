<script lang="ts">
  import { routes } from "../../routes";
  import { login } from "./security.service";

  let username = "";
  const handleLogin = (e) => {
    e.preventDefault();
    login(username).then(() =>
      window.location.replace(`/#${routes.MemberOverview}`)
    );
  };
</script>

<h1>Einloggen</h1>

<div class="card">
  <form on:submit={handleLogin}>
    <input bind:value={username} autofocus />
    <button type="submit">Einloggen</button>
  </form>
</div>
