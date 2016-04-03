# Git-Sync

## Usage
 
```
GIT_SYNC_USERNAME=gituser GIT_SYNC_PASSWORD=gitpassword GIT_SYNC_REPO=https://bitbucket.org/gituser/gitmodule.git GIT_SYNC_DEST=/mycheckout git-sync -loglevel DEBUG 
```

```
./git-sync -repo https://bitbucket.org/gituser/gitmodule.git -username gituser -password gitpassword -dest /mycheckout -loglevel DEBUG
```

## Continuous integration

[Jenkins](https://www.benjamin-borbe.de/jenkins/job/Go-Git-Sync/)
