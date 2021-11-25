package karpenter

// Create creates a Karpenter installer task and waits for it to finish.
func (i *Installer) Create() error {
	taskTree := NewTasksToInstallKarpenter(i.Config, i.StackManager, i.CTL.Provider.EC2(), i.KarpenterInstaller)
	return doTasks(taskTree)
}
