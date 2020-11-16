package kubemon

import (
	"github.com/Dynatrace/dynatrace-operator/pkg/apis/dynatrace/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func buildResources(instance *v1alpha1.DynaKube) corev1.ResourceRequirements {
	limits := buildResourceLimits(instance)
	requests := buildResourceRequests(instance, limits)

	return corev1.ResourceRequirements{
		Limits:   limits,
		Requests: requests,
	}
}

func buildResourceRequests(instance *v1alpha1.DynaKube, limits corev1.ResourceList) corev1.ResourceList {
	cpuMin := resource.MustParse(ResourceCPUMinimum)
	cpuRequest, hasCPURequest := instance.Spec.KubernetesMonitoringSpec.Resources.Requests[corev1.ResourceCPU]
	if !hasCPURequest {
		cpuRequest = cpuMin
	}

	memoryMin := resource.MustParse(ResourceMemoryMinimum)
	memoryRequest, hasMemoryRequest := instance.Spec.KubernetesMonitoringSpec.Resources.Requests[corev1.ResourceCPU]
	if !hasMemoryRequest {
		memoryMin = memoryRequest
	}

	return corev1.ResourceList{
		corev1.ResourceCPU:    getMinResource(getMaxResource(cpuMin, cpuRequest), limits[corev1.ResourceCPU]),
		corev1.ResourceMemory: getMinResource(getMaxResource(memoryMin, memoryRequest), limits[corev1.ResourceMemory]),
	}
}

func buildResourceLimits(instance *v1alpha1.DynaKube) corev1.ResourceList {
	cpuMax := resource.MustParse(ResourceCPUMaximum)
	cpuLimit, hasCPULimit := instance.Spec.KubernetesMonitoringSpec.Resources.Limits[corev1.ResourceCPU]
	if !hasCPULimit {
		cpuLimit = cpuMax
	}

	memoryMax := resource.MustParse(ResourceMemoryMaximum)
	memoryLimit, hasMemoryLimit := instance.Spec.KubernetesMonitoringSpec.Resources.Limits[corev1.ResourceMemory]
	if !hasMemoryLimit {
		memoryLimit = memoryMax
	}

	return corev1.ResourceList{
		corev1.ResourceCPU:    getMinResource(cpuLimit, cpuMax),
		corev1.ResourceMemory: getMinResource(memoryLimit, memoryMax),
	}
}

func getMinResource(a resource.Quantity, b resource.Quantity) resource.Quantity {
	if isSmallerThan(a, b) {
		return a
	}
	return b
}

func getMaxResource(a resource.Quantity, b resource.Quantity) resource.Quantity {
	if isSmallerThan(a, b) {
		return b
	}
	return a
}

func isSmallerThan(a resource.Quantity, reference resource.Quantity) bool {
	return a.Cmp(reference) > 0
}
