package collector

import (
	"context"

	"github.com/rs/zerolog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type RBACData struct {
	ClusterRoles        interface{} `json:"clusterRoles"`
	ClusterRoleBindings interface{} `json:"clusterRoleBindings"`
	Roles               interface{} `json:"roles"`
	RoleBindings        interface{} `json:"roleBindings"`
}

type Collector struct {
	Clientset kubernetes.Interface
	Logger    zerolog.Logger
}

func New(clientset kubernetes.Interface, logger zerolog.Logger) *Collector {
	return &Collector{
		Clientset: clientset,
		Logger:    logger,
	}
}

func (c *Collector) CollectAndLog(ctx context.Context) error {
	listOptions := metav1.ListOptions{}

	clusterRoles, err := c.Clientset.RbacV1().ClusterRoles().List(ctx, listOptions)
	if err != nil {
		return err
	}

	clusterRoleBindings, err := c.Clientset.RbacV1().ClusterRoleBindings().List(ctx, listOptions)
	if err != nil {
		return err
	}

	roles, err := c.Clientset.RbacV1().Roles("").List(ctx, listOptions)
	if err != nil {
		return err
	}

	roleBindings, err := c.Clientset.RbacV1().RoleBindings("").List(ctx, listOptions)
	if err != nil {
		return err
	}

	payload := RBACData{
		ClusterRoles:        clusterRoles.Items,
		ClusterRoleBindings: clusterRoleBindings.Items,
		Roles:               roles.Items,
		RoleBindings:        roleBindings.Items,
	}

	c.Logger.Info().Interface("data", payload).Msg("Successfully collected RBAC data")

	return nil
}
