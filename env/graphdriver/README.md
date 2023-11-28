# env/graphdriver

graphdriver

存储驱动

```
root@b93597d64108:/# ./env_linux_amd64 --colorful graphdriver
===========Overlay===========
✔  Enabled      
✔  Used 
The number of graph driver mounted      # equal to the number of containers
7
The host path of container's rootfs     
/var/lib/docker/overlay2/2f660313765143ee9253b4a93a0ca665ed9288fc3e4aca927426f5ac61ca77bb/merged
===========DeviceMapper===========
✔  Enabled      
✘  Used
```

```
root@4fe779fc104c:~/ctrsploit/bin/release# ./env_linux_amd64 --json g | jq
{
  "devicemapper": {
    "name": {
      "name": "DeviceMapper"
    },
    "enabled": {
      "name": "Enabled",
      "description": "",
      "result": true
    },
    "used": {
      "name": "Used",
      "description": "",
      "result": false
    },
    "number": {
      "name": "",
      "description": "",
      "result": ""
    },
    "host_path": {
      "name": "",
      "description": "",
      "result": ""
    }
  },
  "overlay": {
    "name": {
      "name": "Overlay"
    },
    "enabled": {
      "name": "Enabled",
      "description": "",
      "result": true
    },
    "used": {
      "name": "Used",
      "description": "",
      "result": true
    },
    "number": {
      "name": "The number of graph driver mounted",
      "description": "equal to the number of containers",
      "result": "11"
    },
    "host_path": {
      "name": "The host path of container's rootfs",
      "description": "",
      "result": "/var/lib/docker/overlay2/488d1aede2705505ad41e04aba5a3b42c5ce3905caabfe7da65de1c70178644b/merged"
    }
  }
}
```