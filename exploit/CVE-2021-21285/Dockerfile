FROM alpine as data
# https://stackoverflow.com/questions/52086641/override-a-volume-when-building-docker-image-from-another-docker-image
COPY v2 /var/lib/registry/docker/registry/v2
RUN sed -i s/sha256:b8dfde/sha256b8dfde/g /var/lib/registry/docker/registry/v2/blobs/sha256/1b/1b26826f602946860c279fce658f31050cff2c596583af237d971f4629b57792/data

FROM registry:2.7.1
COPY --from=data /var/lib/registry/docker/registry/v2 /var/lib/registry/docker/registry/v2
