# R2D2

Your loyal CLI droid 🤖  
Helping you monitor and manage Kubernetes deployments like a true Jedi.

---

## 🚀 Installation (macOS)

### 📦 Option 1: Download from Releases (Recommended)

1. Go to the [Releases](https://github.com/yourusername/r2d2/releases) page.
2. Download the latest `r2d2-darwin-amd64.tar.gz` or `r2d2-darwin-arm64.tar.gz` based on your system.
3. Extract and prepare the binary:
   ```bash
   tar -xzf r2d2-darwin-*.tar.gz
   chmod +x r2d2
   xattr -d com.apple.quarantine ./r2d2
   ```
4. Move it to your PATH:
   ```bash
   sudo mv ./r2d2 /usr/local/bin/
   ```

Now you can run `r2d2` from anywhere 🎉

### 🛠 Option 2: Build from Source

1. Clone the repo and build:
   ```bash
   git clone https://github.com/yourusername/r2d2.git
   cd r2d2
   go build -o r2d2
   ```
2. Move it to your PATH:
   ```bash
   sudo mv ./r2d2 /usr/local/bin/
   source <(r2d2 completion zsh)
   ```

(Optional) Use a symlink instead of copying after every build:
```bash
ln -s $(pwd)/r2d2 /usr/local/bin/r2d2
```

---

## 🛠️ Commands

```bash
r2d2 [command]
```

### Available Commands

- `update-tag`   – Update image tag for deployments in a namespace.
- `watch-images` – Watch current container images of deployments.
- `watch-tags`   – Watch image tags of deployments in real-time.

- `completion`   – Generate shell autocompletion script.
- `help`         – Show help for any command.

You can also press tab to see the available commands. It will also help you find your namespace and the services within the namespace. 
Tip: Remember to source the cli tool to make it available in your current shell session.

> 💡 More commands coming soon to enhance your Kubernetes workflow!


