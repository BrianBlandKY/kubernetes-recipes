### Plex Setup Notes

Using a claim will automatically discover the plex server on creation. However, these claims have expiration dates so if the system were to restart at a later date then it may not create the container successfully. (Theory)


See provisioning notes for NFS setup. NFS is managed by the node, not through kubernetes. I would like for this to change at some point. 