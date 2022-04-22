# /bin/sh
# set -x
echo "Start to build ctrsploit binary files..."
echo "The project path is \"${PROJECTPATH}\" in the docker"
mkdir -p $PROJECTPATH/build/bin/release
cd $PROJECTPATH/build/bin/release

### you can use gox --arch-list to get the supported OS/Arch list
# gox --arch-list

# by default build linux/amd64 and linux/arm64 OS/Arch binary
# build ctrsploit
CGO_ENABLED=0 gox -cgo=0 -osarch="linux/amd64" $PROJECTPATH/cmd/ctrsploit/...
CGO_ENABLED=0 gox -cgo=0 -osarch="linux/arm64" $PROJECTPATH/cmd/ctrsploit/...

# build checksec
CGO_ENABLED=0 gox -cgo=0 -osarch="linux/amd64" $PROJECTPATH/cmd/checksec/...
CGO_ENABLED=0 gox -cgo=0 -osarch="linux/arm64" $PROJECTPATH/cmd/checksec/...

# compress the binary files by upx
upx $PROJECTPATH/build/bin/release/*

echo "Done !"
echo "You can find the binary files in ${PROJECTPATH}/build/bin/release/ in the docker\n or in the project folder ./build/bin/release"

ls -alh

# Done
