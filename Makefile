ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

.DEFAULT_GOAL := install
.PHONY: dependencies install

dependencies:
	sudo zypper refresh
	sudo zypper install -y ansible

install: dependencies
	cd "$(ROOT_DIR)" && sudo env ANSIBLE_CONFIG="$(ROOT_DIR)/ansible.cfg" ANSIBLE_ROLES_PATH="$(ROOT_DIR)/roles" \
		ansible-playbook --become --inventory "$(ROOT_DIR)/inventory/hosts.yml" playbooks/site.yml $(ARGS)
