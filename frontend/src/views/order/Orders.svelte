<script lang="ts">
  import { routes } from "../../routes";
  import { link, push } from "svelte-spa-router";
  import Navigation from "../../components/Navigation/Navigation.svelte";
  import { Alert, Card, Label, Select, Spinner } from "flowbite-svelte";
  import EquipmentIcon from "../../components/Equipment/EquipmentIcon.svelte";
  import { getOrders } from "./order.service";

  let orderPromise = getOrders();
</script>

<Navigation />
<h1>Bestellte Ausrüstung</h1>

<a href={routes.EquipmentAdd.link} use:link>Ausrüstung bestellen</a>

{#await orderPromise}
  <Spinner />
{:then orders}
  <div>
    Aktuell {orders.length} Bestellungen offen
    <div class="flex flex-wrap">
      {#each orders as order}
        <div>
          <Card class="m-4">
            <EquipmentIcon
              equipmentType={order.type}
              registrationCode={undefined}
            />
            <h3>
              <a href={`${routes.OrderDetails.link}${order.id}`} use:link>
                {order.id}
              </a>
            </h3>
            <p>Größe: {order.size || "-"}</p>
            <p>Bestelldatum: {order.createdAt}</p>
          </Card>
        </div>
      {/each}
    </div>
  </div>
{:catch}
  <Alert class="mb-4 top-alert" color="red" dismissable style="display: block;">
    Leider konnten die Bestellungen nicht geladen werden.
  </Alert>
{/await}
