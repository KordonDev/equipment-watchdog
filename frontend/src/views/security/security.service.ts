import { fetchApi } from "../apiService";
import { bufferDecode, bufferEncode } from "./transformer.service";

export function login(username: string) {
  return fetchApi(`/login/${username}`)
    .then((credentialRequestOptions) => {
      credentialRequestOptions.publicKey.challenge = bufferDecode(
        credentialRequestOptions.publicKey.challenge
      );
      credentialRequestOptions.publicKey.allowCredentials.forEach(
        (listItem) => (listItem.id = bufferDecode(listItem.id))
      );

      return navigator.credentials.get({
        publicKey: credentialRequestOptions.publicKey,
      });
    })
    .then((assertion) => {
      const authData = assertion.response.authenticatorData;
      const clientDataJSON = assertion.response.clientDataJSON;
      const rawId = assertion.rawId;
      const sig = assertion.response.signature;
      const userHandle = assertion.response.userHandle;

      const body = {
        id: assertion.id,
        rawId: bufferEncode(rawId),
        type: assertion.type,
        response: {
          authenticatorData: bufferEncode(authData),
          clientDataJSON: bufferEncode(clientDataJSON),
          signature: bufferEncode(sig),
          userHandle: bufferEncode(userHandle),
        },
      };

      return fetchApi(`/login/${username}`, {
        method: "POST",
        body: JSON.stringify(body),
      });
    });
}

export function register(username: string) {
  return fetchApi(`/register/${username}`)
    .then((credentialCreationOptions) => {
      credentialCreationOptions.publicKey.challenge = bufferDecode(
        credentialCreationOptions.publicKey.challenge
      );
      credentialCreationOptions.publicKey.user.id = bufferDecode(
        credentialCreationOptions.publicKey.user.id
      );
      if (credentialCreationOptions.publicKey.excludeCredentials) {
        for (
          var i = 0;
          i < credentialCreationOptions.publicKey.excludeCredentials.length;
          i++
        ) {
          credentialCreationOptions.publicKey.excludeCredentials[i].id =
            bufferDecode(
              credentialCreationOptions.publicKey.excludeCredentials[i].id
            );
        }
      }
      return navigator.credentials.create({
        publicKey: credentialCreationOptions.publicKey,
      });
    })
    .then((credential) => {
      let attestationObject = credential.response.attestationObject;
      let clientDataJSON = credential.response.clientDataJSON;
      let rawId = credential.rawId;

      const body = {
        id: credential.id,
        rawId: bufferEncode(rawId),
        type: credential.type,
        response: {
          attestationObject: bufferEncode(attestationObject),
          clientDataJSON: bufferEncode(clientDataJSON),
        },
      };

      return fetchApi(`/register/${username}`, {
        method: "POST",
        body: JSON.stringify(body),
      });
    });
}

export function logout() {
  return fetchApi(`/logout`, {
    method: "POST",
  });
}

const LOCALSTOARGE_USERNAME = "username";
export function getStoredUsername() {
  return localStorage.getItem(LOCALSTOARGE_USERNAME)
}

export function setStoredUsername(username: string) {
  localStorage.setItem(LOCALSTOARGE_USERNAME, username);
}
