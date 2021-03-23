package apply

import "github.com/weaveworks/eksctl/pkg/cfn/manager"

func (r *Reconciler) SetIRSAManager(manager IRSAManager) {
	r.irsaManager = manager
}

func (r *Reconciler) SetStackManager(stackManager manager.StackManager) {
	r.stackManager = stackManager
}
