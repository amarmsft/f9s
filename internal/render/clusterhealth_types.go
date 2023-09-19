package render

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

// ClusterHealthPolicyList is a list of ClusterHealthPolicy resources
// +kubebuilder:object:root=true
type ClusterHealthPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ClusterHealthPolicy `json:"items"`
}

// ClusterHealthPolicy defines the policies to compute health of the mentioned resources/components.
// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=clusterhealthpolicies-api,shortName=chp,scope=Cluster
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +genclient:nonNamespaced
type ClusterHealthPolicy struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Contains the information required by a cluster health policy evaluator/contoller
	// to operate.
	Spec ClusterHealthPolicySpec `json:"spec,omitempty"`

	// List of ClusterHealthPolicyStatus.
	Status ClusterHealthPoliciesStatus `json:"status,omitempty"`
}

type ClusterHealthPoliciesStatus struct {
	ClusterHealthPoliciesStatus []ClusterHealthPolicyStatus `json:"clusterHealthPoliciesStatus,omitempty"`

	OsNodesPresentMap map[OperatingSystem]bool `json:"OsNodesPresentMap,omitempty"`
}

type ClusterHealthPolicyStatus struct {
	// ClusterHealthPolicy Status defines the observed health state of the policy.
	HealthPolicyStatus HealthPolicyStatus `json:"healthPolicyStatus,omitempty"`
}

type HealthPolicyStatus struct {
	// Name of the health policy to which we are reporting status.
	PolicyName string `json:"policyName,omitempty"`

	// The health policy status.Enum type having Failure,PartialFailure and Healthy types.
	PolicyStatus PolicyStatus `json:"policyStatus,omitempty"`

	// More Information about the health policy status.
	Message string `json:"message,omitempty"`

	// The time when the status was computed
	// +kubebuilder:validation:Format=date-time
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`

	// The time when the status was changed from the previous state
	// +kubebuilder:validation:Format=date-time
	LastStatusChangeTime string `json:"lastStatusChangeTime,omitempty"`
}

// +kubebuilder:validation:Enum=Unhealthy;PartialHealthy;Healthy;NotApplicable
type PolicyStatus string

const (
	// Health policy status is in unhelathy state
	PolicyStatusUnhealthy PolicyStatus = "Unhealthy"

	// Health policy status is in partial healthy state
	PolicyStatusPartialHealthy PolicyStatus = "PartialHealthy"

	// Health policy status is in healthy state
	PolicyStatusHealthy PolicyStatus = "Healthy"

	PolicyStatusNotApplicable PolicyStatus = "NotApplicable"
)

type ClusterHealthPolicySpec struct {
	// List of all cluster health status policies
	Policies []PolicySpecObject `json:"policies,omitempty"`
}

// PolicySpecObject is an item of the array defining the properties of policies.
type PolicySpecObject struct {
	// Policy Spec object
	PolicySpec PolicySpec `json:"policySpec"`
}

type PolicySpec struct {
	// Name of the health policy
	// +kubebuilder:validation:Required
	PolicyName string `json:"policyName"`

	// Defines the type of resource this service health status is calculated with respect to
	SelectKind string `json:"selectKind,omitempty"`

	// Defines the name of resource this service health status is calculated with respect to. Match Labels ignored in this case.
	SelectName string `json:"selectName,omitempty"`

	// Defines the name of application this service health status is calculated with respect to.
	// This is matched with maifestwork app name label to check if policy is applicable.
	SelectApplicationName string `json:"selectApplicationName,omitempty"`

	// Defines the namespaces under which the type of resource this service health status is calculated with respect to
	SelectNamespace string `json:"selectNamespace,omitempty"`

	// Defines the Api version of the resource this service health is defined by.
	SelectApiVersion string `json:"selectApiVersion,omitempty"`

	// matchLabels is a map of {key,value} pairs. A single {key,value} in the
	// matchLabels map is equivalent to an element of matchExpressions,
	//  whose key field is "key", the operator is "In", and the values array
	// contains only "value". The requirements are ANDed.
	MatchLabels map[string]string `json:"matchLabels,omitempty"`

	// Defines the OS applicable for this policy, if empty it means applicable for all OS.
	OsApplicable []OperatingSystem `json:"osApplicable,omitempty"`

	// Defines the minimum available resource count below which the service can be declared unhealthy/failure.
	// If the value is an integer, then its considered as absolute number of resources that we need to declare it healthy.
	// If the string like 40% is given, then we calculate the minAvailable resource count as this ratio.
	// +kubebuilder:validation:XIntOrString
	MinAvailable intstr.IntOrString `json:"minAvailable,omitempty"`

	// Defines the minimum available resource count below which the service can be declared partial unhealthy/partial failure.
	// If the value is an integer, then its considered as absolute number of resources that we need to declare it healthy.
	// If the string like 40% is given, then we calculate the minAvailable resource count as this ratio.
	// +kubebuilder:validation:XIntOrString
	MinPartialAvailable intstr.IntOrString `json:"minPartialAvailable,omitempty"`

	// Defines the impact of the current health policy on the overall cluster health.
	// Enum=High;Medium;Low;None
	ClusterHealthImpact ClusterHealthImpact `json:"clusterHealthImpact,omitempty"`

	// Defines the capabilities that the current health policy is contributing to in the cluster.
	// Enum=Infra;Data;Monitoring;Health;Framework;Routing;Logs;Ingress;Discovery
	ComponentCapabilities []ComponentCapability `json:"componentCapabilities,omitempty"`
}

// +kubebuilder:validation:Enum=linux;windows
type OperatingSystem string

const (
	OperatingSystemLinux OperatingSystem = "linux"

	OperatingSystemWindows OperatingSystem = "windows"
)

// +kubebuilder:validation:Enum=High;Medium;Low;None
type ClusterHealthImpact string

const (
	// ClusterHealthImpactHigh represents that the impact of the policy on overall cluster health is high
	ClusterHealthImpactHigh ClusterHealthImpact = "High"
	// ClusterHealthImpactMedium represents that the impact of the policy on overall cluster health is Medium
	ClusterHealthImpactMedium ClusterHealthImpact = "Medium"
	// ClusterHealthImpactLow represents that the impact of the policy on overall cluster health is Low
	ClusterHealthImpactLow ClusterHealthImpact = "Low"
	// ClusterHealthImpactNone represents that the impact of the policy on overall cluster health is Negligible
	ClusterHealthImpactNone ClusterHealthImpact = "None"
)

// +kubebuilder:validation:Enum=Infra;Data;Monitoring;Health;Framework;Routing;Logs;Ingress;Discovery
type ComponentCapability string

const (
	// ComponentCapabilityInfra represents that the policy contributes to Infra capability of the cluster
	ComponentCapabilityInfra ComponentCapability = "Infra"
	// ComponentCapabilityData represents that the policy contributes to Data capability of the cluster
	ComponentCapabilityData ComponentCapability = "Data"
	// ComponentCapabilityMonitoring represents that the policy contributes to Monitoring capability of the cluster
	ComponentCapabilityMonitoring ComponentCapability = "Monitoring"
	// ComponentCapabilityHealth represents that the policy contributes to Health capability of the cluster
	ComponentCapabilityHealth ComponentCapability = "Health"
	// ComponentCapabilityFramework represents that the policy contributes to Framework capability of the cluster
	ComponentCapabilityFramework ComponentCapability = "Framework"
	// ComponentCapabilityRouting represents that the policy contributes to Routing capability of the cluster
	ComponentCapabilityRouting ComponentCapability = "Routing"
	// ComponentCapabilityLogs represents that the policy contributes to logs capability of the cluster
	ComponentCapabilityLogs ComponentCapability = "Logs"
	// ComponentCapabilityIngress represents that the policy contributes to ingress capability of the cluster
	ComponentCapabilityIngress ComponentCapability = "Ingress"
	// ComponentCapabilityDiscovery represents that the policy contributes to discovery capability of the cluster
	ComponentCapabilityDiscovery ComponentCapability = "Discovery"
)

// ClusterHealthReportList is a list of ClusterHealthReport resources
// +kubebuilder:object:root=true
type ClusterHealthReportList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ClusterHealthReport `json:"items"`
}

// ClusterHealthReport document defines the cluster healh state of the standard cluster.
// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=clusterhealthreports-api,shortName=chr,scope=Cluster
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +genclient:nonNamespaced
type ClusterHealthReport struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Contains the information required by a cluster health status reporter
	// to operate.
	Spec ClusterHealthPolicySpec `json:"spec,omitempty"`

	// ClusterHealthReportStatus defines the observed health state of the cluster.
	Status ClusterHealthPoliciesStatus `json:"status,omitempty"`
}
