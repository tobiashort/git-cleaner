# git-cleaner

Attention! Data loss!

1. find all repositories
4. checkout master or main
2. reset --hard --recurse-submodules
3. clean -fd
5. pull -p
6. optionally: delete all local branches

## Build

```shell
go run ./build
```
