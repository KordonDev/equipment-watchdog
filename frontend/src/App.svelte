<script lang="ts">
  import Counter from "./lib/Counter.svelte";
  import Hello from "./lib/Hello.svelte";
  import Router from "svelte-spa-router";
</script>

<main>
  <nav>
    <a href="/">Home</a>
    <a href="/#/about">About</a>
    <a href="/#/blog/colors">Colors</a>
  </nav>
  <div class="card">
    <Counter />
  </div>

  <Router
    routes={{
      "/": Counter,
      "/about": Hello,
    }}
  />
</main>

<style>
</style>
