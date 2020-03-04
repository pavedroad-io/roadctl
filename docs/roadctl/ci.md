# Travis integration

```yaml
dist: bionic                      << build on bionic
language: go                      << With go
go:
- '1.13                           << Go version(s)
env: "-GO111MODULE=on             << Modules on 
before_script:
- mv Makefile Bakefil             << Make Travis not see Makefile until we are ready
install: true                     << Use vendors
script:
- make -f Bakefile build          << build and crosscompile
- ls builds
deploy:
  provider: releases              << deploy GitHub release
  skip_cleanup: true              << don't remove files so we can attach them to release
  draft: true                     << mark the relase as draft so we can edit it
  overwrite: true                 << Overwrite if files exist
  api_key:
    secure: SHXUUm7m0s9/ZLOmFkThlj7LM9E/MBLQT0tSmF07u7Ieo7kgmYUFX62jcC0SywTmvFNIl7PE8l3v5fr797Y+gF3PXVoTb7GG6gI716a5AD9wWOOSMh/w6APyhe87HYXNrR6J34ZuHzR6AggaekjtjJd5nmPYUyBeD0dpdF6wBbQXc6xQQEIBhZgcKdFTUEfjQ5AvJ7JRJBXwcl3HVc6t2rVaZjJQ0c1bBFoK/3+S10t0FDlbCYljEL1Hsj34XEc7M1a7voYuSlh+08D6hZVJqnEhJKcWLSG3ZsIX62fAKTsSatD+iv2+TfGy6rBjtMgRPE8X+CK1ywNlhjx9KmcP1Wv4lvA1f6ErnQyPEjAmWRbDOcnEqjgRS6OZt1ta1ZA3Y5f9Uq4s25OIMnOOxk1b5zxzDZdtwpTXo1SqROIi512fzPSk+RO4PnY9aieXtVX/HzVg97/+2QGg8rA8wJcpwa05sXSOUeNRHUKZKusS8fLbkVIwL+U2ybo2PxDo+vEPPOrRhMpxvdPmb02Y+72TRO1uyiofPaOZLPJrKGN16ljtsnRdeB6GUZTDvkcVJyeOuC+gNAZHXq6RsmPriL3LP2MTePxZlwxerufUH0ElZlyO6PJQ6JyiC5DVmYb8KKvMq5dvabHmMaJwPFM07MFs0YAeS+rabnmLAaUvvksyTg8=
  file_glob: true                 << use glob pattern match
  file: "builds/*"                << copy builds executable to release
  on:
    repo: pavedroad-io/roadctl    << repo to use
    branch: travisSetup2          << branch to use
    tags: true                    << only if a git tag exists
		                              << does not fire on a PR
```

# FOSSA integration
Is execute on make check

# SoncarCloud integration
Is execute on make check.  Will break if GoGopkg files exist.


# GoDoc report card
Is integrated

# README.md
Has badges for all the above
