import { BASE_URL } from "../../constants";

interface Member {
  id: string;
  name: string;
}

export function getMembers(): Promise<Member[]> {
  return fetch(`${BASE_URL}/members/`, {
    credentials: "include",
  }).then((res) => res.json());
}
