import { fetchApi } from "../apiService";

export interface User {
  id: number;
  name: string;
  isApproved: boolean;
  isAdmin: boolean;
}

export function getMe(): Promise<User> {
  return fetchApi(`/me`);
}

export function getAllUser(): Promise<User[]> {
  return fetchApi(`/users/`);
}

export function toggleApproveUser(username: string): Promise<User> {
  return fetchApi(`/users/${username}/toggle-approve`, {
    method: "PATCH",
  });
}

export function toggleAdminUser(username: string): Promise<User> {
  return fetchApi(`/users/${username}/toggle-admin`, {
    method: "PATCH",
  });
}
