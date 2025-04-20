# R2D2

Your loyal CLI droid 🤖  
Helping you monitor and manage Kubernetes deployments like a true Jedi.

---

## 🚀 Installation

### ⚡ Option 1: One-line Installation (Recommended for Linux/macOS)

> **Note:** The installation methods below are for **Linux** and **macOS**.

```bash
curl -fsSL https://raw.githubusercontent.com/pi-prakhar/r2d2/main/install.sh | bash
```

This script will automatically detect your OS (Linux/macOS) and architecture, download the latest release, and install R2D2.

### 📦 Option 2: Download from Releases

1. Go to the [Releases](https://github.com/pi-prakhar/r2d2/releases) page.
2. Download the appropriate file for your system:
   - **macOS Intel**: `r2d2-darwin-amd64.tar.gz`
   - **macOS Apple Silicon**: `r2d2-darwin-arm64.tar.gz`
   - **Linux x86_64**: `r2d2-linux-amd64.tar.gz`
   - **Windows x64**: `r2d2-windows-amd64.zip`

3. Extract the binary:
   
   **For Linux/macOS**:
   ```bash
   tar -xzf r2d2-*.tar.gz
   chmod +x r2d2
   
   # macOS only: remove quarantine attribute
   [[ "$(uname)" == "Darwin" ]] && xattr -d com.apple.quarantine ./r2d2 || true
   ```
   
   **For Windows**:
   ```powershell
   # Extract the zip file using Windows Explorer or PowerShell
   Expand-Archive -Path r2d2-windows-amd64.zip -DestinationPath .
   ```

4. Add to your PATH:
   
   **For Linux/macOS**:
   ```bash
   sudo mv ./r2d2 /usr/local/bin/
   ```
   
   **For Windows**: 
   - Move the extracted `r2d2.exe` to a permanent location (e.g., `C:\Program Files\r2d2\`)
   - Add this location to your system PATH:
     1. Right-click on "This PC" > Properties > Advanced system settings > Environment Variables
     2. Under System variables, find "Path" and click Edit
     3. Add `C:\Program Files\r2d2\` to the list
     4. Click OK to save

### 🛠 Option 3: Build from Source

Requires Go 1.24 or later.

1. Clone the repo and build:
   ```bash
   git clone https://github.com/pi-prakhar/r2d2.git
   cd r2d2
   go build -o r2d2  # On Windows, this creates r2d2.exe
   ```

2. Add to your PATH:
   
   **For Linux/macOS**:
   ```bash
   sudo mv ./r2d2 /usr/local/bin/
   
   # Set up shell completion
   source <(r2d2 completion zsh)  # For zsh users
   source <(r2d2 completion bash) # For bash users
   ```
   
   **For Windows**: Follow the same PATH setup instructions as mentioned in Option 2.

**Optional for Linux/macOS**: Use a symlink instead of copying after every build:
```bash
ln -s $(pwd)/r2d2 /usr/local/bin/r2d2
```

---

## 🛠️ Commands

```bash
r2d2 [command]
```

### Available Commands

#### Deployment Management
- `update-tag`      – Update image tag for deployments in a namespace
- `restart deployment` – Restart deployments by updating annotations
- `restart pod`     – Restart pods (by deleting them; they'll auto-restart if part of a deployment)

#### Monitoring
- `watch-images`    – Watch current container images of deployments
- `watch-tags`      – Watch image tags of deployments in real-time
  - `--pod-level, -p` – Show pod-level details instead of deployment-level
- `watch-logs`      – Watch and save logs from Kubernetes pods

#### System
- `completion`      – Generate shell autocompletion script
- `help`            – Show help for any command

### Examples

```bash
# Watch tags for deployments in a namespace
r2d2 watch-tags -n default -d nginx,redis

# Watch tags with pod-level details
r2d2 watch-tags -n default -d nginx,redis --pod-level

# Update image tag for deployments
r2d2 update-tag -n default -d nginx,redis -t v1.0.1

# Restart deployments
r2d2 restart deployment -n default -d nginx,redis

# Watch logs for pods
r2d2 watch-logs -n default -p nginx-pod-1,redis-pod-1
```

### Status Colors

R2D2 uses color-coded statuses for better visibility:

- 🟢 **Green** - Healthy/Running/Complete
- 🟠 **Orange** - Updating/Progressing
- 🟡 **Yellow** - Scaling/Starting/Pending
- 🔴 **Red** - Failed/Error/Terminated
- ⚪️ **Gray** - Unknown

You can also press tab to see the available commands. It will help you find your namespace and services within the namespace.

> 💡 Tip: Remember to source the CLI tool to make it available in your current shell session.

---

## 🤝 Contributing

We welcome contributions from the community! Whether it's bug fixes, new features, or documentation improvements, your help is appreciated.

### How to Contribute

Please see our [Contributing Guide](CONTRIBUTING.md) for detailed instructions on:

- 🍴 Forking and setting up the repository
- 🌿 Branch naming conventions (`feature/`, `fix/`, `docs/`, etc.)
- 💬 Commit message format and guidelines
- 📝 Pull request process
- 📏 Code style and standards

We follow industry-standard practices for development workflows to make collaboration smooth and efficient.

## 📄 License

This project is licensed under the [MIT License](LICENSE).


