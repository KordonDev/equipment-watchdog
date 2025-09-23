<script lang="ts">
	import { onDestroy } from 'svelte';
  import { page } from '$app/state';

  const links = [
    { href: "/members", label: "Mitglieder" },
    { href: "/equipments", label: "Ausrüstung" },
    { href: "/orders", label: "Bestellungen" },
    { href: "/changes", label: "Änderungen" },
		{ href: "/users", label: "Benutzer" }
  ];
  let showMenu = false;
  let menuRef: any;
  let buttonRef: any;


  function handleDocumentClick(event: any) {
    if (
      showMenu &&
      menuRef &&
      buttonRef &&
      !menuRef.contains(event.target) &&
      !buttonRef.contains(event.target)
    ) {
      showMenu = false;
    }
  }

  $: if (showMenu) {
    document.addEventListener("click", handleDocumentClick);
  } else {
    document.removeEventListener("click", handleDocumentClick);
  }

  onDestroy(() => {
    document.removeEventListener("click", handleDocumentClick);
  });
</script>

<div class="burger-menu-container">
  <button class="burger-menu-button" bind:this={buttonRef} on:click={() => showMenu = !showMenu} aria-label="Menu">
    &#9776;
  </button>
  <nav class="burger-menu" bind:this={menuRef} class:open={showMenu}>
    <ul>
      {#each links as link}
        <li>
          <a href={link.href} class:selected={link.href === page.url.pathname}>{link.label}</a>
        </li>
      {/each}
    </ul>
  </nav>
</div>

<style>
.burger-menu-container {
  position: fixed;
  top: 1rem;
  right: 1rem;
  z-index: 1000;
  display: flex;
  flex-direction: row-reverse;
  align-items: flex-start;
}
.burger-menu-button {
  font-size: 2rem;
  background: none;
  border: none;
  cursor: pointer;
}
.burger-menu {
  position: relative;
  min-width: 180px;
  background: #fff;
  border: 1px solid #ccc;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  transform: translateX(110%);
  transition: transform 0.3s cubic-bezier(.4,0,.2,1);
  opacity: 0;
  pointer-events: none;
}
.burger-menu.open {
  transform: translateX(0);
  opacity: 1;
  pointer-events: auto;
}
.burger-menu ul {
  list-style: none;
  margin: 0;
  padding: 0.5rem 1rem;
}
.burger-menu li {
  margin: 0.5rem 0;
}
.burger-menu a {
  color: #333;
	text-decoration: underline;
}
.burger-menu a.selected {
  font-weight: bold;
  color: #1976d2;
  background: #e3f2fd;
  border-radius: 4px;
  padding: 2px 6px;
}
</style>
