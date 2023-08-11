<script lang="ts">
  import Navigation from "../../components/Navigation/Navigation.svelte";
  import { link } from "svelte-spa-router";
  import { Alert, Spinner } from "flowbite-svelte";
  import { getOrders } from "./order.service";
  import OrderCard from "./OrderCard.svelte";
  import { routes } from "../../routes";

  let orderPromise = getOrders();
</script>

<Navigation />
<h1>Bestellte Ausrüstung</h1>
<p><a href={routes.OrdersFulfilled.link} use:link>Fertige Bestellungen</a></p>

{#await orderPromise}
  <Spinner />
{:then orders}
  <div>
    Aktuell {orders.length} Bestellungen offen
    <div class="flex flex-wrap">
      {#each orders as order}
        <OrderCard {order} withMember={true} />
      {/each}
    </div>
  </div>
{:catch}
  <Alert class="mb-4 top-alert" color="red" dismissable style="display: block;">
    Leider konnten die Bestellungen nicht geladen werden.
  </Alert>
{/await}
