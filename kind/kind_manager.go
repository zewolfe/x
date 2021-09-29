package kind

import (
	"github.com/pkg/errors"
	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/log"
)

type Manager struct{}

func CreateCluster(name string, version string) error {
	nodeImage := getImageForKubernetesVersion(version)
	provider := cluster.NewProvider(
		cluster.ProviderWithLogger(log.NoopLogger{}),
	)

	// create the cluster
	if err := provider.Create(
		name,
		cluster.CreateWithNodeImage(nodeImage),
		cluster.CreateWithDisplayUsage(true),
		cluster.CreateWithDisplaySalutation(true),
	); err != nil {
		return errors.Wrap(err, "failed to create cluster")
	}

	return nil
}

func DeleteCluster(name string) error {
	provider := cluster.NewProvider(
		cluster.ProviderWithLogger(log.NoopLogger{}),
	)

	// Delete individual cluster
	if err := provider.Delete(name, ""); err != nil {
		return errors.Wrapf(err, "failed to delete cluster %q", name)
	}
	return nil
}

func getImageForKubernetesVersion(version string) string {
	kindImageKubVersionsMap := map[string]string{
		"1.22": "kindest/node:v1.22.0@sha256:b8bda84bb3a190e6e028b1760d277454a72267a5454b57db34437c34a588d047",
		"1.21": "kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6",
		"1.20": "kindest/node:v1.20.7@sha256:cbeaf907fc78ac97ce7b625e4bf0de16e3ea725daf6b04f930bd14c67c671ff9",
		"1.19": "kindest/node:v1.19.11@sha256:07db187ae84b4b7de440a73886f008cf903fcf5764ba8106a9fd5243d6f32729",
		"1.18": "kindest/node:v1.18.19@sha256:7af1492e19b3192a79f606e43c35fb741e520d195f96399284515f077b3b622c",
		"1.17": "kindest/node:v1.17.17@sha256:66f1d0d91a88b8a001811e2f1054af60eef3b669a9a74f9b6db871f2f1eeed00",
	}

	return kindImageKubVersionsMap[version]
}
