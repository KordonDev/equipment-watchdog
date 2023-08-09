<script lang="ts">
  import { routes } from "../../routes";
  import { link } from "svelte-spa-router";
  import { Card, Spinner } from "flowbite-svelte";
  import EquipmentIcon from "../../components/Equipment/EquipmentIcon.svelte";
  import type { Order } from "./order.service";
  import { formatToDate } from "../../components/timeHelper";
  import { getMember, type Member } from "../member/member.service";

  export let order: Order;
  export let withMember: boolean = false

  const memberPromise: Promise<Member> = withMember ?
    getMember(order.memberId.toString()) :
    Promise.resolve({} as Member);
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
  {#if withMember}
    {#await memberPromise}
      <Spinner />
    {:then member} 
      <p>Für: {member.name}</p>
    {/await}
  {/if}
</Card>
