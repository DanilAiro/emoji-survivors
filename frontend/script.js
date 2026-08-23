document.addEventListener("DOMContentLoaded", () => {

const AUTH_URL = "http://localhost:8081";
const GAME_WS_URL = "ws://localhost:8082/ws/connect";

const authScreen = document.getElementById("auth-screen");
const gameScreen = document.getElementById("game-screen");
const tabs = document.querySelectorAll(".tab");
const submitBtn = document.getElementById("submit-btn");
const usernameInput = document.getElementById("username");
const passwordInput = document.getElementById("password");
const authError = document.getElementById("auth-error");

let mode = "login"; // "login" | "register"

tabs.forEach(tab => {
  tab.addEventListener("click", () => {
    tabs.forEach(t => t.classList.remove("active"));
    tab.classList.add("active");
    mode = tab.dataset.mode;
    submitBtn.textContent = mode === "login" ? "Войти" : "Зарегистрироваться";
    authError.textContent = "";
  });
});

submitBtn.addEventListener("click", handleSubmit);
passwordInput.addEventListener("keydown", e => { if (e.key === "Enter") handleSubmit(); });

async function handleSubmit() {
  const username = usernameInput.value.trim();
  const password = passwordInput.value;
  authError.textContent = "";

  if (username.length < 3) {
    authError.textContent = "Имя пользователя должно быть не короче 3 символов";
    return;
  }
  if (password.length < 6) {
    authError.textContent = "Пароль должен быть не короче 6 символов";
    return;
  }

  submitBtn.disabled = true;
  submitBtn.textContent = "Подождите...";

  try {
    if (mode === "register") {
      const res = await fetch(`${AUTH_URL}/api/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password })
      });
      if (!res.ok) {
        const body = await safeJson(res);
        throw new Error(body?.error || `Ошибка регистрации (${res.status})`);
      }
      // после успешной регистрации сразу логинимся
    }

    const loginRes = await fetch(`${AUTH_URL}/api/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password })
    });
    if (!loginRes.ok) {
      const body = await safeJson(loginRes);
      throw new Error(body?.error || `Ошибка входа (${loginRes.status})`);
    }

    const { token } = await loginRes.json();
    localStorage.setItem("emoji_survivor_token", token);
    localStorage.setItem("emoji_survivor_username", username);

    startGame(token, username);
  } catch (err) {
    authError.textContent = err.message || "Что-то пошло не так";
  } finally {
    submitBtn.disabled = false;
    submitBtn.textContent = mode === "login" ? "Войти" : "Зарегистрироваться";
  }
}

async function safeJson(res) {
  try { return await res.json(); } catch { return null; }
}

// Если токен уже есть в localStorage - пробуем сразу войти в игру,
// не заставляя вводить логин/пароль заново.
(function tryAutoLogin() {
  const token = localStorage.getItem("emoji_survivor_token");
  const username = localStorage.getItem("emoji_survivor_username");
  if (token && username) {
    startGame(token, username);
  }
})();

// ==================== Игровой экран ====================

const canvas = document.getElementById("game-canvas");
const ctx = canvas.getContext("2d");
const hudUsername = document.getElementById("hud-username");
const hudKills = document.getElementById("hud-kills");
const hudPlayers = document.getElementById("hud-players");
const hpBar = document.getElementById("hp-bar");
const deathOverlay = document.getElementById("death-overlay");
const deathKills = document.getElementById("death-kills");

let socket = null;
let myUserId = null;
let latestState = { players: [], mobs: [] };

function startGame(token, username) {
  authScreen.style.display = "none";
  gameScreen.style.display = "flex";
  hudUsername.textContent = username;

  try {
    const payload = JSON.parse(atob(token.split(".")[1]));
    myUserId = payload.user_id;
  } catch {
    myUserId = null;
  }

  connectWebSocket(token);
  requestAnimationFrame(renderLoop);
}

function connectWebSocket(token) {
  socket = new WebSocket(`${GAME_WS_URL}?token=${encodeURIComponent(token)}`);

  socket.onopen = () => {
    console.log("WebSocket подключён");
  };

  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);
    if (data.type === "state") {
      latestState = data;
      updateHud();
    }
  };

  socket.onclose = () => {
    console.log("WebSocket отключён");
  };

  socket.onerror = (err) => {
    console.error("Ошибка WebSocket:", err);
  };
}

function updateHud() {
  const me = latestState.players.find(p => p.user_id === myUserId);
  hudPlayers.textContent = latestState.players.length;

  if (me) {
    hudKills.textContent = me.kills;
    const hpPercent = Math.max(0, me.hp);
    hpBar.style.width = hpPercent + "%";
    hpBar.style.background = hpPercent > 30 ? "var(--hp)" : "var(--danger)";

    if (me.hp <= 0) {
      deathOverlay.style.display = "flex";
      deathKills.textContent = me.kills;
    } else {
      deathOverlay.style.display = "none";
    }
  }
}

const pressedKeys = new Set();

const KEY_TO_DIRECTION = {
  "w": [0, -1], "arrowup": [0, -1],
  "s": [0, 1], "arrowdown": [0, 1],
  "a": [-1, 0], "arrowleft": [-1, 0],
  "d": [1, 0], "arrowright": [1, 0],
};

window.addEventListener("keydown", (e) => {
  const key = e.key.toLowerCase();

  if (key === " ") {
    sendAttack();
    e.preventDefault();
    return;
  }

  if (KEY_TO_DIRECTION[key]) {
    pressedKeys.add(key);
    sendMove();
  }
});

window.addEventListener("keyup", (e) => {
  const key = e.key.toLowerCase();
  if (KEY_TO_DIRECTION[key]) {
    pressedKeys.delete(key);
    sendMove();
  }
});

canvas.addEventListener("click", sendAttack);

function sendMove() {
  if (!socket || socket.readyState !== WebSocket.OPEN) return;

  let dx = 0, dy = 0;
  for (const key of pressedKeys) {
    const [kx, ky] = KEY_TO_DIRECTION[key];
    dx += kx;
    dy += ky;
  }

  socket.send(JSON.stringify({ action: "move", dx, dy }));
}

function sendAttack() {
  if (!socket || socket.readyState !== WebSocket.OPEN) return;
  socket.send(JSON.stringify({ action: "attack" }));
}


function renderLoop() {
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  drawGrid();

  for (const mob of latestState.mobs) {
    drawEmoji("👾", mob.x, mob.y, 24);
  }

  for (const player of latestState.players) {
    if (player.hp <= 0) continue;
    const emoji = player.user_id === myUserId ? "🙂" : "🙃";
    drawEmoji(emoji, player.x, player.y, 28);
    drawUsername(player.username, player.x, player.y);
  }

  requestAnimationFrame(renderLoop);
}

function drawGrid() {
  ctx.strokeStyle = "rgba(255,255,255,0.04)";
  ctx.lineWidth = 1;
  const step = 40;
  for (let x = 0; x < canvas.width; x += step) {
    ctx.beginPath();
    ctx.moveTo(x, 0);
    ctx.lineTo(x, canvas.height);
    ctx.stroke();
  }
  for (let y = 0; y < canvas.height; y += step) {
    ctx.beginPath();
    ctx.moveTo(0, y);
    ctx.lineTo(canvas.width, y);
    ctx.stroke();
  }
}

function drawEmoji(emoji, x, y, size) {
  ctx.font = `${size}px sans-serif`;
  ctx.textAlign = "center";
  ctx.textBaseline = "middle";
  ctx.fillText(emoji, x, y);
}

function drawUsername(name, x, y) {
  ctx.font = "11px sans-serif";
  ctx.fillStyle = "rgba(232, 230, 240, 0.8)";
  ctx.textAlign = "center";
  ctx.fillText(name, x, y - 22);
}
});