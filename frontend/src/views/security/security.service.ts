import { BASE_URL } from "../../constants";

function bufferDecode(value) {
  return Uint8Array.from(atob(value), (c) => c.charCodeAt(0));
}

function bufferEncode(value) {
  return btoa(String.fromCharCode.apply(null, new Uint8Array(value)))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=/g, "");
}

export function login(username: string) {
  fetch(`${BASE_URL}/login/${username}`, {
    credentials: "include",
  })
    .then((res) => res.json())
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

      return fetch(`${BASE_URL}/login/${username}`, {
        method: "POST",
        body: JSON.stringify(body),
        credentials: "include",
      });
    })
    .then(() => {
      alert("Success");
    })
    .catch(console.error);
}

export function register(username: string) {
  fetch(`${BASE_URL}/register/${username}`, {
    credentials: "include",
  })
    .then((res) => res.json())
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

      return fetch(`${BASE_URL}/register/${username}`, {
        method: "POST",
        body: JSON.stringify(body),
        credentials: "include",
      });
    })
    .then((succ) => alert(`registered ${username} ${succ}`))
    .catch((err) => console.error(err));
}
