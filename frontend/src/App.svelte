<script lang="ts">
  import Router, { replace } from "svelte-spa-router";
  import Login from "./views/security/Login.svelte";
  import { routes } from "./routes";
  import Register from "./views/security/Register.svelte";
  import MemberOverview from "./views/member/MemberOverview.svelte";
  import AddMember from "./views/member/AddMember.svelte";
  import MemberDetail from "./views/member/MemberDetail.svelte";
  import Notification from "./components/Notification/Notification.svelte";
  import UserOverview from "./views/user/UserOverview.svelte";
  import NotApproved from "./views/security/NotApproved.svelte";
  import CreateEquipment from "./views/equipment/NewEquipment.svelte";
  import EquipmentByType from "./views/equipment/EquipmentByType.svelte";
  import EquipmentDetails from "./views/equipment/EquipmentDetails.svelte";
  import Orders from "./views/order/Orders.svelte";
  import OrderDetails from "./views/order/OrderDetails.svelte";
  import AddOrder from "./views/order/NewOrder.svelte";
  import FulfilledOrders from "./views/order/FulfilledOrders.svelte";

  if (window.location.pathname === "/" && window.location.hash === "") {
    replace(routes.MemberOverview.link);
  }
</script>

<main>
  <Router
    routes={{
      [routes.Register.path]: Register,
      [routes.Login.path]: Login,
      [routes.MemberOverview.path]: MemberOverview,
      [routes.AddMember.path]: AddMember,
      [routes.MemberDetail.path]: MemberDetail,
      [routes.Users.path]: UserOverview,
      [routes.EquipmentType.path]: EquipmentByType,
      [routes.EquipmentAdd.path]: CreateEquipment,
      [routes.EquipmentDetails.path]: EquipmentDetails,
      [routes.Orders.path]: Orders,
      [routes.OrdersFulfilled.path]: FulfilledOrders,
      [routes.OrderDetails.path]: OrderDetails,
      [routes.AddOrder.path]: AddOrder,
      [routes.NotApproved.path]: NotApproved,
    }}
  />
</main>

<Notification />

<style>
</style>
