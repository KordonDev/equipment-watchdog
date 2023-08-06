<script lang="ts">
  import { routes } from "../../routes";
  import { link } from "svelte-spa-router";
  import { Card } from "flowbite-svelte";
  import EquipmentIcon from "../../components/Equipment/EquipmentIcon.svelte";
  import type { Order } from "./order.service";
  import { formatToDate } from "../../components/timeHelper";

  export let order: Order;
</script>

<Card class="m-4 flex items-center">
  <EquipmentIcon equipmentType={order.type} registrationCode={undefined} />
  <h3>
    <a href={`${routes.OrderDetails.link}${order.id}`} use:link>
      {order.id}
    </a>
  </h3>
  <p>Größe: {order.size || "-"}</p>
  <p>Bestelldatum: {formatToDate(order.createdAt)}</p>
</Card>
