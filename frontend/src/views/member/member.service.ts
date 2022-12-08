import { fetchApi } from "../apiService";

export interface Member {
  id: string;
  name: string;
  group: string;
}

export function getMembers(): Promise<Member[]> {
  return fetchApi("/members/");
}

export function getMember(id: string): Promise<Member> {
  return fetchApi(`/members/${id}`);
}

export function createMember(member: Member): Promise<Member> {
  return fetchApi("/members/", {
    method: "POST",
    body: JSON.stringify({
      ...member,
      id: 0,
    }),
  });
}

export function updateMember(member: Member): Promise<Member> {
  return fetchApi(`/members/${member.id}`, {
    method: "PUT",
    body: JSON.stringify(member),
  });
}

export function deleteMember(id: string): Promise<void> {
  return fetchApi(`/members/${id}`, {
    method: "DELETE",
  });
}

export enum Group {
  FRIDAY = "friday",
  MONDAY = "monday",
  MINI = "mini"
}

export enum Equipment {
  Helmet = "helmet",
  Jacket = "jacket",
  Gloves = "gloves",
  Trousers = "trousers",
  Boots = "boots",
  TShirt = "tshirt"
}

export interface Groups {
  [Group.FRIDAY]: Equipment[];
  [Group.MONDAY]: Equipment[];
  [Group.MINI]: Equipment[];
}
export function getGroups(): Promise<Groups> {
  return fetchApi("/members/groups");
}
