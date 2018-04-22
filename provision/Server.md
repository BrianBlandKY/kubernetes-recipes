# Server Setup

- Standard Fedora Installation
- No Swap, not compatible with Kubernetes
- All storage drives (at install) are LVM/XFS and mounted to /vault

## Static IP Configuration (Example)

``` vi
/etc/sysconfig/network-scripts/ifcfg-en*
TYPE=Ethernet
PROXY_METHOD=none
BROWSER_ONLY=no
BOOTPROTO=none
DEFROUTE=yes
IPV4_FAILURE_FATAL=no
IPV6INIT=yes
IPV6_AUTOCONF=yes
IPV6_DEFROUTE=yes
IPV6_FAILURE_FATAL=no
IPV6_ADDR_GEN_MODE=stable-privacy
NAME=eno1
UUID=6bfce6ad-7583-454a-9f70-746912f71d58
DEVICE=eno1
ONBOOT=yes
ETHTOOL_OPTS="autoneg on speed 1000 duplex full"
IPADDR=10.0.0.10
PREFIX=24
GATEWAY=10.0.0.1
DNS1=10.0.0.1
IPV6_PRIVACY=no
```

## SSHD Config for SSH key authentication only

``` sh
sudo nano /etc/ssh/sshd_config
# PermitRootLogin no
# StrictModes yes
# PubkeyAuthentication yes
# AuthorizedKeysFile %h/.ssh/authorized_keys
# PasswordAuthentication no

mkdir ~/.ssh
echo "ssh-rsa PUBLIC_KEY" >> ~/.ssh/authorized_keys
chmod -R 700 ~/.ssh
reboot
```

## Kubernetes Installation

``` sh
# Add Yum Repo
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOF

# Disable SELinux
setenforce 0

# At this time, we could only get it working with docker 17.03.02
cd /tmp
wget https://download.docker.com/linux/centos/7/x86_64/stable/Packages/docker-ce-selinux-17.03.2.ce-1.el7.centos.noarch.rpm

wget https://download.docker.com/linux/centos/7/x86_64/stable/Packages/docker-ce-17.03.2.ce-1.el7.centos.x86_64.rpm

yum install -y docker-ce-selinux-17.03.2.ce-1.el7.centos.noarch.rpm
yum install -y docker-ce-17.03.2.ce-1.el7.centos.x86_64.rpm

# Confirm docker is running
systemctl enable docker && systemctl start docker

# Ensure Docker compatibility
cat << EOF > /etc/docker/daemon.json
{
  "exec-opts": ["native.cgroupdriver=systemd"]
}
EOF

# Install Kubernetes
yum install -y kubelet kubeadm kubectl

# IPTables Fix
sysctl net.bridge.bridge-nf-call-iptables=1
sysctl --system

# Stop and Disable Firewall
systemctl stop firewalld && systemctl disable firewalld

# Restart Services
systemctl enable docker && systemctl start docker
systemctl enable kubelet && systemctl start kubelet
```

## Master Node

``` sh
# Create Master Node
kubeadm init —pod-network-cidr=10.244.0.0/16

# Save Token

# Setup Normal User (switch to normal user)
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# Install Pod Network (Flannel)
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/v0.9.1/Documentation/kube-flannel.yml

# Untaint Master Node
kubectl taint nodes --all node-role.kubernetes.io/master-

# NOTE: Confirm kube-dns pods are running before doing anything else.
# Otherwise everything will fail with cni errors
```

## Slave Node

``` sh
# Join Master
kubeadm join 10.0.0.10:6443 --token MASTER_TOKEN
```

Kubernetes Notes

``` sh
systemctl status docker     # Check docker service status
systemctl status kubelet    # Check kubelet service status
```

View service logs

``` sh
journalctl -xeu kubelet
journalctl -xeu docker
```