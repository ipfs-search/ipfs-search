# Ansible deployment
Requires one, or two, Ubuntu 16.04 (virtual) machine(s) with root access and a local install of [Ansible 2.6](https://docs.ansible.com/ansible/2.6/installation_guide/intro_installation.html).

## Hosts configuration
In order to use Ansible deployment, `backend` and `frontend` hosts need to be configured in your [inventory](https://docs.ansible.com/ansible/2.6/user_guide/intro_inventory.html). For example:
```
[frontend]
my-canoncial-host.name.com

[backend]
other-or-same-host.name.com
```

You may refer to the hosts file by setting `ANSIBLE_INVENTORY=<hosts_file>` in your environment.

## Bootstrappeing
The `bootstrap` playbook will make sure `sudo` and `Ansible` are available on the machine, assuming initial SSH root access. It creates a remote SSH user with pubkey authentication and sudo rights on the server(s) and completely disables password login.

After bootstrapping, the `site` playbook will do common configuration, including the backend and the frontend.

Then, using , execute:
```bash
$ ansible-playbook bootstrap.yml --user root --ask-pass
$ ansible-playbook site.yml
```
