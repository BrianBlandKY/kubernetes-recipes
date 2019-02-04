# Setup NFS


```
# NFS Tools
$ yum install nfs-utils nfs-utils-lib


# Auto mount
$ vi /etc/fstab

# Add this line
# 192.168.0.100:/nfsshare /mnt  nfs defaults 0 0

```