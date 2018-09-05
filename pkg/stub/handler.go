package stub

import (
	"context"
	"os/exec"
	"strings"

	"github.com/rthallisey/kubevirt-operator/pkg/apis/app/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch event.Object.(type) {
	case *v1alpha1.App:
		kv_yaml := "/etc/kubevirt/kubevirt.yaml"
		kcmd := "kubectl"

		// According to https://github.com/operator-framework/operator-sdk/issues/270
		// Deleted field will be removed from event
		if event.Deleted {
			// According to https://kubernetes.io/docs/concepts/workloads/controllers/garbage-collection/#background-cascading-deletion
			// we could rely on background cascading deletion for kubevirt dependent objects
			// if metadata.ownerReferences field is properly set on them.
			cmd := exec.Command(kcmd, "delete", "-f", kv_yaml, "--cascade=true")
			out, _ := cmd.CombinedOutput()
			logrus.Infof(string(out))
			return nil
		}

		// Create kubevirt manifest using the client
		cmd := exec.Command(kcmd, "create", "-f", kv_yaml)
		// Error is outputed in plain text in out
		out, _ := cmd.CombinedOutput()
		if strings.Contains(string(out), "Error from server (AlreadyExists)") {
			logrus.Debugf("Resources from kubevirt.yaml already exist")
		} else {
			logrus.Infof(string(out))
		}
	}
	return nil
}
