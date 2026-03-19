import subprocess
import sys

from pyinfra import host
from pyinfra.operations import apt, files, systemd, server

DEB_PATH = "godeps/target/deb/xrayvpn_0.1.0_amd64.deb"
REMOTE_DEB_PATH = "/tmp/xrayvpn_0.1.0_amd64.deb"

subprocess.run(["make", "deb"], cwd="godeps", check=True)

apt.update(name="Update apt cache", cache_time=3600, _env={"DEBIAN_FRONTEND": "noninteractive"})
apt.packages(name="Install deps", packages=["nftables", "dnsmasq", "hostapd"], present=True, _env={"DEBIAN_FRONTEND": "noninteractive"})

files.put(name="Upload xrayvpn .deb package", src=DEB_PATH, dest=REMOTE_DEB_PATH, mode="0644")
server.shell(name="Install xrayvpn .deb package", commands=[f"dpkg -i {REMOTE_DEB_PATH}"], _env={"DEBIAN_FRONTEND": "noninteractive"})


files.template(name="Deploy /etc/network/interfaces", src="templates/interfaces.j2", dest="/etc/network/interfaces", mode="0644", user="root", group="root")
files.template(name="Deploy /etc/nftables.conf", src="templates/nftables.conf.j2", dest="/etc/nftables.conf", mode="0644", user="root", group="root")
files.template(name="Deploy /etc/dnsmasq.conf", src="templates/dnsmasq.conf.j2", dest="/etc/dnsmasq.conf", mode="0644", user="root", group="root")
files.template(name="Deploy /etc/hostapd/hostapd.conf", src="templates/hostapd.conf.j2", dest="/etc/hostapd/hostapd.conf", mode="0600", user="root", group="root")
files.template(name="Deploy /etc/ssh/sshd_config", src="templates/sshd_config.j2", dest="/etc/ssh/sshd_config", mode="0644", user="root", group="root")

files.directory(name="Create nftables.service.d directory", path="/etc/systemd/system/nftables.service.d", mode="0755", user="root", group="root")
files.put(
    name="Deploy nftables service override",
    src="templates/nftables-override.conf",
    dest="/etc/systemd/system/nftables.service.d/override.conf",
    mode="0644",
    user="root",
    group="root",
)

server.shell(name="Add gleb to xrayvpn group", commands=["usermod -aG xrayvpn gleb"])

for svc in ["nftables", "dnsmasq", "hostapd", "xrayvpnd"]:
    systemd.service(name=f"Enable and start {svc}", service=svc, running=True, enabled=True, restarted=True, daemon_reload=True)
