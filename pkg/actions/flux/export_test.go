package flux

func (ti *Installer) SetFluxClient(client InstallerClient) {
	ti.fluxClient = client
}
