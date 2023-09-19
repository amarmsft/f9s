package render

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterProvisioner is used by a cluster provisioner to provide top level configuration used by a provisioner for any standard cluster it creates
// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=category-api,shortName=clpro,scope=Cluster
// +kubebuilder:storageversion
// +genclient:nonNamespaced
type ClusterProvisioner struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Contains the information required by a cluster provisioner to operate.
	Spec ClusterProvisionerSpec `json:"spec,omitempty"`
}

// ClusterProvisionerList is a list of ClusterProvisioner resources
// +kubebuilder:object:root=true
type ClusterProvisionerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ClusterProvisioner `json:"items"`
}

// ClusterProvisionerSpec contains information required by a cluster provisioner to operate.
type ClusterProvisionerSpec struct {
	// Priority is used by provisioner resolution process for ordering of operations.
	Priority int `json:"priority,omitempty"`

	// Contains name of the default definition to be used by provisioner if a cluster definition does not provide specific definition.
	DefaultDefinitionName string `json:"defaultDefinitionName,omitempty"`

	// Contains map of configuration key/value pairs used by the provisioner.
	Properties map[string]string `json:"properties,omitempty"`
}

// ClusterDefinitionTask contains information for a task executed by Provisioner. A task is only understood by the provisioner, and the actual effect of the task can vary by provisioner.
type ClusterDefinitionTask struct {
	// Unique name for the task to be referenced by other tasks.
	Name string `json:"name"`

	// Type of the task that needs to be executed.
	TaskType string `json:"taskType"`

	// Properties contains map of configuration key/value pairs with information necessary for the task.
	Properties map[string]ClusterDefinitionTaskPropertyValue `json:"properties,omitempty"`
}

// ClusterDefinitionTaskPropertyValue contains information for Task Property Value. It can have a value or a value reference from a ClusterDefinitionTaskRef.
type ClusterDefinitionTaskPropertyValue struct {
	// Value of the property.
	Value string `json:"value,omitempty"`

	// Reference to output of another task.
	ValueFrom ClusterDefinitionTaskRef `json:"valueFrom,omitempty"`
}

// ClusterDefinitionTaskRef contains information for a task reference. Used by another task to get output value.
type ClusterDefinitionTaskRef struct {
	// Name of the Task.
	TaskName string `json:"taskName"`

	// Name of the generated output property of the referenced task.
	OutputProperty string `json:"outputProperty"`
}

// ClusterDefinition is a general purpose definition of a cluster, with the specific effects of it being determined by the provisioner.
// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=category-api,shortName=cldef,scope=Cluster
// +kubebuilder:storageversion
// +genclient:nonNamespaced
type ClusterDefinition struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Contains information required by a cluster provisioner to operate
	Spec ClusterDefinitionSpec `json:"spec,omitempty"`
}

// ClusterDefinitionList is a list of ClusterDefinition resources
// +kubebuilder:object:root=true
type ClusterDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ClusterDefinition `json:"items"`
}

// ClusterDefinitionSpec contains information required by a cluster provisioner to operate.
type ClusterDefinitionSpec struct {
	// Set of tasks provisioner needs to execute prior to creating cluster
	PreTasks []ClusterDefinitionTask `json:"preTasks,omitempty"`

	// Set of tasks provisioner needs to execute in order to create cluster
	Tasks []ClusterDefinitionTask `json:"tasks,omitempty"`

	// Set of tasks provisioner needs to execute after creating cluster
	PostTasks []ClusterDefinitionTask `json:"postTasks,omitempty"`

	// Properties contains map of configuration key/value pairs for creating clusters, not specific to a single cluster instance
	Properties map[string]string `json:"properties,omitempty"`
}

// Cluster contains information to create a cluster, identify a cluster, and provide its statuses.
// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=category-api,shortName=cl,scope=Cluster
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +genclient:nonNamespaced
type Cluster struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Contains information used when creating the cluster.
	Spec ClusterSpec `json:"spec,omitempty"`

	// Contains the status details for the cluster.
	Status ClusterStatus `json:"status,omitempty"`
}

// ClusterList is a list of Cluster resources
// +kubebuilder:object:root=true
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Cluster `json:"items"`
}

// ClusterSpec contains information used when creating the cluster.
type ClusterSpec struct {
	// Provisioner to be used for cluster.
	Provisioner string `json:"provisioner,omitempty"`

	// ClusterDefinition to be used for creating cluster.
	ClusterDefinition string `json:"clusterDefinition,omitempty"`

	// Contains map of property key/value pairs for the cluster
	Properties map[string]string `json:"properties,omitempty"`
}

// ClusterHealthStatus represents the health of the cluster.
type ClusterHealthStatus string

const (
	// HealthyClusterHealth indicates cluster can host and run existing and new workload.
	HealthyClusterHealth ClusterHealthStatus = "Healthy"

	// PartialFailedClusterHealth indicates cluster is up, running currently deployed workloads but should not be used to host new workloads.
	PartialFailedClusterHealth ClusterHealthStatus = "Partial-Failed"

	// FailedClusterHealth indicates cluster has failed, and workloads should migrate off cluster.
	FailedClusterHealth ClusterHealthStatus = "Failed"
)

// ClusterHealthStatus represents the health of the cluster.
type ClusterActivityStatus string

const (
	// ActiveClusterActivity indicates cluster can host and run existing and new workload.
	ActiveClusterActivity ClusterActivityStatus = "Active"

	// DrainClusterActivity indicates cluster is marked as Drain, and workloads should migrate off cluster.
	DrainClusterActivity ClusterActivityStatus = "Drain"
)

// ClusterStatus contains the status details for the cluster.
type ClusterStatus struct {
	// map of all identification keys used to identify this cluster (on top of its name).
	KeyIdentifier map[string]string `json:"keyIdentifiers,omitempty"`

	// lastStatusChange is a timestamp of the last status seen for this cluster
	// it represents the latest any of cluster nodes has been updated.
	LastStatusChange string `json:"lastStatusChange,omitempty"`

	// Health Policy to be used for determining cluster status.
	ClusterHealthPolicy string `json:"clusterHealthPolicy,omitempty"`

	// ClusterStatus is calculated by LifecycleController based on node status and ClusterHealthPolicy.
	ClusterHealthStatus ClusterHealthStatus `json:"clusterStatus,omitempty"`

	// ClusterActivityStatus is used to indicate wheather the cluster is Active or Drain.
	ClusterActivityStatus ClusterActivityStatus `json:"clusterActivityStatus,omitempty"`

	// Contains status information for all steps/tasks during cluster creation. That information
	//  is then used for cluster deletion. This status should only be updated
	//  by a Provisioner.
	ProvisioningStatus []StepStatus `json:"provisioningStatus,omitempty"`

	// Runtime status for the cluster and all of its nodes.
	RuntimeStatus RuntimeStatus `json:"runtimeStatus,omitempty"`
}

// StepStatus contains the status for a step/task performed by a provisioner.
type StepStatus struct {
	Name         string            `json:"name,omitempty"`
	Type         string            `json:"type,omitempty"`
	Status       string            `json:"status,omitempty"`
	ErrorMessage string            `json:"errorMessage,omitempty"`
	Properties   map[string]string `json:"properties,omitempty"`
}

// RuntimeStatus contains key runtime details of a cluster.
type RuntimeStatus struct {
	// ClusterVersion contains Kubernetes cluster version information.
	ClusterVersion ClusterVersion `json:"clusterVersion,omitempty"`

	// NodeStatus contains an overview of healthy/unhealthy/total node counts.
	NodeStatus NodeStatus `json:"nodeStatus,omitempty"`

	// OsDistribution contains list of all OSes discovered within the cluster.
	OsDistribution OsDistribution `json:"osDistribution,omitempty"`

	ClusterStateOverride bool `json:"clusterStateOverride,omitempty"`

	// Cluster state based on cluster health report
	ClusterState ClusterState `json:"clusterState,omitempty"`

	// Cluster sub state based on cluster health report
	ClusterSubstate ClusterSubstate `json:"clusterSubstate,omitempty"`

	// The time when the state was changed
	// +kubebuilder:validation:Format=date-time
	LastClusterStateChangeTime string `json:"lastClusterStateChangeTime,omitempty"`

	// dictionary of capabilities and their status
	ClusterCapabilitiesStatus map[ComponentCapability]PolicyStatus `json:"clusterCapabilitiesStatus,omitempty"`

	// dictionary of capabilities and their status per OS
	ClusterCapabilitiesStatusPerOs map[OperatingSystem]map[ComponentCapability]PolicyStatus `json:"clusterCapabilitiesStatusPerOs,omitempty"`

	// The time when the capabilities status was computed
	// +kubebuilder:validation:Format=date-time
	LastCapabilitiesStatusTransitionTime string `json:"lastCapabilitiesStatusTransitionTime,omitempty"`
}

// +kubebuilder:validation:Enum=Provisioned;Ready;Failed;Drain;Delete;Unknown
type ClusterState string

// +kubebuilder:validation:Enum=Running;Degraded
type ClusterSubstate string

const (
	// cluster state as provisioned
	ClusterStateProvisioned ClusterState = "Provisioned"

	// cluster state as Ready
	ClusterStateReady ClusterState = "Ready"

	// cluster state as Failed
	ClusterStateFailed ClusterState = "Failed"

	// cluster state as Drain
	ClusterStateDrain ClusterState = "Drain"

	// cluster state as Delete
	ClusterStateDelete ClusterState = "Delete"

	// cluster state as Unknown
	ClusterStateUnknown ClusterState = "Unknown"
)

const (
	// cluster sub state as running
	ClusterSubstateRunning ClusterSubstate = "Running"

	// cluster sub state as degraded
	ClusterSubstateDegraded ClusterSubstate = "Degraded"
)

// NodeStatus contains the Healthy and Totalnodes counts for a cluster.
type NodeStatus struct {
	HealthyNodeCount int `json:"healthyNodeCount,omitempty"`
	TotalNodeCount   int `json:"totalNodeCount,omitempty"`
}

// ClusterVersion contains Kubernetes cluster version information for a cluster.
type ClusterVersion struct {
	ControlPlaneVersion string         `json:"controlPlaneVersion,omitempty"`
	NodeVersions        map[string]int `json:"nodeVersions,omitempty"`
}

// OsDistribution contains a summary of OS details for cluster nodes.
type OsDistribution struct {
	Os          string `json:"os,omitempty"`
	Type        string `json:"type,omitempty"`
	VersionName string `json:"versionName,omitempty"`
	NodeCount   int    `json:"nodeCount,omitempty"`
}

// The FleetDisruption ensures that application updates are rolled out with an acceptable number of removed clusters
// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=category-api,shortName=fdb
// +kubebuilder:storageversion
type FleetDisruption struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Contains information when creating the FleetDisruption
	Spec FleetDisruptionSpec `json:"spec,omitempty"`
}

type FleetDisruptionSpec struct {
	// Selector to which resources are matched
	ClusterSelector map[string]string `json:"clusterSelector,omitempty"`

	// This number reflects the # of clusters that an update will be rolled out to in case of a rolling update.
	// Validated as int or int% sign. When % is used then it uses the ceiling of the value:
	//		e.g., ceiling(101 * 0.10).
	// Example:
	//		A resource is deployed to 10 clusters. MaxUnavailble is 2.
	//		The resource will be applied on two clusters over 5 steps.
	// If a cluster is already down, the rollout will go one cluster at a time until the cluster recovers
	MaxUnavailable string `json:"maxUnavailable,omitempty"`
}

// FleetDisruption is a list of FleetDisruption resources
// +kubebuilder:object:root=true
type FleetDisruptionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []FleetDisruption `json:"items"`
}

// The FineGrainedSyncer ensures that application dependencies are provisioned with applications to standard clusters
// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=category-api,shortName=fgs
// +kubebuilder:storageversion
type FineGrainedSyncer struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Contains information when creating the FineGrainedSyncer
	Spec FineGrainedSyncerSpec `json:"spec,omitempty"`
}

// FineCraingedSyncerSpec contains information when creating the FineGrainedSyncer
type FineGrainedSyncerSpec struct {
	// Selector to which resources are matched
	Selector map[string]string `json:"selector,omitempty"`

	// A list of fleet resources to sync with services
	Follow []string `json:"follow,omitempty"`
}

// FineGrainedSyncer is a list of FineGrainedSyncer resources
// +kubebuilder:object:root=true
type FineGrainedSyncerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []FineGrainedSyncer `json:"items"`
}

// CCS
// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=category-api,shortName=ccs
// +kubebuilder:storageversion
type CrossClusterService struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Contains information used when creating the cluster.
	Spec CrossClusterServiceSpec `json:"spec,omitempty"`
}

// CrossClusterServiceSpec contains the information used for deploying a Cross Cluster Service
type CrossClusterServiceSpec struct {
	// Enum:
	// Load-Balanced (DEFAULT): Service load balance across all endpoints that matches this service.
	// Sharded: service represents “book-keeping” for all endpoints that matches this service
	Type string `json:"provisioner,omitempty"`

	// note on selectors: we cannot use traditional selectors due to https://github.com/kubernetes/kubernetes/issues/53459
	// cluster selector are labels used to match clusters this allows for scoping the service and its endpoint to be read and routed only on cluster that match certain labels e.g., env:Prod.
	ClusterSelector []map[string]string `json:"clusterSelector,omitempty"`

	// pod selectors are used to match pods to this cross cluster service
	PodSelector []map[string]string `json:"podSelector,omitempty"`

	// Applicable only in case of sharded services. Must be nil in case of load-balanced service.
	ShardResolver ShardResolver `json:"shardResolver,omitempty"`
}

// ShardResolver is used for deciding which shard a service should belong to.
type ShardResolver struct {
	// Enum:
	// byLabel: (DEFAULT) shared name is available as a label on each of the matching pods.
	// This allows destination to be shaped as multiple of deployments or statefulsets or
	// any other type of Kubernetes application models as long # as pods are labeled
	// correctly within each deployment.
	// byOrdinal: shard name is available as an ordinal in pod name.
	// Typically as statefulset pod. For example, a statefulset named store will yield into
	// pods named “store-1”, “store-2”, “store-3”; “1”, “2” and “3” will become shared names.
	ResolverType string `json:"type,omitempty"`

	// label name carrying shard name. must be empty string in case of “byOrdinal”
	LabelName string `json:"labelName,omitempty"`
}

// CrossClusterService is a list of CrossClusterService resources
// +kubebuilder:object:root=true
type CrossClusterServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CrossClusterService `json:"items"`
}

// CrossClusterShard defines a shard for a CrossClusterService
// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=category-api,shortName=ccsh
// +kubebuilder:storageversion
type CrossClusterShard struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Contains information when creating the cross cluster shard
	Spec CrossClusterShardSpec `json:"spec,omitempty"`
}

type CrossClusterShardSpec struct {
	// CCS service that owns the shard
	ServiceName string `json:"serviceName,omitempty"`
}

// CrossClusterShard is a list of CrossClusterShard resources
// +kubebuilder:object:root=true
type CrossClusterShardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CrossClusterShard `json:"items"`
}
