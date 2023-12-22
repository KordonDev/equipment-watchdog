import { fetchApi } from "../apiService";
let chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_';

    // Use a lookup table to find the index.
    let lookup = new Uint8Array(256);
    for (let i = 0; i < chars.length; i++) {
        lookup[chars.charCodeAt(i)] = i;
    }

function bufferDecode(base64string) {
  let bufferLength = base64string.length * 0.75,
        len = base64string.length, i, p = 0,
        encoded1, encoded2, encoded3, encoded4;

        let bytes = new Uint8Array(bufferLength);

        for (i = 0; i < len; i+=4) {
            encoded1 = lookup[base64string.charCodeAt(i)];
            encoded2 = lookup[base64string.charCodeAt(i+1)];
            encoded3 = lookup[base64string.charCodeAt(i+2)];
            encoded4 = lookup[base64string.charCodeAt(i+3)];

            bytes[p++] = (encoded1 << 2) | (encoded2 >> 4);
            bytes[p++] = ((encoded2 & 15) << 4) | (encoded3 >> 2);
            bytes[p++] = ((encoded3 & 3) << 6) | (encoded4 & 63);
        }

        return bytes.buffer
}

function bufferEncode(arraybuffer) {
  let bytes = new Uint8Array(arraybuffer),
        i, len = bytes.length, base64url = '';

        for (i = 0; i < len; i+=3) {
            base64url += chars[bytes[i] >> 2];
            base64url += chars[((bytes[i] & 3) << 4) | (bytes[i + 1] >> 4)];
            base64url += chars[((bytes[i + 1] & 15) << 2) | (bytes[i + 2] >> 6)];
            base64url += chars[bytes[i + 2] & 63];
        }

        if ((len % 3) === 2) {
            base64url = base64url.substring(0, base64url.length - 1);
        } else if (len % 3 === 1) {
            base64url = base64url.substring(0, base64url.length - 2);
        }

        return base64url;
}

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
