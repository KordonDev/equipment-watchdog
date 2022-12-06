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

export function toggleApproveUser(userId: number): Promise<User> {
  return fetchApi(`/users/${userId}/toogle-approve`, {
    method: "PATCH",
  });
}

export function toggleAdminUser(userId: number): Promise<User> {
  return fetchApi(`/users/${userId}/toogle-admin`, {
    method: "PATCH",
  });
}