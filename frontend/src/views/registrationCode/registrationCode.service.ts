import { fetchApi } from "../apiService";
import { getDate } from "../../components/timeHelper";

export interface RegistrationCode<T = Date> {
  id: string;
  reservedUntil: T;
}

let cachedRegistrationCode: RegistrationCode;

export function useRegistrationCode(code: string) {
  if (code === cachedRegistrationCode.id) {
    cachedRegistrationCode = undefined;
  }
}

export function getRegistrationCode(): Promise<RegistrationCode> {
  if (cachedRegistrationCode && cachedRegistrationCode.reservedUntil.getTime() > new Date().getTime()) {
    return Promise.resolve(cachedRegistrationCode);
  }
  return getNewRegistrationCode();
}

export function getNewRegistrationCode(): Promise<RegistrationCode> {
  return fetchApi(`/registrationCode/`)
    .then(rc => {
      const rcD = {
        ...rc,
        reservedUntil: new Date(rc.reservedUntil)
      };
      cachedRegistrationCode = rcD;
      return rcD;
    });
}

