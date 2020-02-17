package images

import (
	"testing"

	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	exutil "github.com/openshift/origin/test/extended/util"
	e2e "k8s.io/kubernetes/test/e2e/framework"

	clientimageregistryv1 "github.com/openshift/client-go/imageregistry/clientset/versioned/typed/imageregistry/v1"
)

const (
	RegistryOperatorDeploymentNamespace = "openshift-image-registry"
	RegistryOperatorDeploymentName      = "cluster-image-registry-operator"
	ImageRegistryName                   = "image-registry"
	ImageRegistryResourceName           = "cluster"
	ImageRegistryOperatorResourceName   = "image-registry"
)

//Make sure Registry Operator is Available
func EnsureRegistryOperatorStatusIsAvailable(oc *exutil.CLI) {
	err := oc.AsAdmin().WithoutNamespace().Run("describe").Args("co", ImageRegistryOperatorResourceName).Execute()
	o.Expect(err).NotTo(o.HaveOccurred())
	g.By("No error for Image Registry Operator")

	availablestatus, err := oc.AsAdmin().WithoutNamespace().Run("get").Args("co", ImageRegistryOperatorResourceName, "-o=jsonpath={range .status.conditions[0]}{.status}").Output()
	o.Expect(err).NotTo(o.HaveOccurred())
	progressingstatus, err := oc.AsAdmin().WithoutNamespace().Run("get").Args("co", ImageRegistryOperatorResourceName, "-o=jsonpath={range .status.conditions[1]}{.status}").Output()
	o.Expect(err).NotTo(o.HaveOccurred())
	degradestatus, err := oc.AsAdmin().WithoutNamespace().Run("get").Args("co", ImageRegistryOperatorResourceName, "-o=jsonpath={range .status.conditions[2]}{.status}").Output()
	o.Expect(err).NotTo(o.HaveOccurred())
	if availablestatus == "True" && progressingstatus == "False" && degradestatus == "False" {
		g.By("Image registry operator is available")
	}

}

func RegistryConfigClient(oc *exutil.CLI) clientimageregistryv1.ImageregistryV1Interface {
	return clientimageregistryv1.NewForConfigOrDie(oc.UserConfig())
}

//Configure Image Registry Operator
func ConfigureImageRegistryStorage(oc *exutil.CLI) {
	var t *testing.T
	config, err := RegistryConfigClient(oc).Configs().Get(
		ImageRegistryResourceName,
		metav1.GetOptions{},
	)
	if err != nil {
		t.Fatal(err)
	}
	var hasstorage string
	if config.Status.Storage.EmptyDir != nil {
		e2e.Logf("Image Registry is already using EmptyDir")
	} else {
		switch {
		case config.Status.Storage.S3 != nil:
			hasstorage = "s3"
		case config.Status.Storage.Swift != nil:
			hasstorage = "swift"
		case config.Status.Storage.GCS != nil:
			hasstorage = "GCS"
		case config.Status.Storage.Azure != nil:
			hasstorage = "azure"
		case config.Status.Storage.PVC != nil:
			hasstorage = "pvc"
		default:
			e2e.Logf("Image Registry is using unknown storage type")
		}
		err = oc.AsAdmin().WithoutNamespace().Run("patch").Args("configs.imageregistry.operator.openshift.io", ImageRegistryResourceName, "-p", `{"spec":{"storage":{"`+hasstorage+`":null,"emptyDir":{}}}}`, "--type=merge").Execute()
		o.Expect(err).NotTo(o.HaveOccurred())
		if err != nil {
			e2e.Logf("Image Registry is not using EmptyDir")
		}
	}
	return
}

func EnableRegistryPublicRoute(oc *exutil.CLI) {
	err := oc.AsAdmin().WithoutNamespace().Run("patch").Args("configs.imageregistry.operator.openshift.io", ImageRegistryResourceName, "-p", `{"spec":{"defaultRoute":true}}`, "--type=merge").Execute()
	o.Expect(err).NotTo(o.HaveOccurred())
	if err != nil {
		e2e.Logf("Default route for Image Registry is failed to be enabled")
	}
}

/*func ConfigureImageRegistryStorage() {
	return
}*/
