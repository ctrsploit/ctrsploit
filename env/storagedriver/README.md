# env/storage-driver

storage-driver

存储驱动

```shell
root@913f6ea27a11:~/app# ./bin/latest/ctrsploit_linux_amd64 env storage-driver
===========Storage Driver===========
Type:                   overlay 
[Y]  Enabled    
[Y]  Used       
The number of graph driver mounted      # equal to the number of containers
20
The host path of container's rootfs     
/var/lib/docker/overlay2/6046653621c4e2dc878b50b9e17cf14fc8cb4b397e9d649d64c07c385ffdb1a5/merged

root@913f6ea27a11:~/app# ./bin/latest/ctrsploit_linux_amd64 --colorful env storage-driver
===========Storage Driver===========
Type:                   overlay 
✔  Enabled      
✔  Used 
The number of graph driver mounted      # equal to the number of containers
20
The host path of container's rootfs     
/var/lib/docker/overlay2/6046653621c4e2dc878b50b9e17cf14fc8cb4b397e9d649d64c07c385ffdb1a5/merged


root@913f6ea27a11:~/app# ./bin/latest/ctrsploit_linux_amd64 --json env storage-driver | jq
{
  "storage_driver": {
    "type": 1,
    "enabled": true,
    "used": true,
    "number": 20,
    "rootfs": "/var/lib/docker/overlay2/6046653621c4e2dc878b50b9e17cf14fc8cb4b397e9d649d64c07c385ffdb1a5/merged"
  }
}
```