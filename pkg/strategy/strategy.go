package strategy

import (
	"fmt"
	"sort"

	"github.com/kubesphere/porter/pkg/apis/network/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

type EIPSelectStrategyType string

const (
	DefaultStrategy EIPSelectStrategyType = "Default"
)

type EIPSelectStrategy interface {
	Select(*corev1.Service, *v1alpha1.EIPList) (*v1alpha1.EIP, error)
}

func GetStrategy(name EIPSelectStrategyType) (EIPSelectStrategy, error) {
	switch name {
	case DefaultStrategy:
		return defaultStrategy{}, nil
	default:
		return nil, fmt.Errorf("Strategy %s not found", name)
	}
}

type defaultStrategy struct{}

func (defaultStrategy) Select(serv *corev1.Service, eips *v1alpha1.EIPList) (*v1alpha1.EIP, error) {
	if len(eips.Items) == 0 {
		return nil, fmt.Errorf("Not enough ips to select")
	}
	for _, ip := range eips.Items {
		if !ip.Spec.Disable {
			if len(ip.Status.PortsUsage) == 0 {
				return &ip, nil
			}
			chosen := false
			ports := make([]int32, 0, len(ip.Status.PortsUsage))
			for key := range ip.Status.PortsUsage {
				ports = append(ports, key)
			}
			for _, port := range serv.Spec.Ports {
				index := sort.Search(len(ports), func(i int) bool {
					return ports[i] >= port.Port
				})
				if ports[index] == port.Port {
					chosen = false
					break
				}
				chosen = true
			}
			if chosen {
				return &ip, nil
			}
		}
	}
	return nil, fmt.Errorf("No suitable ip has empty ports for service %s", serv.Name)
}
