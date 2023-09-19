package render

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Workload + strategy
type ManifestWithStrategy struct {
	ManifestItem Manifest `json:"manifest,omitempty"`
}

type RollbackStrategy struct {
	// Minimum percentage of replicas that should be ready before the deployment can be considered for LKG after the minreadyseconds period
	// default is 100%
	MinReadyReplicas int32 `json:"minReadyReplicas,omitempty"`

	// Minimum time in seconds after the deployment is complete to mark the deployment version LKG
	MinReadySeconds int32 `json:"minReadySeconds,omitempty"`

	// Total Deployment timeout
	ProgressDeadlineSeconds int32 `json:"progressDeadlineSeconds,omitempty"`

	// The number of application versions to store for rollbacks, default is 5
	RevisionHistoryLimit int32 `json:"revisionHistoryLimit,omitempty"`
}

type RolloutStrategy struct {
	// Minimum percentage of replicas that should be available before we mark the deployment completed
	// +kubebuilder:default="100%"
	MinAvailableReplicas string `json:"minAvailableReplicas,omitempty"`
}

// ApplicationSpec defines the desired state of Application
type ApplicationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make generate-apis-and-code" to regenerate code after modifying this file

	Paused bool `json:"paused,omitempty"`

	// Workloads represents the MainfestWithStrategy to be deployed on a standard cluster.
	Workload []ManifestWithStrategy `json:"workload,omitempty"`
	// PacementRef is the name of the Placement resource, from which a PlacementDecision will be found and used
	// to distribute the ManifestWork
	// PlacementRef *corev1.LocalObjectReference `json:"placementRef,omitempty"`

	// Version represent the current application version
	Version string `json:"version,omitempty"`

	// Rollout Strategy represents the parameters that control rollout
	RolloutStrategy RolloutStrategy `json:"rolloutStrategy,omitempty"`

	// Rollback Strategy represents the parameters that control rollback
	RollbackStrategy RollbackStrategy `json:"rollbackStrategy,omitempty"`

	// Cluster selector based on the capabilities required by the application
	ClusterSelectors ClusterSelectors `json:"clusterSelectors,omitempty"`
}

type ClusterSelectors struct {
	// List of capabilities needed in the cluster
	Capabilities []string `json:"capabilities,omitempty"`
}

// +kubebuilder:validation:Enum=Progressing;Applied;Available;Degraded
type ApplicationState string

const (
	// ApplicationProgressing represents that the resource is processing on the control cluster
	ApplicationProgressing ApplicationState = "Progressing"
	// ApplicationApplied represents that the resource object is applied/splitted into standard clusters.
	// on the standard cluster.
	ApplicationApplied ApplicationState = "Applied"
	// ApplicationAvailable represents that the resource object exists on the standard cluster.
	ApplicationAvailable ApplicationState = "Available"
	// ApplicationDegraded represents that the current state of resource object does not match the desired state for a certain period.
	ApplicationDegraded ApplicationState = "Degraded"
)

// +kubebuilder:validation:Enum=Migrating;Working
type SingletonApplicationState string

const (
	// This happends when scheduler find manifestwork is degraded, scheduler will set statefulset replicas as 0 temporarily
	// and set SingletonApplicationState as SingletonApplicationMigrating, that means scheduler is dealing with and migrating
	// degraded manifestwork, there will be inconsistent between application and manifestwork, replicas in manifestwork is 0 temporarily.
	// When manifestwork is available again, SingletonApplicationState will be changed to SingletonApplicationWorking.
	SingletonApplicationMigrating SingletonApplicationState = "Migrating"
	// Default state for singleton application
	// This is normal state of singleton application.
	// It is an opposite state against SingletonApplicationMigrating, SingletonApplicationWorking represents that SingletonApplication is working normally,
	// Unlike state SingletonApplicationMigrating, the replicas in application and manifestwork is consistent.
	SingletonApplicationWorking SingletonApplicationState = "Working"
)

// +kubebuilder:validation:Enum=InProgress;Completed;RollingBack;RollingBackCompleted
type RolloutStatus string

const (
	// InProgress represents that new version of application is rolling out
	InProgress RolloutStatus = "InProgress"
	// Completed represents that new version of application has been rolled out successfully
	CompletedRollout RolloutStatus = "Completed"
	// RollingBack represents that application is being rolled back to LKG version
	RollingBack RolloutStatus = "RollingBack"
	// RollingBackCompleted represents that application rolled back to LKG successfully
	RollingBackCompleted RolloutStatus = "RollingBackCompleted"
)

type ApplicationClusterStatus struct {
	// Human readable name of the standard cluster
	Cluster string `json:"cluster,omitempty"`

	// Conditions contains the different condition statuses for this work.
	// Valid condition types are:
	// 1. Ready - All the Manifest resources satisfy the readiness conditions.
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// LastManifestStatusObservedTime is the last time the Manifests status is observed.
	LastManifestStatusObservedTime metav1.Time `json:"lastManifestStatusObservedTime,omitempty"`

	// The observed generation of Application which scheduled the ManifestWork in the cluster.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// Per Manifest's status in the standard cluster.
	ManifestStatuses []ManifestStatus `json:"manifestStatuses,omitempty"`
}

// ApplicationStatus defines the observed state of Application
type ApplicationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make generate-apis-and-code" to regenerate code after modifying this file

	// Status of the Application manifests in the scheduled standard clusters
	Clusters []ApplicationClusterStatus `json:"clusters,omitempty"`

	// Conditions contains the different condition statuses for distrbution of ManifestWork resources
	// Valid condition types are:
	// 1. Provisioned represents Application is scheduled on to the standard clusters.
	// 2. AppliedManifestWorks represents ManifestWorks have been distributed as per All, Partial, None, Problem
	// 3. Ready  represents ManifestWorks across all the standard clusters has "Ready" condition status as "True"
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// The generation observed by the application controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// LKG version for the application
	LastKnownGoodVersion string `json:"lastKnownGoodVersion,omitempty"`

	RolloutStatus RolloutStatus `json:"rolloutStatus,omitempty"`

	// 1. Applied represents workload in Application is applied/splitted successfully into standard clusters.
	// 2. Progressing represents workload in Application is being processed/splitted on control cluster.
	// 3. Available represents workload in Application exists on the standard clusters.
	// 4. Degraded represents the current state of workload does not match the desired state for a certain period.
	ApplicationState ApplicationState `json:"applicationState,omitempty"`
	// 1. SingletonApplicationMigrating represents that SingletonApplication is migrating for reason like manifestwork is degraded
	// 2. SingletonApplicationWorking is default state for singleton application, representing singleton application it normal.
	SingletonApplicationState SingletonApplicationState `json:"SingletonApplicationState,omitempty"`

	// Per Manifest's status in the standard cluster.
	ManifestStatuses []ManifestStatus `json:"manifestStatuses,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=application-api,shortName=app
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`
// +kubebuilder:printcolumn:name="LKG",type=string,JSONPath=`.status.lastKnownGoodVersion`
// +kubebuilder:printcolumn:name="RolloutStatus",type=string,JSONPath=`.status.rolloutStatus`
// Application is the Schema for the applications API
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// ApplicationList contains a list of Application
type ApplicationList struct {
	metav1.TypeMeta `              json:",inline"`
	metav1.ListMeta `              json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}
