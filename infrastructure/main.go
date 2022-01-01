package main

import (
	"log"

	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cluster, err := digitalocean.NewKubernetesCluster(ctx, "do-challenge", &digitalocean.KubernetesClusterArgs{
			NodePool: &digitalocean.KubernetesClusterNodePoolArgs{
				Name:      pulumi.String("ops-pool"),
				Size:      pulumi.String("s-2vcpu-4gb"),
				MinNodes:  pulumi.Int(2),
				NodeCount: pulumi.Int(2),
			},
			Region:  pulumi.String("sfo3"),
			Version: pulumi.String("1.21.5-do.0"),
		})

		if err != nil {
			log.Fatal(err)
		}

		ctx.Export("clusterIP", cluster.Ipv4Address)
		ctx.Export("clusterID", cluster.ID())
		ctx.Export("clusterName", cluster.Name)

		return nil
	})
}
