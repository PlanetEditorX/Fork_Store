# SPDX-License-Identifier: GPL-2.0+

# Device Tree targets for Amlogic Meson architecture
dtb-$(CONFIG_ARCH_MESON) += meson-g12b-dangbei-b3-pro.dtb

# Build targets
targets += $(dtb-y)

# Add any required device tree compiler flags here
DTC_FLAGS += -a 0x8

# Phony targets
PHONY += dtbs

# Device Tree build rule
dtbs: $(addprefix $(obj)/, $(dtb-y))
    @:

# Clean rule
clean-files := *.dtb *.dtbo *_HS
