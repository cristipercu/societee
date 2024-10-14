export function makeid(length) {
  let result = '';
  const characters = 'abcdefghijklmnopqrstuvwxyz0123456789';
  const charactersLength = characters.length;
  let counter = 0;
  while (counter < length) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
    counter += 1;
  }
  return result;
}

export function getInternalApiHeaders(withLogin) {
  if (withLogin) {
    return { 
      'Content-Type': 'application/json',
      'Authorization': localStorage.getItem('token')
    }
  } else {
    return {
      'Content-Type': 'application/json'
    }
  }
}

export async function handleLogin(email, password, message, url) {
      const response = await fetch(url , {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });

      const data = await response.json();   
      if (response.ok) {
        localStorage.setItem('token', data.token);
        localStorage.setItem('displayName', data.user);
        window.location.href = '/'; 
      } else {
        message = data.error;
      }
  }

export function connectToWebsocket(url, onOpen, onMessage, onClose, onError) {
  const ws = new WebSocket(url);

  ws.onopen = () => {
    if (onOpen) {
      onOpen(ws);
    }
  };

  ws.onmessage = (event) => {
    if (onMessage) {
      onMessage(event.data);
    } else {
      console.log(event.data);
    }
  };

  ws.onclose = () => {
    if (onClose) {
      onClose();
    } else {
      ws.close()
    }

  };

  ws.onerror = (error) => {
    console.error('ws error:', error);
    if (onError) {
      onError(error);
    }
  };

  return ws;
}

export function sendMessage(ws, user, type, message) {
  if (ws.readyState === WebSocket.OPEN) {
    const messageData = {
      user: user,
      type: type,
      message: message
    };
    ws.send(JSON.stringify(messageData));
  } else {
    console.error('no ws connection');
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

export function clearGameLogicFromLocalStorage() {
  localStorage.removeItem('roomCode');
  localStorage.removeItem('roomID');
  localStorage.removeItem('roomAdmin');
}

export function AmITheAdmin() {
  return localStorage.getItem('displayName') === localStorage.getItem('roomAdmin')
}
