package sqlite

import (
	"context"
	"os"
	"testing"
	
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	
	"github.com/operator-framework/operator-registry/pkg/registry"
)

func TestConfigMapLoader(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	store, err := NewSQLLiteLoader("test.db")
	require.NoError(t, err)
	defer os.Remove("test.db")

	path := "../../configmap.example.yaml"
	fileReader, err := os.Open(path)
	require.NoError(t, err, "unable to load configmap from file %s", path)

	decoder := yaml.NewYAMLOrJSONDecoder(fileReader, 30)
	manifest := v1.ConfigMap{}
	err = decoder.Decode(&manifest)
	require.NoError(t, err, "could not decode contents of file %s into configmap", path)

	loader := NewSQLLoaderForConfigMap(store, manifest)
	require.NoError(t, loader.Populate())
}

func TestQuerierForConfigmap(t *testing.T) {
	load, err := NewSQLLiteLoader("test.db")
	require.NoError(t, err)
	defer os.Remove("test.db")

	path := "../../configmap.example.yaml"
	fileReader, err := os.Open(path)
	require.NoError(t, err, "unable to load configmap from file %s", path)

	decoder := yaml.NewYAMLOrJSONDecoder(fileReader, 30)
	manifest := v1.ConfigMap{}
	err = decoder.Decode(&manifest)
	require.NoError(t, err, "could not decode contents of file %s into configmap", path)

	loader := NewSQLLoaderForConfigMap(load, manifest)
	require.NoError(t, loader.Populate())

	store, err := NewSQLLiteQuerier("test.db")
	require.NoError(t, err)

	foundPackages, err := store.ListPackages(context.TODO())
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"elasticsearch-operator"}, foundPackages)

	etcdPackage, err := store.GetPackage(context.TODO(), "elasticsearch-operator")
	require.NoError(t, err)
	require.EqualValues(t, &registry.PackageManifest{
		PackageName:        "elasticsearch-operator",
		DefaultChannelName: "4.2",
		Channels: []registry.PackageChannel{
			{
				Name:           "4.2",
				CurrentCSVName: "elasticsearch-operator.v4.2.0",
			},
		},
	}, etcdPackage)

	etcdBundleByChannel, err := store.GetBundleForChannel(context.TODO(), "elasticsearch-operator", "4.2")
	require.NoError(t, err)
	expectedBundle :=  "{\"apiVersion\":\"operators.coreos.com/v1alpha1\",\"kind\":\"ClusterServiceVersion\",\"metadata\":{\"annotations\":{\"alm-examples\":\"[\\n    {\\n        \\\"apiVersion\\\": \\\"logging.openshift.io/v1\\\",\\n        \\\"kind\\\": \\\"Elasticsearch\\\",\\n        \\\"metadata\\\": {\\n          \\\"name\\\": \\\"elasticsearch\\\"\\n        },\\n        \\\"spec\\\": {\\n          \\\"managementState\\\": \\\"Managed\\\",\\n          \\\"nodeSpec\\\": {\\n            \\\"image\\\": \\\"quay.io/openshift/origin-logging-elasticsearch5:latest\\\",\\n            \\\"resources\\\": {\\n              \\\"limits\\\": {\\n                \\\"memory\\\": \\\"1Gi\\\"\\n              },\\n              \\\"requests\\\": {\\n                \\\"memory\\\": \\\"512Mi\\\"\\n              }\\n            }\\n          },\\n          \\\"redundancyPolicy\\\": \\\"SingleRedundancy\\\",\\n          \\\"nodes\\\": [\\n            {\\n                \\\"nodeCount\\\": 1,\\n                \\\"roles\\\": [\\\"client\\\",\\\"data\\\",\\\"master\\\"]\\n            }\\n          ]\\n        }\\n    }\\n]\",\"capabilities\":\"Seamless Upgrades\",\"categories\":\"OpenShift Optional, Logging \\u0026 Tracing\",\"certified\":\"false\",\"containerImage\":\"quay.io/openshift/origin-elasticsearch-operator:latest\",\"createdAt\":\"2019-02-20T08:00:00Z\",\"description\":\"The Elasticsearch Operator for OKD provides a means for configuring and managing an Elasticsearch cluster for tracing and cluster logging.\\n## Prerequisites and Requirements\\n### Elasticsearch Operator Namespace\\nThe Elasticsearch Operator must be deployed to the global operator group namespace\\n### Memory Considerations\\nElasticsearch is a memory intensive application.  The initial\\nset of OKD nodes may not be large enough to support the Elasticsearch cluster.  Additional OKD nodes must be added\\nto the OKD cluster if you desire to run with the recommended (or better) memory. Each ES node can operate with a\\nlower memory setting though this is not recommended for production deployments.\",\"olm.skipRange\":\"\\u003e=4.1.0 \\u003c4.2.0\",\"support\":\"AOS Cluster Logging, Jaeger\"},\"creationTimestamp\":null,\"name\":\"elasticsearch-operator.v4.2.0\",\"namespace\":\"openshift-logging\"},\"spec\":{\"customresourcedefinitions\":{\"owned\":[{\"description\":\"An Elasticsearch cluster instance\",\"displayName\":\"Elasticsearch\",\"kind\":\"Elasticsearch\",\"name\":\"elasticsearches.logging.openshift.io\",\"resources\":[{\"kind\":\"Deployment\",\"version\":\"v1\"},{\"kind\":\"StatefulSet\",\"version\":\"v1\"},{\"kind\":\"ReplicaSet\",\"version\":\"v1\"},{\"kind\":\"Pod\",\"version\":\"v1\"},{\"kind\":\"ConfigMap\",\"version\":\"v1\"},{\"kind\":\"Service\",\"version\":\"v1\"},{\"kind\":\"Route\",\"version\":\"v1\"}],\"specDescriptors\":[{\"description\":\"Limits describes the minimum/maximum amount of compute resources required/allowed\",\"displayName\":\"Resource Requirements\",\"path\":\"nodeSpec.resources\",\"x-descriptors\":[\"urn:alm:descriptor:com.tectonic.ui:resourceRequirements\"]}],\"statusDescriptors\":[{\"description\":\"The current Status of the Elasticsearch Cluster\",\"displayName\":\"Status\",\"path\":\"cluster.status\",\"x-descriptors\":[\"urn:alm:descriptor:io.kubernetes.phase\"]},{\"description\":\"The number of Active Primary Shards for the Elasticsearch Cluster\",\"displayName\":\"Active Primary Shards\",\"path\":\"cluster.activePrimShards\",\"x-descriptors\":[\"urn:alm:descriptor:text\"]},{\"description\":\"The number of Active Shards for the Elasticsearch Cluster\",\"displayName\":\"Active Shards\",\"path\":\"cluster.activeShards\",\"x-descriptors\":[\"urn:alm:descriptor:text\"]},{\"description\":\"The number of Initializing Shards for the Elasticsearch Cluster\",\"displayName\":\"Initializing Shards\",\"path\":\"cluster.initializingShards\",\"x-descriptors\":[\"urn:alm:descriptor:text\"]},{\"description\":\"The number of Data Nodes for the Elasticsearch Cluster\",\"displayName\":\"Number of Data Nodes\",\"path\":\"cluster.numDataNodes\",\"x-descriptors\":[\"urn:alm:descriptor:text\"]},{\"description\":\"The number of Nodes for the Elasticsearch Cluster\",\"displayName\":\"Number of Nodes\",\"path\":\"cluster.numNodes\",\"x-descriptors\":[\"urn:alm:descriptor:text\"]},{\"description\":\"The number of Relocating Shards for the Elasticsearch Cluster\",\"displayName\":\"Relocating Shards\",\"path\":\"cluster.relocatingShards\",\"x-descriptors\":[\"urn:alm:descriptor:text\"]},{\"description\":\"The number of Unassigned Shards for the Elasticsearch Cluster\",\"displayName\":\"Unassigned Shards\",\"path\":\"cluster.unassignedShards\",\"x-descriptors\":[\"urn:alm:descriptor:text\"]},{\"description\":\"The status for each of the Elasticsearch pods with the Client role\",\"displayName\":\"Elasticsearch Client Status\",\"path\":\"pods.client\",\"x-descriptors\":[\"urn:alm:descriptor:com.tectonic.ui:podStatuses\"]},{\"description\":\"The status for each of the Elasticsearch pods with the Data role\",\"displayName\":\"Elasticsearch Data Status\",\"path\":\"pods.data\",\"x-descriptors\":[\"urn:alm:descriptor:com.tectonic.ui:podStatuses\"]},{\"description\":\"The status for each of the Elasticsearch pods with the Master role\",\"displayName\":\"Elasticsearch Master Status\",\"path\":\"pods.master\",\"x-descriptors\":[\"urn:alm:descriptor:com.tectonic.ui:podStatuses\"]}],\"version\":\"v1\"}]},\"description\":\"The Elasticsearch Operator for OKD provides a means for configuring and managing an Elasticsearch cluster for use in tracing and cluster logging.\\nThis operator only supports OKD Cluster Logging and Jaeger.  It is tightly coupled to each and is not currently capable of\\nbeing used as a general purpose manager of Elasticsearch clusters running on OKD.\\n\\nIt is recommended this operator be deployed to the **openshift-operators** namespace to properly support the Cluster Logging and Jaeger use cases.\\n\\nOnce installed, the operator provides the following features:\\n* **Create/Destroy**: Deploy an Elasticsearch cluster to the same namespace in which the Elasticsearch custom resource is created.\\n\",\"displayName\":\"Elasticsearch Operator\",\"install\":{\"spec\":{\"clusterPermissions\":[{\"rules\":[{\"apiGroups\":[\"logging.openshift.io\"],\"resources\":[\"*\"],\"verbs\":[\"*\"]},{\"apiGroups\":[\"\"],\"resources\":[\"pods\",\"pods/exec\",\"services\",\"endpoints\",\"persistentvolumeclaims\",\"events\",\"configmaps\",\"secrets\",\"serviceaccounts\"],\"verbs\":[\"*\"]},{\"apiGroups\":[\"apps\"],\"resources\":[\"deployments\",\"daemonsets\",\"replicasets\",\"statefulsets\"],\"verbs\":[\"*\"]},{\"apiGroups\":[\"monitoring.coreos.com\"],\"resources\":[\"prometheusrules\",\"servicemonitors\"],\"verbs\":[\"*\"]},{\"apiGroups\":[\"rbac.authorization.k8s.io\"],\"resources\":[\"clusterroles\",\"clusterrolebindings\"],\"verbs\":[\"*\"]},{\"nonResourceURLs\":[\"/metrics\"],\"verbs\":[\"get\"]},{\"apiGroups\":[\"authentication.k8s.io\"],\"resources\":[\"tokenreviews\",\"subjectaccessreviews\"],\"verbs\":[\"create\"]},{\"apiGroups\":[\"authorization.k8s.io\"],\"resources\":[\"subjectaccessreviews\"],\"verbs\":[\"create\"]}],\"serviceAccountName\":\"elasticsearch-operator\"}],\"deployments\":[{\"name\":\"elasticsearch-operator\",\"spec\":{\"replicas\":1,\"selector\":{\"matchLabels\":{\"name\":\"elasticsearch-operator\"}},\"template\":{\"metadata\":{\"labels\":{\"name\":\"elasticsearch-operator\"}},\"spec\":{\"containers\":[{\"command\":[\"elasticsearch-operator\"],\"env\":[{\"name\":\"WATCH_NAMESPACE\",\"valueFrom\":{\"fieldRef\":{\"fieldPath\":\"metadata.annotations['olm.targetNamespaces']\"}}},{\"name\":\"POD_NAME\",\"valueFrom\":{\"fieldRef\":{\"fieldPath\":\"metadata.name\"}}},{\"name\":\"OPERATOR_NAME\",\"value\":\"elasticsearch-operator\"},{\"name\":\"PROXY_IMAGE\",\"value\":\"quay.io/openshift/origin-oauth-proxy:latest\"},{\"name\":\"ELASTICSEARCH_IMAGE\",\"value\":\"quay.io/openshift/origin-logging-elasticsearch5:latest\"}],\"image\":\"quay.io/openshift/origin-elasticsearch-operator:latest\",\"imagePullPolicy\":\"IfNotPresent\",\"name\":\"elasticsearch-operator\",\"ports\":[{\"containerPort\":60000,\"name\":\"metrics\"}]}],\"serviceAccountName\":\"elasticsearch-operator\"}}}}]},\"strategy\":\"deployment\"},\"installModes\":[{\"supported\":true,\"type\":\"OwnNamespace\"},{\"supported\":false,\"type\":\"SingleNamespace\"},{\"supported\":false,\"type\":\"MultiNamespace\"},{\"supported\":true,\"type\":\"AllNamespaces\"}],\"keywords\":[\"elasticsearch\",\"jaeger\"],\"links\":[{\"name\":\"Elastic\",\"url\":\"https://www.elastic.co/\"},{\"name\":\"Elasticsearch Operator\",\"url\":\"https://github.com/openshift/elasticsearch-operator\"}],\"maintainers\":[{\"email\":\"aos-logging@redhat.com\",\"name\":\"Red Hat, AOS Logging\"}],\"provider\":{\"name\":\"Red Hat, Inc\"},\"version\":\"4.2.0\"}}\n{\"apiVersion\":\"apiextensions.k8s.io/v1beta1\",\"kind\":\"CustomResourceDefinition\",\"metadata\":{\"creationTimestamp\":null,\"name\":\"elasticsearches.logging.openshift.io\"},\"spec\":{\"group\":\"logging.openshift.io\",\"names\":{\"kind\":\"Elasticsearch\",\"listKind\":\"ElasticsearchList\",\"plural\":\"elasticsearches\",\"singular\":\"elasticsearch\"},\"scope\":\"Namespaced\",\"validation\":{\"openAPIV3Schema\":{\"properties\":{\"spec\":{\"description\":\"Specification of the desired behavior of the Elasticsearch cluster\",\"properties\":{\"managementState\":{\"description\":\"Indicator if the resource is 'Managed' or 'Unmanaged' by the operator\",\"enum\":[\"Managed\",\"Unmanaged\"],\"type\":\"string\"},\"nodeSpec\":{\"description\":\"Default specification applied to all Elasticsearch nodes\",\"properties\":{\"image\":{\"description\":\"The image to use for the Elasticsearch nodes\",\"type\":\"string\"},\"resources\":{\"description\":\"The resource requirements for the Elasticsearch nodes\",\"properties\":{\"limits\":{\"description\":\"Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/\",\"type\":\"object\"},\"requests\":{\"description\":\"Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. More info: https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/\",\"type\":\"object\"}}}}},\"nodes\":{\"description\":\"Specification of the different Elasticsearch nodes\",\"items\":{\"properties\":{\"nodeCount\":{\"description\":\"Number of nodes to deploy\",\"format\":\"int32\",\"type\":\"integer\"},\"nodeSelector\":{\"description\":\"Define which Nodes the Pods are scheduled on.\",\"type\":\"object\"},\"nodeSpec\":{\"description\":\"Specification of a specific Elasticsearch node\",\"properties\":{\"image\":{\"description\":\"The image to use for the Elasticsearch node\",\"type\":\"string\"},\"resources\":{\"description\":\"The resource requirements for the Elasticsearch node\",\"properties\":{\"limits\":{\"description\":\"Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/\",\"type\":\"object\"},\"requests\":{\"description\":\"Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. More info: https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/\",\"type\":\"object\"}}}}},\"roles\":{\"description\":\"The specific Elasticsearch cluster roles the node should perform\",\"items\":{\"type\":\"string\"},\"type\":\"array\"},\"storage\":{\"description\":\"The type of backing storage that should be used for the node\",\"properties\":{\"size\":{\"description\":\"The max storage capacity for the node\",\"type\":\"string\"},\"storageClassName\":{\"description\":\"The name of the storage class to use with creating the node's PVC\",\"type\":\"string\"}}}},\"type\":\"object\"},\"type\":\"array\"},\"redundancyPolicy\":{\"description\":\"The policy towards data redundancy to specify the number of redundant primary shards\",\"enum\":[\"FullRedundancy\",\"MultipleRedundancy\",\"SingleRedundancy\",\"ZeroRedundancy\"],\"type\":\"string\"}}}}}},\"version\":\"v1\",\"versions\":[{\"name\":\"v1\",\"served\":true,\"storage\":true}]},\"status\":{\"acceptedNames\":{\"kind\":\"\",\"plural\":\"\"},\"conditions\":null,\"storedVersions\":null}}\n"
	require.Equal(t, expectedBundle, etcdBundleByChannel)

	etcdBundle, err := store.GetBundle(context.TODO(), "elasticsearch-operator", "4.2", "elasticsearch-operator.v4.2.0")
	require.NoError(t, err)
	require.Equal(t, expectedBundle, etcdBundle)

	// etcdChannelEntries, err := store.GetChannelEntriesThatReplace(context.TODO(), "etcdoperator.v0.9.0")
	// require.NoError(t, err)
	// require.ElementsMatch(t, []*registry.ChannelEntry{{"etcd", "alpha", "etcdoperator.v0.9.2", "etcdoperator.v0.9.0"}}, etcdChannelEntries)

	// etcdBundleByReplaces, err := store.GetBundleThatReplaces(context.TODO(), "etcdoperator.v0.9.0", "etcd", "alpha")
	// require.NoError(t, err)
	// require.EqualValues(t, expectedBundle, etcdBundleByReplaces)

	// etcdChannelEntriesThatProvide, err := store.GetChannelEntriesThatProvide(context.TODO(), "etcd.database.coreos.com", "v1beta2", "EtcdCluster")
	// require.ElementsMatch(t, []*registry.ChannelEntry{
	// 	{"etcd", "alpha", "etcdoperator.v0.6.1", ""},
	// 	{"etcd", "alpha", "etcdoperator.v0.9.0", "etcdoperator.v0.6.1"},
	// 	{"etcd", "alpha", "etcdoperator.v0.9.2", "etcdoperator.v0.9.0"}}, etcdChannelEntriesThatProvide)

	// etcdChannelEntriesThatProvideAPIServer, err := store.GetChannelEntriesThatProvide(context.TODO(), "etcd.database.coreos.com", "v1beta2", "FakeEtcdObject")
	// require.ElementsMatch(t, []*registry.ChannelEntry{{"etcd", "alpha", "etcdoperator.v0.9.0", "etcdoperator.v0.6.1"}}, etcdChannelEntriesThatProvideAPIServer)
	//
	// etcdLatestChannelEntriesThatProvide, err := store.GetLatestChannelEntriesThatProvide(context.TODO(), "etcd.database.coreos.com", "v1beta2", "EtcdCluster")
	// require.NoError(t, err)
	// require.ElementsMatch(t, []*registry.ChannelEntry{{"etcd", "alpha", "etcdoperator.v0.9.2", "etcdoperator.v0.9.0"}}, etcdLatestChannelEntriesThatProvide)
	//
	// etcdBundleByProvides, entry, err := store.GetBundleThatProvides(context.TODO(), "etcd.database.coreos.com", "v1beta2", "EtcdCluster")
	// require.NoError(t, err)
	// require.Equal(t, expectedBundle, etcdBundleByProvides)
	// require.Equal(t, &registry.ChannelEntry{"etcd", "alpha","etcdoperator.v0.9.2", ""}, entry)
}
