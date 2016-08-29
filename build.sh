#!/bin/sh

SOURCEDIRECTORY="github.com/bborbe/git-sync"
VERSION="1.0.1-b${BUILD_NUMBER}"
NAME="git-sync"

################################################################################

echo "use workspace ${WORKSPACE}"

export GOROOT=/opt/go
export PATH=/opt/go2xunit/bin/:/opt/utils/bin/:/opt/aptly_utils/bin/:/opt/aptly/bin/:/opt/debian_utils/bin/:/opt/debian/bin/:$GOROOT/bin:$PATH
export GOPATH=${WORKSPACE}
export REPORT_DIR=${WORKSPACE}/test-reports
INSTALLS=`cd src && find $SOURCEDIRECTORY/bin -name "*.go" | dirof | unique`
DEB="${NAME}_${VERSION}.deb"
rm -rf $REPORT_DIR ${WORKSPACE}/*.deb ${WORKSPACE}/pkg
mkdir -p $REPORT_DIR
PACKAGES=`cd src && find $SOURCEDIRECTORY -name "*_test.go" | dirof | unique`
FAILED=false
for PACKAGE in $PACKAGES
do
  XML=$REPORT_DIR/`pkg2xmlname $PACKAGE`
  OUT=$XML.out
  go test -i $PACKAGE
  go test -v $PACKAGE | tee $OUT
  cat $OUT
  go2xunit -fail=true -input $OUT -output $XML
  rc=$?
  if [ $rc != 0 ]
  then
    echo "Tests failed for package $PACKAGE"
    FAILED=true
  fi
done

if $FAILED
then
  echo "Tests failed => skip install"
  exit 1
else
  echo "Tests success"
fi

echo "Tests completed, install to ${GOPATH}"

go install $INSTALLS

echo "Install completed, create debian package"

create_debian_package \
-logtostderr \
-v=2 \
-version=$VERSION \
-config=src/$SOURCEDIRECTORY/create_debian_package_config.json || exit 1

echo "Create debian package completed, start upload to aptly"

aptly_upload \
-logtostderr \
-v=2 \
-url=http://aptly-api.aptly.svc.cluster.local:3845 \
-file=$DEB \
-repo=unstable || exit 1

echo "Upload completed"
