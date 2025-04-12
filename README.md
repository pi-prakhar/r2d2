# R2D2

```
          ___
       ,-'___'-.
     ,'  [(_)]  '.
    |_]||[][O]o[][|
  * |*____________| _
 | []   _______   [] |
 | []   _______   [] |
[| ||      _      || |]
 |_|| =   [=]     ||_|
 | || =   [|]     || |
 | ||      _      || |
 | ||||   (+)    (|| |
 | ||_____________|| |
 |_| \___________/ |_|
 / \      | |      / \
/___\    /___\    /___\
```

> Your loyal CLI droid 🤖  
> Helping you monitor and manage Kubernetes deployments like a true Jedi.

---

## 📦 Features

For available commands and options, run:
```
r2d2 --help
```
    
---

## 🚀 Installation

### 1. Clone the repository and build
```
git clone https://github.com/yourusername/r2d2.git
cd r2d2
go build -o r2d2
```

### 2. Move the binary to a directory in your PATH
```
sudo mv r2d2 /usr/local/bin/
```

Now you can run `r2d2` from anywhere 🎉

### (Optional) Use a symlink to avoid copying after every build
```
ln -s $(pwd)/r2d2 /usr/local/bin/r2d2
```

---

## 🧪 Usage

### `watch-tags`
Watch and display deployment image tags from the Kubernetes cluster in real-time.
```
r2d2 watch-tags --namespace <namespace> --services <service1,service2>
```

This will show a live-updating table with the image tags of the specified deployments.

### (More commands coming soon!)
Future commands will be added to enhance your workflow with Kubernetes.
