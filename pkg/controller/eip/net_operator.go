package eip

import (
	"context"
	"strings"

	"github.com/kubesphere/porter/pkg/controller/eip/nettool"

	networkv1alpha1 "github.com/kubesphere/porter/pkg/apis/network/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcileEIP) AddRoute(instance *networkv1alpha1.EIP) error {
	checkedService := make(map[string]bool)
	if len(instance.Status.PortsUsage) != 0 {
		for _, value := range instance.Status.PortsUsage {
			if _, ok := checkedService[value]; !ok {
				checkedService[value] = true
			} else {
				continue
			}
			splits := strings.Split(value, "/")
			service := &corev1.Service{}
			err := r.Get(context.TODO(), types.NamespacedName{Namespace: splits[0], Name: splits[1]}, service)
			if err != nil {
				return err
			}
			clusterIP := service.Spec.ClusterIP
			route := nettool.NewEIPRoute(instance.Spec.Address, clusterIP, 32)
			yes, err := route.IsExist()
			if err != nil {
				return err
			}
			if !yes {
				if err := route.Add(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *ReconcileEIP) DelRoute(instance *networkv1alpha1.EIP) error {
	checkedService := make(map[string]bool)
	if len(instance.Status.PortsUsage) != 0 {
		for _, value := range instance.Status.PortsUsage {
			if _, ok := checkedService[value]; !ok {
				checkedService[value] = true
			} else {
				continue
			}
			splits := strings.Split(value, "/")
			service := &corev1.Service{}
			err := r.Get(context.TODO(), types.NamespacedName{Namespace: splits[0], Name: splits[1]}, service)
			if err != nil {
				return err
			}
			clusterIP := service.Spec.ClusterIP
			route := nettool.NewEIPRoute(instance.Spec.Address, clusterIP, 32)
			yes, err := route.IsExist()
			if err != nil {
				return err
			}
			if yes {
				if err := route.Delete(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
