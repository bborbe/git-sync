# Git-Sync

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

## Continuous integration

[Jenkins](https://www.benjamin-borbe.de/jenkins/job/Go-Git-Sync/)
