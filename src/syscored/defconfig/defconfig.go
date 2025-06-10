package defconfig

import (
	"sync"

	InitService "github.com/Continu-OS/syscored/src/syscored"
	defconfigarc "github.com/Continu-OS/syscored/src/syscored/defconfig/arc"
)

var localOptionsMap *sync.Map = new(sync.Map)

func InitWithBootArgs(bootArgs *InitService.BootloaderStartParameters) error {
	// MEMORY/EMC Features
	localOptionsMap.Store("memory.auto_resize_emc_enable", defconfigarc.MEMORY_AUTO_RESIZE_ON_EMC)
	localOptionsMap.Store("memory.auto_resize_enable", defconfigarc.MEMORY_AUTO_RESIZE_ON_EMC)

	// CONTAINER Features
	localOptionsMap.Store("container.full_oscontainer_enable", defconfigarc.CONTAINER_FULL_OS)
	localOptionsMap.Store("container.servicecontainer_enable", defconfigarc.CONTAINER_RUN_SERVICES_DEFAULT)

	// SECURITY/LOGIN Features
	localOptionsMap.Store("security.allow_root_login", defconfigarc.SECURITY_ALLOW_SUDO_ROOT_LOGIN)

	// VIRTUALIZATION Features
	localOptionsMap.Store("virt.enable", defconfigarc.VIRTUALIZATION_ENABLE_CONTOS)

	// USER Features
	localOptionsMap.Store("user.multi_user_support_enable", defconfigarc.USER_MULTI_SESSION_SUPPORT)

	// HARDWARE Features
	localOptionsMap.Store("hardware.arm_gpio_support_enable", defconfigarc.HARDWARE_ARM_GPIO_SUPPORT)

	// KERNEL Features
	localOptionsMap.Store("kernel.force_cent_fork_enable", defconfigarc.KERNEL_FORCE_CENTOS_FORK)

	// SYSTEM Features
	localOptionsMap.Store("system.static_login_promt", "")

	// Es ist kein Fehler aufgetreten
	return nil
}

func GetOption(name string) (interface{}, bool) {
	val, ok := localOptionsMap.Load(name)
	if !ok {
		return nil, false
	}
	return val, true
}

func GetBoolOption(name string) bool {
	val, isOk := GetOption(name)
	if !isOk {
		return false
	}
	conv, ok := val.(bool)
	if !ok {
		return false
	}
	return conv
}
