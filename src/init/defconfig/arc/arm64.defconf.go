package defconfigarc

const (
	// MEMORY/EMC Features
	MEMORY_AUTO_RESIZE_ON_EMC bool = true // Prüft beim Starten, ob der Speicherplatz korrekt initialisiert ist, ggf. Resize (z.B. SD-Karten)

	// CONTAINER Features
	CONTAINER_RUN_SERVICES_DEFAULT bool = false // Startet Systemdienste standardmäßig in einem eigenen Container
	CONTAINER_FULL_OS              bool = true  // Startet alles nach dem Basissystem in Containern (z.B. für sichere Updates)

	// SECURITY/LOGIN Features
	SECURITY_ALLOW_SUDO_ROOT_LOGIN bool = false // Erlaubt Root-Login über SUDO

	// VIRTUALIZATION Features
	VIRTUALIZATION_ENABLE_CONTOS bool = false // Aktiviert Virtualisierung (KVM/QEMU)

	// USER Features
	USER_MULTI_SESSION_SUPPORT bool = false // Erlaubt mehrere Benutzer-Sessions auf dem Gerät

	// HARDWARE Features
	HARDWARE_ARM_GPIO_SUPPORT bool = false // Aktiviert Zugriff auf GPIO Pins (ARM)

	// KERNEL Features
	KERNEL_FORCE_CENTOS_FORK bool = true // Erzwingt speziellen (CentOS-basierten) Kernel
)
