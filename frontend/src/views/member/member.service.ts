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
  debugger;
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

interface Groups {
  name: string[];
}
export function getGroups(): Promise<Groups> {
  return fetchApi("/members/groups");
}
