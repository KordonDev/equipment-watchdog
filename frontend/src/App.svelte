<script lang="ts">
  import Router, { link, replace } from "svelte-spa-router";
  import Login from "./views/security/Login.svelte";
  import { routes } from "./routes";
  import Register from "./views/security/Register.svelte";
  import MemberOverview from "./views/member/MemberOverview.svelte";
  import AddMember from "./views/member/AddMember.svelte";
  import MemberDetail from "./views/member/MemberDetail.svelte";
  import Notification from "./components/Notification/Notification.svelte";
  import UserOverview from "./views/user/UserOverview.svelte";

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
    }}
  />
</main>

<Notification />

<style>
</style>
