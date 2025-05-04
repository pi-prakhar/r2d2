## 🤖 Auto-Deploy Guide

The `auto-deploy` command automates the process of deploying to Kubernetes when GitHub workflows complete successfully. It monitors workflow status in real-time and automatically updates your deployments.

### Prerequisites

1. **GitHub Token**

   ```bash
   # Set your GitHub token as an environment variable
   export GITHUB_TOKEN="your_github_token"
   ```

   > 💡 The token needs `repo` scope to access workflow information

2. **Kubernetes Access**
   - Make sure you have access to your Kubernetes cluster
   - The correct context is selected (`kubectl config current-context`)
   - Proper permissions to update deployments

### Usage

```bash
r2d2 auto-deploy [flags]
```

#### Required Flags

- `-r, --repo` - GitHub repository name
- `-t, --tag` - Git tag to watch
- `-n, --namespace` - Kubernetes namespace
- `-d, --names` - Comma-separated list of deployment names

#### Optional Flags

- `-i, --interval` - Polling interval in seconds (default: 10)

### Examples

1. **Basic Usage**

   ```bash
   r2d2 auto-deploy \
     -r my-service \
     -t v1.0.0 \
     -n production \
     -d frontend,backend
   ```

2. **Custom Polling Interval**
   ```bash
   r2d2 auto-deploy \
     -r my-service \
     -t v1.0.0 \
     -n production \
     -d frontend,backend \
     -f 30
   ```

### Workflow

1. **Start Monitoring**

   - Command starts watching GitHub workflows for the specified tag
   - Shows real-time status with modern UI

2. **Final Summary Display**

   ```
    🚀 [Tag] v-tag
    📦 [Repository] Owner/repo
    🏁 Final Summary:

    ✅ [ECR Push] - Success
    ✅ [Deployment] - All deployments updated successfully!

    🔗 GitHub Actions: https://github.com/Owner/repo/actions

   ```

3. **Automatic Deployment**
   - When all workflows complete successfully:
     - Updates specified Kubernetes deployments
     - Provides final success/failure summary
