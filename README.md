# Docker Hosts
Update a `/etc/hosts`-style file to add docker containers to resolve locally on the host. The application will run until an error occurs and will watch for events from the docker daemon, keeping the file up to date.

It will not remove existing host definitions from a file, so it can be safely used with `/etc/hosts` directly.

**Usage:** `docker-hosts [/etc/hosts [docker]]`

The first argument specifies the target file, the second argument specifies the TLD (`docker` by default, you might want to use an actual domain because someone registering the `docker` TLD might lead to some trouble).

Each docker container will be added as `<container-name>.docker`, `<container-short-id>.docker`.  
Which network will be used for containers with multiple networks is undefined behaviour.

You can also set additional names (seprated by `;`) via the label `docker-hosts.hosts`.

**Tip:** if you want to use this on a dedicated DNS server instead of using the hosts file directly, I suggest using [CoreDNS](https://coredns.io/) with the [hosts plugin](https://coredns.io/plugins/hosts/).

## Build and installation

If [Go](https://golang.org/), is installed, Docker Hosts can be installed simply by running `git clone https://github.com/moqmar/docker-hosts && cd docker-hosts && go install .`.

If you want to install the binary version (for 64-bit systems), use the following instead:
```
wget -qO- https://github.com/moqmar/docker-hosts/releases/latest | grep -Eo '/moqmar/docker-hosts/releases/download/[^/]+/docker-hosts' | sudo wget --base=http://github.com/ -i - -O /usr/local/bin/docker-hosts
sudo chmod +x /usr/local/bin/docker-hosts
sudo wget -O /etc/systemd/system/docker-hosts.service https://raw.githubusercontent.com/moqmar/docker-hosts/master/docker-hosts.service
sudo systemctl enable docker-hosts.service
sudo systemctl start docker-hosts.service
```
