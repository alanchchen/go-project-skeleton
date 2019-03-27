DEP_INSTALL_URL := https://raw.githubusercontent.com/golang/dep/v0.5.1/install.sh

$(HOSTBIN_DIR)/dep:
	$(Q)curl -s $(DEP_INSTALL_URL) | INSTALL_DIRECTORY=$(dir $@) sh