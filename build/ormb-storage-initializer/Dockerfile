FROM debian:stretch

COPY bin/ormb-storage-initializer /usr/local/bin/ormb-storage-initializer

ENTRYPOINT ["ormb-storage-initializer", "pull-and-export"]
