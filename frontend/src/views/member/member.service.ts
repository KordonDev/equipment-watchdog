import { fetchApi } from "../apiService";

interface Member {
  id: string;
  name: string;
}

export function getMembers(): Promise<Member[]> {
  return fetchApi("/members/");
}
