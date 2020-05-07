# Ansible deployment
Requires one, or two, Ubuntu 16.04 (virtual) machine(s) with root access and a local install of [Ansible 2.6](https://docs.ansible.com/ansible/2.6/installation_guide/intro_installation.html).

## Hosts configuration
In order to use Ansible deployment, `backend` and `frontend` hosts need to be configured in your [inventory](https://docs.ansible.com/ansible/2.6/user_guide/intro_inventory.html). For example:
```toml
[frontend]
short_name my-canoncial-host.name.com

[backend]
short_name_2 other-or-same-host.name.com
```

You may refer to the hosts file by setting `ANSIBLE_INVENTORY=<hosts_file>` in your environment.

## Amazon S3 credentials
The `snapshots` role expects the Amazon credentials to be configured in `vault/secrets.yml`. You can might either disable/skip the `snapshots` role or set up a bucket and authentication and configure the credentials based on `secrets_example.yml`, which you can encrypt with `ansible-vault encrypt`.

See:
- https://www.elastic.co/guide/en/elasticsearch/plugins/6.5/repository-s3-client.html
- https://github.com/harobed/ansible-vault-tutorial

## Bootstrappeing
The `bootstrap` playbook will make sure `sudo` and `Ansible` are available on the machine, assuming initial SSH root access. It creates a remote SSH user with pubkey authentication and sudo rights on the server(s) and completely disables password login.

After bootstrapping, the `site` playbook will do common configuration, including the backend and the frontend.

Then, using , execute:
```bash
$ ansible-playbook bootstrap.yml --user root --ask-pass
$ ansible-playbook site.yml
```

## Frontend deploy

### Staging vs. production certificates
By default, staging certificates are requested from LetsEncrypt. Once the above process proceeds, the variable `certbot_test` should be set to `false`. The best way to do this, is to define it in the inventory file, for example:
```toml
[frontend]
short_name my-canoncial-host.name.com certbot_test=false
```

### Deploying frontend updates
By default, the `v2` branch of the [frontend repository] is deployed, using the following command (the `-t frontend` makes sure that only the actual frontend code is deployed, rather than the entire frontend server setup):

```bash
$ ansible-playbook frontend.yml -t frontend
```
