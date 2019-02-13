package test

import (
	"github.com/kubesphere/porter/pkg/bgp/routes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BGP routes test", func() {
	Context("Reconcile Routes", func() {
		FIt("Should correctly add/delete all routes", func() {
			ip := "100.100.100.100"
			nexthops := []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"}
			add, delete, err := routes.ReconcileRoutes(ip, nexthops)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(add).Should(HaveLen(3))
			Expect(delete).Should(HaveLen(0))
			err = routes.AddMultiRoutes(ip, 32, nexthops)
			Expect(err).ShouldNot(HaveOccurred())
			add, delete, err = routes.ReconcileRoutes(ip, nexthops)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(add).Should(HaveLen(0))
			Expect(delete).Should(HaveLen(0))
		})
	})
})
