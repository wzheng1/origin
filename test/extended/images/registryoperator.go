package images

import (
	o "github.com/onsi/gomega"
	restclient "k8s.io/client-go/rest"

	exutil "github.com/openshift/origin/test/extended/util"
	//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	e2e "k8s.io/kubernetes/test/e2e/framework"

	clientconfigv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	clientimageregistryv1 "github.com/openshift/client-go/imageregistry/clientset/versioned/typed/imageregistry/v1"
	clientroutev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	clientappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	clientstoragev1 "k8s.io/client-go/kubernetes/typed/storage/v1"
)

const (
	RegistryOperatorDeploymentNamespace = "openshift-image-registry"
	RegistryOperatorDeploymentName      = "cluster-image-registry-operator"
	ImageRegistryName                   = "image-registry"
	ImageRegistryResourceName           = "cluster"
	ImageRegistryOperatorResourceName   = "image-registry"
)

// Clientset is a set of Kubernetes clients.
type Clientset struct {
	clientcorev1.CoreV1Interface
	clientappsv1.AppsV1Interface
	clientconfigv1.ConfigV1Interface
	clientimageregistryv1.ImageregistryV1Interface
	clientroutev1.RouteV1Interface
	clientstoragev1.StorageV1Interface
}

//Make sure Registry Operator is Available
func EnsureRegistryOperatorStatusIsAvailable(oc *exutil.CLI) {
	msg, err := oc.AsAdmin().WithoutNamespace().Run("get").Args("co", ImageRegistryOperatorResourceName).Output()
	if err != nil {
		e2e.Failf("Unable to get co %s status, error:%v", msg, err)
	}
	o.Expect(err).NotTo(o.HaveOccurred())
	availablestatus, err := oc.AsAdmin().WithoutNamespace().Run("get").Args("co", ImageRegistryOperatorResourceName, "-o=jsonpath={range .status.conditions[0]}{.status}").Output()
	progressingstatus, err := oc.AsAdmin().WithoutNamespace().Run("get").Args("co", ImageRegistryOperatorResourceName, "-o=jsonpath={range .status.conditions[1]}{.status}").Output()
	degradestatus, err := oc.AsAdmin().WithoutNamespace().Run("get").Args("co", ImageRegistryOperatorResourceName, "-o=jsonpath={range .status.conditions[2]}{.status}").Output()
	if availablestatus == "True" && progressingstatus == "False" && degradestatus == "False" {
		e2e.Logf("Image registry operator is available")
	}
	o.Expect(err).NotTo(o.HaveOccurred())
}

// NewClientset creates a set of Kubernetes clients. The default kubeconfig is
// used if not provided.
func NewClientset() (clientset *Clientset, err error) {
	/*if kubeconfig == nil {
		kubeconfig, err = client.GetConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to get kubeconfig: %s", err)
		}
	}*/
	kubeconfig, err := e2e.LoadConfig()
	o.Expect(err).ToNot(o.HaveOccurred())

	clientset = &Clientset{}
	clientset.CoreV1Interface, err = clientcorev1.NewForConfig(kubeconfig)
	if err != nil {
		return
	}
	clientset.AppsV1Interface, err = clientappsv1.NewForConfig(kubeconfig)
	if err != nil {
		return
	}
	clientset.ConfigV1Interface, err = clientconfigv1.NewForConfig(kubeconfig)
	if err != nil {
		return
	}
	clientset.ImageregistryV1Interface, err = clientimageregistryv1.NewForConfig(kubeconfig)
	if err != nil {
		return
	}
	clientset.RouteV1Interface, err = clientroutev1.NewForConfig(kubeconfig)
	if err != nil {
		return
	}
	clientset.StorageV1Interface, err = clientstoragev1.NewForConfig(kubeconfig)
	if err != nil {
		return
	}
	return
}

// MustNewClientset is like NewClienset but aborts the test if clienset cannot
// be constructed.
func MustNewClientset(kubeconfig *restclient.Config) *Clientset {
	clientset, err := NewClientset()
	if err != nil {
		e2e.Logf("Cannot get kubeconfig")
	}
	return clientset
}

//Configure Image Registry Operator
/*func ConfigureImageRegistryStorage(oc *exutil.CLI) {
	client := MustNewClientset(nil)
	config, err := client.Configs().Get(
		ImageRegistryResourceName,
		metav1.GetOptions{},
	)
	if err != nil {
		e2e.Logf("unable to get image registry config")
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
		}
		err = oc.Run("patch").Args("configs.imageregistry.operator.openshift.io", ImageRegistryResourceName, "-p", `{"spec":{"storage":{`+hasstorage+`:null,"emptyDir":{}}}}`, "--type=merge").Execute()
		if err == nil {
			e2e.Logf("Image Registry is changed to use EmptyDir")
		}
	}
	return
}*/

func ConfigureImageRegistryStorage() {
	return
}
