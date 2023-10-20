package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/discovery/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"unsafe"
)

var (
	// TODO: move SchemeBuilder with zz_generated.deepcopy.go to k8s.io/api.
	// localSchemeBuilder and AddToScheme will stay in k8s.io/kubernetes.
	localSchemeBuilder = &SchemeBuilder
)

func Convert_v1beta1_Endpoint_To_v1_Endpoint(in *Endpoint, out *v1.Endpoint, s conversion.Scope) error {
	out.Addresses = *(*[]string)(unsafe.Pointer(&in.Addresses))
	if err := Convert_v1beta1_EndpointConditions_To_v1_EndpointConditions(&in.Conditions, &out.Conditions, s); err != nil {
		return err
	}
	out.Hostname = (*string)(unsafe.Pointer(in.Hostname))
	out.TargetRef = (*corev1.ObjectReference)(unsafe.Pointer(in.TargetRef))
	out.DeprecatedTopology = *(*map[string]string)(unsafe.Pointer(&in.Topology))
	out.NodeName = (*string)(unsafe.Pointer(in.NodeName))
	out.Hints = (*v1.EndpointHints)(unsafe.Pointer(in.Hints))
	return nil
}

func Convert_v1_Endpoint_To_v1beta1_Endpoint(in *v1.Endpoint, out *Endpoint, s conversion.Scope) error {
	out.Addresses = *(*[]string)(unsafe.Pointer(&in.Addresses))
	if err := Convert_v1_EndpointConditions_To_v1beta1_EndpointConditions(&in.Conditions, &out.Conditions, s); err != nil {
		return err
	}
	out.Hostname = (*string)(unsafe.Pointer(in.Hostname))
	out.TargetRef = (*corev1.ObjectReference)(unsafe.Pointer(in.TargetRef))
	out.Topology = *(*map[string]string)(unsafe.Pointer(&in.DeprecatedTopology))
	out.NodeName = (*string)(unsafe.Pointer(in.NodeName))
	// WARNING: in.Zone requires manual conversion: does not exist in peer-type
	out.Hints = (*EndpointHints)(unsafe.Pointer(in.Hints))
	return nil
}
