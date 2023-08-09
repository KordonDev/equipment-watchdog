<script lang="ts">
  import Navigation from "../../components/Navigation/Navigation.svelte";
  import { Alert, Spinner } from "flowbite-svelte";
  import { getOrders } from "./order.service";
    import OrderCard from "./OrderCard.svelte";

  let orderPromise = getOrders();
</script>

<Navigation />
<h1>Bestellte Ausrüstung</h1>

{#await orderPromise}
  <Spinner />
{:then orders}
  <div>
    Aktuell {orders.length} Bestellungen offen
    <div class="flex flex-wrap">
      {#each orders as order}
        <OrderCard order={order} withMember={true} />
      {/each}
    </div>
  </div>
{:catch}
  <Alert class="mb-4 top-alert" color="red" dismissable style="display: block;">
    Leider konnten die Bestellungen nicht geladen werden.
  </Alert>
{/await}
