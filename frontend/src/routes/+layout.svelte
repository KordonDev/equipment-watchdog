<script>
  import "../app.css";
  import favicon from "$lib/assets/favicon.svg";
  import Notification from "$lib/components/Notification.svelte";

  let { children } = $props();
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>

<Notification />
{@render children?.()}
