// Code generated by helmit-generate. DO NOT EDIT.

package v1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	"time"
)

var ClusterRoleKind = resource.Kind{
	Group:   "rbac.authorization.k8s.io",
	Version: "v1",
	Kind:    "ClusterRole",
	Scoped:  false,
}

var ClusterRoleResource = resource.Type{
	Kind: ClusterRoleKind,
	Name: "clusterroles",
}

func NewClusterRole(clusterRole *rbacv1.ClusterRole, client resource.Client) *ClusterRole {
	return &ClusterRole{
		Resource: resource.NewResource(clusterRole.ObjectMeta, ClusterRoleKind, client),
		Object:   clusterRole,
	}
}

type ClusterRole struct {
	*resource.Resource
	Object *rbacv1.ClusterRole
}

func (r *ClusterRole) Delete() error {
	client, err := kubernetes.NewForConfig(r.Config())
	if err != nil {
		return err
	}
	return client.RbacV1().
		RESTClient().
		Delete().
		NamespaceIfScoped(r.Namespace, ClusterRoleKind.Scoped).
		Resource(ClusterRoleResource.Name).
		Name(r.Name).
		VersionedParams(&metav1.DeleteOptions{}, metav1.ParameterCodec).
		Timeout(time.Minute).
		Do().
		Error()
}