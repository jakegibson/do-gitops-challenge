# DigitalOcean Kubernetes Challenge
## Deploy a GitOps CI/CD implementation
### Challenge Instructions
GitOps is today the way you automate deployment pipelines within Kubernetes itself, and ArgoCD  is currently one of the leading implementations. Install it to create a CI/CD solution, using tekton and kaniko for actual image building. https://medium.com/dzerolabs/using-tekton-and-argocd-to-set-up-a-kubernetes-native-build-release-pipeline-cf4f4d9972b0


### Infrastructure
1. We use [Pulumi's DigitalOcean Go SDk](https://www.pulumi.com/registry/packages/digitalocean/api-docs/) to deploy the kuberntes cluster and get the ID `ctx.Export("clusterID", cluster.ID())`
2. [Install doctl](https://docs.digitalocean.com/reference/doctl/how-to/install/)
3. Update kubeconfig/context with `doctl kubernetes cluster kubeconfig save [ID exported in step 1 (see stdout)]`

### GitOps Resources
1. We use [Pulumi's DigitalOcean Go SDk](https://www.pulumi.com/registry/packages/digitalocean/api-docs/) to deploy ArgoCD and Tekton manifests and related resources
2. `brew install argocd` for ArgoCD CLI