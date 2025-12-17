Rất tốt 👍
Dưới đây là **README.md viết theo chuẩn OSS chuyên nghiệp**, tone **engineering / production / maintainable**, giống phong cách repo của **HashiCorp / Docker / Kubernetes tool**.

Bạn có thể **copy 100% dùng ngay** để public GitHub.

---

```md
# Go Port Manager

A **production-ready Chrome Extension + Go Native Messaging Host**  
for inspecting local network ports and managing running processes securely.

> Built for developers who need **visibility and control over local ports**
> without running insecure local servers.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Security Model](#security-model)
- [Project Structure](#project-structure)
- [Requirements](#requirements)
- [Installation](#installation)
  - [Build Native Host](#build-native-host)
  - [Install Native Host](#install-native-host)
  - [Register Native Messaging Host (Windows)](#register-native-messaging-host-windows)
  - [Load Chrome Extension](#load-chrome-extension)
- [Usage](#usage)
- [Troubleshooting](#troubleshooting)
- [Production Considerations](#production-considerations)
- [FAQ](#faq)
- [License](#license)

---

## Overview

**Go Port Manager** is a local developer utility that provides:

- Visibility into **active TCP/UDP ports**
- Associated **PID and process information**
- Ability to **terminate processes safely**

Unlike traditional solutions, this project:

- ❌ Does **not** expose a local HTTP server
- ❌ Does **not** require users to manually start background services
- ✅ Uses **Chrome Native Messaging**, the only secure and approved way
  for browser extensions to interact with the operating system

---

## Features

- List active ports with protocol, PID, process name, and state
- Real-time search (port / PID / process)
- Kill processes directly from the extension UI
- Secure communication via stdin/stdout
- Chrome & Edge compatible
- Dark, developer-friendly UI
- No external dependencies at runtime

---

## Architecture

Browser extensions are sandboxed and **cannot access OS processes**.

This project uses **Chrome Native Messaging**:

```

┌────────────────────────┐
│ Chrome Extension (UI)  │
└──────────▲─────────────┘
│ JSON messages
│ (stdin / stdout)
┌──────────┴─────────────┐
│ Go Native Host Binary  │
└──────────▲─────────────┘
│ OS syscalls
┌──────────┴─────────────┐
│ Operating System       │
└────────────────────────┘

```

Chrome:
- Spawns the native process on demand
- Restricts access to whitelisted extensions
- Terminates the process automatically

---

## Security Model

This project follows Chrome’s **official security model**:

- No open network ports
- No IPC over localhost
- No remote access
- Only explicitly whitelisted extensions can connect
- Native binary path must be absolute and user-controlled

This approach is **Chrome Web Store compliant**.

---

## Project Structure

```

go-port-manager/
├── extension/                 # Chrome Extension UI
│   ├── manifest.json
│   ├── popup.html
│   ├── popup.css
│   └── popup.js
│
├── native-host/               # Go Native Messaging Host
│   ├── main.go
│   ├── go.mod
│   └── go.sum
│
└── README.md

````

---

## Requirements

- Go ≥ 1.20
- Google Chrome or Microsoft Edge
- Windows / macOS / Linux
- Administrator privileges (required for killing processes)

---

## Installation

### Build Native Host

```bash
cd native-host
go mod tidy
````

#### Windows

```bash
go build -o port-manager.exe
```

#### macOS / Linux

```bash
go build -o port-manager
```

---

### Install Native Host

Create an installation directory.

#### Windows (recommended)

```
C:\Program Files\PortManager\
```

Copy:

* `port-manager.exe`
* `com.port.manager.json`

---

### Native Messaging Host Manifest

Create:

```
C:\Program Files\PortManager\com.port.manager.json
```

```json
{
  "name": "com.port.manager",
  "description": "Go Port Manager Native Host",
  "path": "C:\\Program Files\\PortManager\\port-manager.exe",
  "type": "stdio",
  "allowed_origins": [
    "chrome-extension://<EXTENSION_ID>/"
  ]
}
```

> **Note**
>
> * The path must be absolute
> * The extension ID must match exactly
> * A trailing slash is required

---

### Register Native Messaging Host (Windows)

Open `regedit` and navigate to:

```
HKEY_CURRENT_USER
└─ Software
   └─ Google
      └─ Chrome
         └─ NativeMessagingHosts
```

Create a key:

```
com.port.manager
```

Set its default value to:

```
C:\Program Files\PortManager\com.port.manager.json
```

---

### Load Chrome Extension

1. Open `chrome://extensions`
2. Enable **Developer Mode**
3. Click **Load unpacked**
4. Select the `extension/` directory
5. Copy the generated **Extension ID**
6. Update `allowed_origins`
7. Restart Chrome

---

## Usage

Open the extension popup to:

* View active ports
* Filter results in real time
* Terminate processes with confirmation

Chrome manages the lifecycle of the native host automatically.

---

## Troubleshooting

### Common Errors

| Error                                          | Cause                             |
| ---------------------------------------------- | --------------------------------- |
| `Specified native messaging host not found`    | Registry key missing or incorrect |
| `Access to native messaging host is forbidden` | `allowed_origins` mismatch        |
| `Disconnected port object`                     | Native host rejected connection   |

### Quick Diagnostics

```cmd
C:\Program Files\PortManager\port-manager.exe
```

If the binary runs without exiting immediately, the host is valid.

---

## Production Considerations

* Package native host via installer (`.msi`, `.pkg`)
* Sign binaries (Windows/macOS)
* Remove wildcard origins before release
* Add logging to stderr for diagnostics
* Implement graceful termination (SIGTERM)
* Consider feature-gating destructive actions

---

## FAQ

### Why not use a local HTTP server?

* Requires manual startup
* Exposes attack surface
* Not Chrome Store friendly

### Can this be published to Chrome Web Store?

Yes. Native Messaging is officially supported.

### Does this require admin privileges?

Yes, for process termination.

---

## License

MIT License

---

## Maintainers

This project is intended to serve as a **reference implementation**
for secure local tooling via Chrome Native Messaging.

Contributions and improvements are welcome.

```
