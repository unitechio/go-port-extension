const list = document.getElementById("list");
const search = document.getElementById("search");
const empty = document.getElementById("empty");
const refreshBtn = document.getElementById("refresh");

const port = chrome.runtime.connectNative("com.port.manager");

let data = [];

port.onMessage.addListener((msg) => {
  if (!msg.success) return;
  data = msg.data;
  render();
});

function render() {
  list.innerHTML = "";
  const q = search.value.toLowerCase();

  const filtered = data.filter((p) =>
    `${p.port} ${p.pid} ${p.process}`.toLowerCase().includes(q),
  );

  empty.classList.toggle("hidden", filtered.length > 0);

  filtered.forEach((p) => {
    const div = document.createElement("div");
    div.className = "card";

    div.innerHTML = `
      <div class="info">
        <div class="row1">
          <div class="port">${p.protocol} ${p.port}</div>
          <div class="status">${p.status}</div>
        </div>
        <div class="row2">
          ${p.process} • PID ${p.pid}
        </div>
      </div>
      <button class="kill">Kill</button>
    `;

    div.querySelector(".kill").onclick = () => kill(p);
    list.appendChild(div);
  });
}

function kill(p) {
  const ok = confirm(
    `Kill process?\n\n` + `${p.process}\nPID ${p.pid}\nPort ${p.port}`,
  );
  if (!ok) return;

  port.postMessage({
    action: "kill",
    pid: p.pid,
  });

  setTimeout(load, 300);
}

function load() {
  port.postMessage({ action: "list_ports" });
}

search.addEventListener("input", render);
refreshBtn.onclick = load;

load();
