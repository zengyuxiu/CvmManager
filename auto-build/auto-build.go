package auto_build

import (
	"CvmManager/config"
	docker_api "CvmManager/docker-api"
	lxc_api "CvmManager/lxc-api"
)

func AutoBuild() {
	networkConfig := config.GetConfig()
	networkConfig.OspfdConfig()
	lxc_api.BridgeCreate(networkConfig)
	docker_api.NetworkCreate(networkConfig)
	lxc_api.CreateRoute(networkConfig)
	lxc_api.BridgeAttach(networkConfig)
	lxc_api.InstanceOspfConfig(networkConfig)
	lxc_api.RestartRoute(networkConfig)
}
