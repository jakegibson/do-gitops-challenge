# DigitalOcean Kubernetes Challenge
## Deploy a GitOps CI/CD implementation
### Challenge Instructions
GitOps is today the way you automate deployment pipelines within Kubernetes itself, and ArgoCD  is currently one of the leading implementations. Install it to create a CI/CD solution, using tekton and kaniko for actual image building. https://medium.com/dzerolabs/using-tekton-and-argocd-to-set-up-a-kubernetes-native-build-release-pipeline-cf4f4d9972b0


### Infrastructure
1. We use [Pulumi's DigitalOcean Go SDk](https://www.pulumi.com/registry/packages/digitalocean/api-docs/) to deploy the kuberntes cluster and get the ID `ctx.Export("clusterID", cluster.ID())`
2. [Install doctl](https://docs.digitalocean.com/reference/doctl/how-to/install/)
3. Update kubeconfig/context with `doctl kubernetes cluster kubeconfig save [ID exported in step 1 (see stdout)]`

### GitOps Resources
1. We use [Pulumi's DigitalOcean Go SDk](https://www.pulumi.com/registry/packages/digitalocean/api-docs/) to deploy ArgoCD and Tekton manifests and related resources. Traefik is also installed with annotations that create a DO loadbalancer and a letsencrypt cert provisioned. The plan here was to expose the Argo and Tekton UIs via domains with Traefik but did not finish. See port forward commands below.
2. `brew install argocd` for ArgoCD CLI
3. `brew install tekton-cli` for Tekton cli `tkn`
4. `kubectl port-forward svc/argocd-server -n argocd 8080:443` expose argocd UI on localhost:8080
5. `kubectl --namespace tekton-pipelines port-forward svc/tekton-dashboard 9097:9097` expose Tekton Dashboard on localhost:9097
6. Follow getting started guides in ArgoCD and Tekton to deploy example app, tasks, and pipeline.

### Results
![Pods](/screenshots/pods.png)
![ArgoCD](/screenshots/argocd-example.png)
![TektonTasks](/screenshots/tekton-tasks.png)
![TektonPipelines](/screenshots/tekton-pipelines.png)