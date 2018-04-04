# Git-Sync

Sync Git repo to directory

## Install 

```
go get github.com/bborbe/git-sync
```

## Usage

With params

```
git-sync \
-logtostderr \
-v=2 \
-repo https://github.com/bborbe/git-sync.git \
-dest /tmp/git-sync
```

With username and password

```
git-sync \
-logtostderr \
-v=2 \
-repo https://github.com/bborbe/git-sync.git \
-dest /tmp/git-sync \
-username gituser \
-password gitpassword
```

With env

```
GIT_SYNC_USERNAME=gituser \
GIT_SYNC_PASSWORD=gitpassword \
GIT_SYNC_REPO=https://github.com/bborbe/git-sync.git \
GIT_SYNC_DEST=/tmp/git-sync \
git-sync \
-logtostderr \
-v=2
```

## Docker

```
mkdir -p /tmp/git-sync
docker run \
-v /tmp/git-sync:/git \
-e GIT_SYNC_DEST=/git \
-e GIT_SYNC_REPO=https://github.com/bborbe/git-sync.git \
bborbe/git-sync:1.1.6 \
-logtostderr \
-v=2
```
