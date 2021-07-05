/*
Copyright 2019 The OpenEBS Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tests

import (
	"flag"

	"os"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"github.com/openebs/dynamic-localpv-provisioner/tests/artifacts"
	ns "github.com/openebs/maya/pkg/kubernetes/namespace/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	kubeConfigPath                  string
	namespace                       = "localpv-ns"
	namespaceObj                    *corev1.Namespace
	err                             error
	LocalPVProvisionerLabelSelector = "openebs.io/component-name=openebs-localpv-provisioner"
)

func TestSource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test application deployment")
}

func init() {
	flag.StringVar(&kubeConfigPath, "kubeconfig", os.Getenv("KUBECONFIG"), "path to kubeconfig to invoke kubernetes API calls")
}

var ops *Operations

var _ = BeforeSuite(func() {

	ops = NewOperations(WithKubeConfigPath(kubeConfigPath))

	By("waiting for openebs-localpv-provisioner pod to come into running state")
	provPodCount := ops.GetPodRunningCountEventually(
		string(artifacts.OpenebsNamespace),
		string(artifacts.OpenEBSLocalPVProvisionerLabelSelector),
		1,
	)
	Expect(provPodCount).To(Equal(1))

	By("waiting for openebs-ndm-operator pod to come into running state")
	ndmPodCount := ops.GetPodRunningCountEventually(
		string(artifacts.OpenebsNamespace),
		string(artifacts.OpenEBSNDMOperatorLabelSelector),
		1,
	)
	Expect(ndmPodCount).To(Equal(1))

	By("building a namespace")
	namespaceObj, err = ns.NewBuilder().
		WithGenerateName(namespace).
		APIObject()
	Expect(err).ShouldNot(HaveOccurred(), "while building namespace {%s}", namespaceObj.GenerateName)

	By("creating above namespace")
	namespaceObj, err = ops.NSClient.Create(namespaceObj)
	Expect(err).To(BeNil(), "while creating namespace {%s}", namespaceObj.GenerateName)

})

var _ = AfterSuite(func() {

	By("deleting namespace")
	err = ops.NSClient.Delete(namespaceObj.Name, &metav1.DeleteOptions{})
	Expect(err).To(BeNil(), "while deleting namespace {%s}", namespaceObj.Name)

})

// Declarations for Ginkgo DSL
type Done ginkgo.Done
type Benchmarker ginkgo.Benchmarker

var GinkgoWriter = ginkgo.GinkgoWriter
var GinkgoRandomSeed = ginkgo.GinkgoRandomSeed
var GinkgoParallelNode = ginkgo.GinkgoParallelNode
var GinkgoT = ginkgo.GinkgoT
var CurrentGinkgoTestDescription = ginkgo.CurrentGinkgoTestDescription
var RunSpecs = ginkgo.RunSpecs
var RunSpecsWithDefaultAndCustomReporters = ginkgo.RunSpecsWithDefaultAndCustomReporters
var RunSpecsWithCustomReporters = ginkgo.RunSpecsWithCustomReporters
var Skip = ginkgo.Skip
var Fail = ginkgo.Fail
var GinkgoRecover = ginkgo.GinkgoRecover
var Describe = ginkgo.Describe
var FDescribe = ginkgo.FDescribe
var PDescribe = ginkgo.PDescribe
var XDescribe = ginkgo.XDescribe
var Context = ginkgo.Context
var FContext = ginkgo.FContext
var PContext = ginkgo.PContext
var XContext = ginkgo.XContext
var When = ginkgo.When
var FWhen = ginkgo.FWhen
var PWhen = ginkgo.PWhen
var XWhen = ginkgo.XWhen
var It = ginkgo.It
var FIt = ginkgo.FIt
var PIt = ginkgo.PIt
var XIt = ginkgo.XIt
var Specify = ginkgo.Specify
var FSpecify = ginkgo.FSpecify
var PSpecify = ginkgo.PSpecify
var XSpecify = ginkgo.XSpecify
var By = ginkgo.By
var Measure = ginkgo.Measure
var FMeasure = ginkgo.FMeasure
var PMeasure = ginkgo.PMeasure
var XMeasure = ginkgo.XMeasure
var BeforeSuite = ginkgo.BeforeSuite
var AfterSuite = ginkgo.AfterSuite
var SynchronizedBeforeSuite = ginkgo.SynchronizedBeforeSuite
var SynchronizedAfterSuite = ginkgo.SynchronizedAfterSuite
var BeforeEach = ginkgo.BeforeEach
var JustBeforeEach = ginkgo.JustBeforeEach
var JustAfterEach = ginkgo.JustAfterEach
var AfterEach = ginkgo.AfterEach

// Declarations for Gomega DSL
var RegisterFailHandler = gomega.RegisterFailHandler
var RegisterFailHandlerWithT = gomega.RegisterFailHandlerWithT
var RegisterTestingT = gomega.RegisterTestingT
var InterceptGomegaFailures = gomega.InterceptGomegaFailures
var Ω = gomega.Ω
var Expect = gomega.Expect
var ExpectWithOffset = gomega.ExpectWithOffset
var Eventually = gomega.Eventually
var EventuallyWithOffset = gomega.EventuallyWithOffset
var Consistently = gomega.Consistently
var ConsistentlyWithOffset = gomega.ConsistentlyWithOffset
var SetDefaultEventuallyTimeout = gomega.SetDefaultEventuallyTimeout
var SetDefaultEventuallyPollingInterval = gomega.SetDefaultEventuallyPollingInterval
var SetDefaultConsistentlyDuration = gomega.SetDefaultConsistentlyDuration
var SetDefaultConsistentlyPollingInterval = gomega.SetDefaultConsistentlyPollingInterval
var NewWithT = gomega.NewWithT
var NewGomegaWithT = gomega.NewGomegaWithT
var Default = gomega.Default

// Declarations for Gomega Matchers
var Equal = gomega.Equal
var BeEquivalentTo = gomega.BeEquivalentTo
var BeIdenticalTo = gomega.BeIdenticalTo
var BeNil = gomega.BeNil
var BeTrue = gomega.BeTrue
var BeFalse = gomega.BeFalse
var HaveOccurred = gomega.HaveOccurred
var Succeed = gomega.Succeed
var MatchError = gomega.MatchError
var BeClosed = gomega.BeClosed
var Receive = gomega.Receive
var BeSent = gomega.BeSent
var MatchRegexp = gomega.MatchRegexp
var ContainSubstring = gomega.ContainSubstring
var HavePrefix = gomega.HavePrefix
var HaveSuffix = gomega.HaveSuffix
var MatchJSON = gomega.MatchJSON
var MatchXML = gomega.MatchXML
var MatchYAML = gomega.MatchYAML
var BeEmpty = gomega.BeEmpty
var HaveLen = gomega.HaveLen
var HaveCap = gomega.HaveCap
var BeZero = gomega.BeZero
var ContainElement = gomega.ContainElement
var BeElementOf = gomega.BeElementOf
var ConsistOf = gomega.ConsistOf
var ContainElements = gomega.ContainElements
var HaveKey = gomega.HaveKey
var HaveKeyWithValue = gomega.HaveKeyWithValue
var BeNumerically = gomega.BeNumerically
var BeTemporally = gomega.BeTemporally
var BeAssignableToTypeOf = gomega.BeAssignableToTypeOf
var Panic = gomega.Panic
var PanicWith = gomega.PanicWith
var BeAnExistingFile = gomega.BeAnExistingFile
var BeARegularFile = gomega.BeARegularFile
var BeADirectory = gomega.BeADirectory
var HaveHTTPStatus = gomega.HaveHTTPStatus
var And = gomega.And
var SatisfyAll = gomega.SatisfyAll
var Or = gomega.Or
var SatisfyAny = gomega.SatisfyAny
var Not = gomega.Not
var WithTransform = gomega.WithTransform
var Satisfy = gomega.Satisfy
