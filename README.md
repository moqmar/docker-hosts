# Docker Hosts
Update a `/etc/hosts`-style file to add docker containers to resolve locally on the host. The application will run until an error occurs and will watch for events from the docker daemon, keeping the file up to date.

**Usage:** `docker-hosts [file=/etc/hosts] [tld=docker]`

Each docker container will be added as `<container-name>.docker`, `<container-short-id>.docker` and `<container-long-id>.docker`.  
Which network will be used for containers with multiple networks is undefined.

**Tip:** if you want to use this on a dedicated DNS server instead of using the hosts file, I suggest using [CoreDNS](https://coredns.io/) with the [hosts plugin](https://coredns.io/plugins/hosts/).

## Build and installation

If [Go](https://golang.org/), [dep](https://golang.github.io/dep/) and [run](https://github.com/moqmar/run) are installed, Docker Hosts can be installed simply by running `git clone github.com/moqmar/docker-hosts && cd docker-hosts && dep ensure && run install`.

If you want to install the binary version, use the following instead:
```
wget -qO- https://github.com/moqmar/docker-hosts/releases/latest | grep -Eo '/moqmar/docker-hosts/releases/download/[^/]+/docker-hosts' | sudo wget --base=http://github.com/ -i - -O /usr/local/bin/docker-hosts
sudo chmod +x /usr/local/bin/docker-hosts
sudo wget -O /etc/systemd/system/docker-hosts.service https://raw.githubusercontent.com/moqmar/docker-hosts/master/docker-hosts.service
sudo systemctl enable docker-hosts.service
sudo systemctl start docker-hosts.service
```
