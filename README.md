# everycommit

[![CircleCI](https://circleci.com/gh/bensallen/everycommit.svg?style=svg)](https://circleci.com/gh/bensallen/everycommit)

A helper application to iterate through every commit of a Github PR, running a given command. Queries the GitHub API to find commits in a PR. Useful for ensuring every commit in a PR compiles and successfully can run tests.

Example:

```
$ everycommit -r /path/to/repo_checkout -u https://github.com/owner/project/pull/1 -- git show-ref --head ^HEAD
Checking out commit: 6c5e7c97b4eead4fabaea34b2e7e37d6d966369a
Running: /usr/local/bin/git show-ref --head ^HEAD
6c5e7c97b4eead4fabaea34b2e7e37d6d966369a HEAD
Checking out commit: 5c6bef13efa6367630e50bcc225fafc7560d0e85
Running: /usr/local/bin/git show-ref --head ^HEAD
5c6bef13efa6367630e50bcc225fafc7560d0e85 HEAD
Checking out commit: f5024614a16db5959f0cea6c63880e5d5d1b4f22
Running: /usr/local/bin/git show-ref --head ^HEAD
f5024614a16db5959f0cea6c63880e5d5d1b4f22 HEAD
Checking out commit: bc87d4e8df79e0475204325258cffc6f11e4f71e
Running: /usr/local/bin/git show-ref --head ^HEAD
bc87d4e8df79e0475204325258cffc6f11e4f71e HEAD
```
