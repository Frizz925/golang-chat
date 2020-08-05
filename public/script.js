const protocol = window.location.protocol;
const scheme = protocol.substring(0, protocol.length - 1);
const wsScheme = scheme === "https" ? "wss" : "ws";

const host = window.location.host;
const wsUrl = `${wsScheme}://${host}/ws/messages`;
const apiUrl = `${scheme}://${host}/api/messages`;

const setupWebSocket = () => {
  const status = document.getElementById("status");
  const output = document.getElementById("output");

  const socket = new WebSocket(wsUrl);

  socket.addEventListener("open", () => {
    status.innerHTML = "Connected";
  });

  socket.addEventListener("close", () => {
    status.innerHTML = "Disconnected";
  });

  socket.addEventListener("message", (event) => {
    output.innerText += `\n${event.data}`;
  });

  status.innerHTML = "Connecting...";
};

const setupInputForm = () => {
  const form = document.getElementById("form");
  const input = document.getElementById("input");

  form.addEventListener("submit", (event) => {
    event.preventDefault();
    const xhr = new XMLHttpRequest();
    xhr.onload = () => {
      input.value = "";
    };
    xhr.open("POST", apiUrl);
    xhr.send(input.value);
  });
};

const setupOutput = () => {
  const output = document.getElementById("output");
  const xhr = new XMLHttpRequest();
  xhr.open("GET", apiUrl);
  xhr.onload = () => {
    const res = JSON.parse(xhr.responseText);
    output.innerText = res.result.join("\n");
  };
  xhr.send();
};

window.onload = () => {
  setupWebSocket();
  setupInputForm();
  setupOutput();
};
