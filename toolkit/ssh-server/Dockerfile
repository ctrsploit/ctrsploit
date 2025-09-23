FROM ubuntu
ADD --chmod=755 https://github.com/ctrsploit/ctrsploit/releases/latest/download/ctrsploit_linux_amd64 /usr/bin/ctrsploit
RUN apt update && apt install -y \
    openssh-server
RUN mkdir /var/run/sshd && \
    echo 'root:ced53a64-9861-11f0-9d84-e89c258aa5af' | chpasswd && \
    echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config && \
    echo 'PasswordAuthentication yes' >> /etc/ssh/sshd_config && \
    mkdir -p /root/.ssh && \
    chmod 700 /root/.ssh
EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]
