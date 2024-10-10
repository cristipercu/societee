export function makeid(length) {
  let result = '';
  const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  const charactersLength = characters.length;
  let counter = 0;
  while (counter < length) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
    counter += 1;
  }
  return result;
}

export async function handleLogin(email, password, message, url) {
    try {
      const response = await fetch(url , {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });

      if (response.ok) {
        const data = await response.json();   

        localStorage.setItem('token', data.token);
        localStorage.setItem('displayName', data.user);
        window.location.href = '/';
      } else {
        const data = await response.json();
        message = data.error;
      }
    } catch (error) {
      message = 'error connecting to server, call CP!';
    }
  }


export function parseJwt (token) {
    var base64Url = token.split('.')[1];
    var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    var jsonPayload = decodeURIComponent(window.atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
}

export function isTokenExpired (token) {
  if (!token) {
    return true;
  }
  try {
    let expDate = parseJwt(token)['expiresAt'];
    const currentTime = Math.floor(Date.now() / 1000);
    return expDate < currentTime;
  } catch (error) {
    return true
  }
}

