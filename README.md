# Checker

## Overview

Checker is a utility to help ensure documentation quality at MongoDB while reducing error prone, tedious tasks.

## Use

The intended use is to check links in changed files. This can be accomplished with:

```sh
git diff --name-only HEAD master | tr "\n" "," | xargs checker -p --path . --changes
```

The above commands first get the list of file names with changes by comparing your current branch to master,
then converts that into a comma separated list, lastly passing the list to the program via `xargs`.

You can also check _all_ links by omitting the `--changes` flag, though this can take a very long time depending
on the size of the project.

See the `--help` flag for more info.

```sh
checker --help
```

## What it does

Specifically, it checks to ensure all links are valid. It does this in the
following ways:

- It will find all raw links and check them (https?...).
- It will find all [role uses](https://www.sphinx-doc.org/en/master/usage/restructuredtext/roles.html)
  defined in the latest release version of [rstspec.toml](https://github.com/mongodb/snooty-parser/blob/master/snooty/rstspec.toml)
  and check resulting interpreted urls.
- It will optionally check uses of `:doc:` and `:ref:` targets. **Note**: checker DOES NOT ignore rst comments. Use the
  optional `-d` and `-r` flags to check for `:doc:` and `:ref:` targets, respectively.

## How it does it

Once it scans files for checkable items, it begins checking them. URL checing is performed in a pool of workers
(default 10), configurable with the `-w` flag. Each worker is throttled (default 10), configurable with the `-t` flag, so that no
worker can issue more than (1e9 / (throttle / workers)) requests per second. **Setting this value too high can result in
inadvertent DOS attacks.**. `:ref:` targets are only checked for existence, since the URL is guaranteed to be accurate based
on the way they are generated. `:doc:` targets check whether the target is in the list of scanned files.
