import { EquipmentType } from "./views/equipment/equipment.service";

export const routes = {
  Login: {
    path: "/login",
    link: "/login",
  },
  Register: {
    path: "/register",
    link: "/register",
  },
  MemberOverview: {
    path: "/member/overview",
    link: "/member/overview",
  },
  MemberDetail: {
    path: "/member/detail/:id",
    link: "/member/detail/",
  },
  AddMember: {
    path: "/member/add",
    link: "/member/add",
  },
  Users: {
    path: "/users/",
    link: "/users",
  },
  EquipmentAdd: {
    path: "/equipment/add",
    link: "/equipment/add",
  },
  EquipmentType: {
    path: "/equipment/type/:type",
    link: `/equipment/type/`,
  },
  EquipmentDetails: {
    path: "/equipment/details/:id",
    link: `/equipment/details/`,
  },
  NotApproved: {
    path: "/not-approved",
    link: "/not-approved",
  },
};
