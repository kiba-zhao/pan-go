package remoteitem

type VFSFUSERemoteFileRuntime interface {
	RemoteFileInfoService() *RemoteFileInfoService
	RemoteFileStreamService() *RemoteFileStreamService
}
