package main

import (
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var domain = "do-challenge.fulldeploy.dev"

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cert, err := createCert(ctx)
		if err != nil {
			return err
		}

		err = deployTraefik(ctx, cert.Uuid)
		if err != nil {
			return err
		}
		// err = deployLinkered(ctx)
		// if err != nil {
		// 	return err
		// }
		err = deployArgoCD(ctx)
		if err != nil {
			return err
		}
		// err = deployTekton(ctx)
		// if err != nil {
		// 	return err
		// }

		ctx.Export("certificateName", cert.Name)
		ctx.Export("certificateID", cert.ID())
		ctx.Export("certificateUUID", cert.Uuid)
		ctx.Export("certificateUrn", cert.URN())
		return nil
	})
}

// Create cert for domain that will be pointed to kuberntes cluster load balancer
func createCert(ctx *pulumi.Context) (Cert *digitalocean.Certificate, Error error) {
	cert, err := digitalocean.NewCertificate(ctx, "cert", &digitalocean.CertificateArgs{
		Domains: pulumi.StringArray{
			pulumi.String(domain),
		},
		Type: pulumi.String("lets_encrypt"),
	})
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// Deploy Traefik proxy as a service load balancer
func deployTraefik(ctx *pulumi.Context, certId pulumi.StringOutput) error {
	_, err := corev1.NewNamespace(ctx, "networking", &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("networking"),
		},
	})

	if err != nil {
		return err
	}

	values := pulumi.Map{
		"service": pulumi.Map{
			"annotations": pulumi.StringMap{
				"service.beta.kubernetes.io/do-loadbalancer-certificate-id":         certId,
				"service.beta.kubernetes.io/do-loadbalancer-protocol":               pulumi.String("https"),
				"service.beta.kubernetes.io/do-loadbalancer-hostname":               pulumi.String(domain),
				"service.beta.kubernetes.io/do-loadbalancer-redirect-http-to-https": pulumi.String("true"),
			},
		},
	}
	_, err = helm.NewChart(ctx, "traefik", helm.ChartArgs{
		Chart: pulumi.String("traefik"),

		FetchArgs: helm.FetchArgs{
			Repo: pulumi.String("https://helm.traefik.io/traefik"),
		},
		Namespace: pulumi.String("networking"),
		Values:    values,
		Version:   pulumi.String("10.7.1"),
	})

	if err != nil {
		return err
	}

	return nil
}

// Deploy Linkered
func deployLinkered(ctx *pulumi.Context) error {
	values := pulumi.Map{
		"clusterDomain":    pulumi.String(domain),
		"installNamespace": pulumi.Bool(false),
	}
	_, err := helm.NewChart(ctx, "linkered", helm.ChartArgs{
		Chart: pulumi.String("linkered"),

		FetchArgs: helm.FetchArgs{
			Repo: pulumi.String("https://helm.linkerd.io/stable"),
		},
		Namespace: pulumi.String("networking"),
		Version:   pulumi.String("2.11.1"),
		Values:    values,
	})
	if err != nil {
		return err
	}
	return nil
}

// Deploy ArgoCD
func deployArgoCD(ctx *pulumi.Context) error {
	_, err := corev1.NewNamespace(ctx, "argocd", &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("argocd"),
		},
	})
	if err != nil {
		return err
	}
	_, err = yaml.NewConfigFile(ctx, "argocd",
		&yaml.ConfigFileArgs{
			File: "argocd/argocd-v2.2.1.yaml",
			Transformations: []yaml.Transformation{
				func(state map[string]interface{}, opts ...pulumi.ResourceOption) {
					state["metadata"].(map[string]interface{})["namespace"] = "argocd"
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// Deploy Tekton
func deployTekton(ctx *pulumi.Context) error {
	_, err := yaml.NewConfigGroup(ctx, "tekton",
		&yaml.ConfigGroupArgs{
			Files: []string{
				"tekton/pipeline-v0.31.0.yaml",
				"tekton/triggers-v0.17.1.yaml",
				"tekton/interceptors-v0.17.1.yaml",
				"tekton/dashboard-v0.23.0.yaml",
			},
		},
	)

	_, err = yaml.NewConfigFile(ctx, "configmap-pv",
		&yaml.ConfigFileArgs{
			File: "tekton/configmap.yaml",
		},
	)
	if err != nil {
		return err
	}
	return nil
}
