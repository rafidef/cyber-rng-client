# CYBER UPLINK – TERMINAL CLIENT

```txt
  ______   ______  _____ ____    _   _ ____  _     ___ _   _ _  __
 / ___\ \ / / __ )| ____|  _ \  | | | |  _ \| |   |_ _| \ | | |/ /
| |    \ V /|  _ \|  _| | |_) | | | | | |_) | |    | ||  \| | ' / 
| |___  | | | |_) | |___|  _ <  | |_| |  __/| |___ | || |\  | . \ 
 \____| |_| |____/|_____|_| \_\  \___/|_|   |_____|___|_| \_|_|\_\
                      C Y B E R   U P L I N K
```

<div align="center">

### The Interface to the Shadow Network.

<p>
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Platform-Windows%20%7C%20Mac%20%7C%20Linux-333?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Client-Type-CLI%20%2F%20TUI-00ff99?style=for-the-badge" />
  <img src="https://img.shields.io/badge/Relayer-Gasless-8247E5?style=for-the-badge" />
</p>

<p>
  <img src="https://img.shields.io/github/last-commit/rafidef/cyber-rng-client?color=green&style=flat-square" />
  <img src="https://img.shields.io/github/languages/top/rafidef/cyber-rng-client?style=flat-square" />
  <img src="https://img.shields.io/github/repo-size/rafidef/cyber-rng-client?style=flat-square" />
</p>

</div>

---

## 📟 Table of Contents
- [🔓 Access Log](#-access-log)
- [👁️ Visuals](#️-visuals)
- [💾 Installation](#-installation)
- [🎮 Operations Manual](#-operations-manual)
- [📡 CLI Architecture](#-cli-architecture)
- [⚙️ Configuration](#️-configuration)

---

## 🔓 Access Log

**Cyber Uplink** is a Go-powered terminal client for the **CyberRNG Network**, providing:

- Local wallet generation  
- Offline signing & identity  
- Secure message relaying  
- TUI/HUD mining interface  

> **System Integrity:** `100%`  
> **Interface:** `CLI / TUI`  
> **Latency Mode:** `Ultra Low`  

---

## 👁️ Visuals

HUD Information:
- Wallet Address  
- $HASH Balance  
- Mining Cooldown  
- Rig Stats (GH/s)  
- Mission Tracking  

Rarity Colors:

| Rarity | Color |
|--------|--------|
| COMMON | Gray |
| UNCOMMON | Green |
| RARE | Blue |
| EPIC | Purple |
| LEGENDARY | Gold |

---

## 💾 Installation

### **Prerequisite**
Go 1.20+ installed.

### **1. Clone Uplink**
```bash
git clone https://github.com/yourusername/cyber-rng-client.git
cd cyber-rng-client
```

### **2. Install Modules**
```bash
go mod tidy
```

### **3. Run Client**
```bash
go run main.go
```

Or build binary:
```bash
go build -o uplink main.go
./uplink
```

---

## 🎮 Operations Manual

### **[1] HACK_NODE (Mining)**  
Bruteforce mining operations.

### **[2] CYBERDECK (Loadout)**  
Equip GPU / VPN.

### **[3] WORKSHOP (Overclock)**  
Enhance equipment stats.

### **[4] INVENTORY**  
Use, destroy, salvage items.

### **[5] SHADOW_NET**  
Daily Contracts & Leaderboard.

### **[6] SERVER_ROOM**  
Stake hardware for income.

---

## 📡 CLI Architecture

```mermaid
flowchart TD
    A["User Terminal / CLI"] --> B["Local Wallet Generator"]
    B --> C["session.key Identity Store"]
    A --> D["Action Selection: Mine / Equip / Enchant"]
    D --> E["Sign Payload Locally (EIP-191)"]
    E --> F["Send Signed Message to Backend"]
    F --> G["CyberRNG Core Relayer"]
    G --> H["Polygon Amoy Smart Contracts"]
    H --> I["Rewards / State Update"]
    I --> A
```

---

## ⚙️ Configuration

### **Backend Endpoint**
```go
const SERVER_URL = "http://localhost:3000"
```

### **Identity File (session.key)**

- Auto-generated on first run  
- Encrypted private key (Hex)  
- Delete file to regenerate identity  

---

<div align="center">
<img src="https://img.shields.io/badge/Uplink_Status-Operational-00ff99?style=for-the-badge" /><br><br>
<sub>Welcome to the Shadow Network.</sub>
</div>
