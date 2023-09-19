package render

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=manifestwork-api,shortName={manifestwork,mw}
// +kubebuilder:storageversion
// ManifestWork represents a manifests workload that hub wants to deploy on the standard cluster.
// A manifest workload is defined as a set of Kubernetes resources.
// ManifestWork must be created in the cluster namespace on the hub, so that agent on the
// corresponding standard cluster can access this resource and deploy on the standard
// cluster.
type ManifestWork struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec represents a desired configuration of work to be deployed on the standard cluster.
	Spec ManifestWorkSpec `json:"spec"`

	// Status represents the current status of work.
	// +optional
	Status ManifestWorkStatus `json:"status,omitempty"`
}

// ManifestWorkSpec represents a desired configuration of manifests to be deployed on the standard cluster.
type ManifestWorkSpec struct {
	// Workload represents the manifest workload to be deployed on a standard cluster.
	Workload ManifestsTemplate `json:"workload,omitempty"`
}

// Manifest represents a resource to be deployed on standard cluster.
type Manifest struct {
	// +kubebuilder:validation:EmbeddedResource
	// +kubebuilder:pruning:PreserveUnknownFields
	runtime.RawExtension `json:",inline"`
}

// ManifestsTemplate represents the manifest workload to be deployed on a standard cluster.
type ManifestsTemplate struct {
	// Manifests represents a list of kuberenetes resources to be deployed on a standard cluster.
	// +optional
	Manifests []Manifest `json:"manifests,omitempty"`
}

// +kubebuilder:validation:Enum=Progressing;Applied;Available;Degraded
type ManifestWorkState string

const (
	// ManifestProgressing represents that the resource is being applied on the standard cluster
	ManifestProgressing ManifestWorkState = "Progressing"
	// ManifestApplied represents that the resource object is applied
	// on the standard cluster.
	ManifestApplied ManifestWorkState = "Applied"
	// ManifestAvailable represents that the resource object exists
	// on the standard cluster.
	ManifestAvailable ManifestWorkState = "Available"
	// ManifestDegraded represents that the current state of resource object does not
	// match the desired state for a certain period.
	ManifestDegraded ManifestWorkState = "Degraded"
)

type ManifestDegradedReason string

const (
	// ManifestDegradedNoEnoughResource represents that the current state of the resource object
	// does not match the desired state for a certain period,
	// and the reason is no enough resource.
	ManifestDegradedNoEnoughResource ManifestDegradedReason = "NoEnoughResouce"
)

// Status of the Manifest in a standard cluster
type ManifestStatus struct {
	// Name of the Maniest
	Name string `json:"name,omitempty"`

	// Namespace where the Manifest is created.
	Namespace string `json:"namespace,omitempty"`

	// Manifest kind
	Kind string `json:"kind,omitempty"`

	// Status section of the Manifest resource in standard cluster
	// +kubebuilder:pruning:PreserveUnknownFields
	Status runtime.RawExtension `json:"status,omitempty"`
}

// ManifestWorkStatus represents the current status of standard cluster ManifestWork.
type ManifestWorkStatus struct {
	// The generation observed by the manifestwork controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// Conditions contains the different condition statuses for this work.
	// Valid condition types are:
	// 1. Applied represents workload in ManifestWork is applied successfully on standard cluster.
	// 2. Ready - All the Manifest resources satisfy the readiness conditions.
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// 1. Applied represents workload in ManifestWork is applied successfully on standard cluster.
	// 2. Progressing represents workload in ManifestWork is being applied on standard cluster.
	// 3. Available represents workload in ManifestWork exists on the standard cluster.
	// 4. Degraded represents the current state of workload does not match the desired
	// state for a certain period.
	ManifestState ManifestWorkState `json:"manifestState,omitempty"`

	// LastManifestStatusObservedTime is the last time the ManifestWorkStatus is observed.
	LastManifestStatusObservedTime metav1.Time `json:"lastManifestStatusObservedTime,omitempty"`

	// Per Manifest's status in the standard cluster.
	ManifestStatuses []ManifestStatus `json:"manifestStatuses,omitempty"`

	// health pods,
	// For each Stae value:
	// 1. Applied:   HealthReplicas = 0
	// 2. Progressing:  0 <= HealthReplicas <= spec.replicas
	// 3. Available:  HealthReplicas = spec.replicas
	// 4. Degraded: 0 <= HealthReplicas <= spec.replicas
	HealthReplicas int `json:"healthReplicas,omitempty"`

	//Not schedulable pending replicas indicate not enough resource
	NonSchedulableReplicas int `json:"nonSchedulableReplicas,omitempty"`

	//total pods with the manifestwork revision
	TotalReplicas int `json:"totalReplicas,omitempty"`

	// why the State is Degraded
	DegradedReason ManifestDegradedReason `json:"degradedReason,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ManifestWorkList is a collection of manifestworks.
type ManifestWorkList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of manifestworks.
	Items []ManifestWork `json:"items"`
}
