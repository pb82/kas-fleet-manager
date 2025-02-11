package integration

import (
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/internal/kafka/test"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/internal/kafka/test/common"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/internal/kafka/test/mocks/kasfleetshardsync"
	"net/http"
	"testing"

	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/internal/kafka/constants"

	coreTest "github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/test"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/test/mocks"
	. "github.com/onsi/gomega"
	clustersmgmtv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

func TestClusterComputeNodesScaling(t *testing.T) {
	// create a mock ocm api server, keep all endpoints as defaults
	// see the mocks package for more information on the configurable mock server
	ocmServerBuilder := mocks.NewMockConfigurableServerBuilder()
	ocmServer := ocmServerBuilder.Build()
	defer ocmServer.Close()

	// setup the test environment, if OCM_ENV=integration then the ocmServer provided will be used instead of actual
	// ocm
	h, _, teardown := test.NewKafkaHelper(t, ocmServer)
	defer teardown()

	kasFleetshardSyncBuilder := kasfleetshardsync.NewMockKasFleetshardSyncBuilder(h, t)
	kasfFleetshardSync := kasFleetshardSyncBuilder.Build()
	kasfFleetshardSync.Start()
	defer kasfFleetshardSync.Stop()

	clusterID, getClusterErr := common.GetRunningOsdClusterID(h, t)
	if getClusterErr != nil {
		t.Fatalf("Failed to retrieve cluster details: %v", getClusterErr)
	}
	if clusterID == "" {
		panic("No cluster found")
	}

	expectedReplicas := mocks.MockClusterComputeNodes + constants.DefaultClusterNodeScaleIncrement

	overrideClusterMockResponse(ocmServerBuilder, expectedReplicas)

	scaleUpComputeNodes(h, expectedReplicas, clusterID, constants.DefaultClusterNodeScaleIncrement)

	expectedReplicas = expectedReplicas - constants.DefaultClusterNodeScaleIncrement

	overrideClusterMockResponse(ocmServerBuilder, expectedReplicas)

	scaleDownComputeNodes(h, expectedReplicas, clusterID, constants.DefaultClusterNodeScaleIncrement)
}

// get mock Cluster with specified Compute replicas number
func getClusterForScaleTest(replicas int) *clustersmgmtv1.Cluster {
	nodesBuilder := clustersmgmtv1.NewClusterNodes().
		Compute(replicas)
	mockClusterBuilder := mocks.GetMockClusterBuilder(func(builder *clustersmgmtv1.ClusterBuilder) {
		(*builder).Nodes(nodesBuilder)
	})
	cluster, err := mockClusterBuilder.Build()
	if err != nil {
		panic(err)
	}
	return cluster
}

// scaleUpComputeNodes and confirm that it is scaled without error
func scaleUpComputeNodes(h *coreTest.Helper, expectedReplicas int, clusterID string, increment int) {
	_, err := test.TestServices.ClusterService.ScaleUpComputeNodes(clusterID, increment)
	Expect(err).To(BeNil())
}

// scaleDownComputeNodes and confirm that it is scaled without error
func scaleDownComputeNodes(h *coreTest.Helper, expectedReplicas int, clusterID string, decrement int) {
	_, err := test.TestServices.ClusterService.ScaleDownComputeNodes(clusterID, decrement)
	Expect(err).To(BeNil())
}

// overrideClusterMockResponse - override mock response for Cluster patch
func overrideClusterMockResponse(ocmServerBuilder *mocks.MockConfigurableServerBuilder, expectedReplicas int) {
	mockCluster := getClusterForScaleTest(expectedReplicas)
	ocmServerBuilder.SwapRouterResponse(mocks.EndpointPathCluster, http.MethodPatch, mockCluster, nil)
}
