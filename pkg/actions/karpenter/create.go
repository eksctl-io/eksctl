package karpenter

// Create creates a Karpenter installer task and waits for it to finish.
func (i *Installer) Create() error {
	taskTree := NewTasksToInstallKarpenter(i.cfg, i.stackManager, i.ctl.Provider.EC2(), i.karpenterInstaller)
	err := doTasks(taskTree)
	return err
}
