# Docs Link Checker

## Overview

Checker is designed to quickly check hyperlinks, :ref: directives, and :doc: 
directives to ensure they are valid. Based on Nathan Leniz's [tool of the same 
name](https://github.com/terakilobyte/checker), this fork improves performance, 
adds the ability to exclude specific URLS, and changes the default values.

## Install

```sh
go install github.com/MongoCaleb/checker@latest
```

## Use

There are two ways to use checker. This first is to run it against the entire 
docset. With multithreading, this process takes a matter of seconds. To do this, 
simply navigate to the root directory of your docs repo and run ``checker``. 

You can also configure the tool to only check files changed in the current 
diff. This can be accomplished with:

```sh
git diff --name-only | tr "\n" "," | xargs checker  --changes
```

**NOTE:** To check recent files, be sure to run this before adding the files to 
the current commit (before running ``git add``.)

See the `--help` flag for more info.

```sh
checker --help
```

## Excluding links

There are times when you may want to not check URLs. For example, if your docset 
has examples that use fake URLs, you want to make sure those URLs are ignored. 
One common example is to exclude checking http://example.com URLs.

To exclude URLS, create the following file:
``./config/link_checker_bypass_list.json``

In this file, add the URLs to be excluded and the reason for the exclusion in the 
following format:

```
[
    {
        "exclude":"example.com",
        "reason":"is not real url"
    },
    {
        "exclude":"api/client/v2.0",
        "reason":"will always return 400"
    }
]
```

## Running as a Github Action.

TBD. See https://github.com/actions/setup-go.

## What it does

Specifically, it checks to ensure all links are valid. It does this in the
following ways:

- It will find all raw links and check them (https?...).
- It will find all [role uses](https://www.sphinx-doc.org/en/master/usage/restructuredtext/roles.html)
  defined in the latest release version of [rstspec.toml](https://github.com/mongodb/snooty-parser/blob/master/snooty/rstspec.toml)
  and check resulting interpreted urls.
- It will optionally check uses of `:doc:` and `:ref:` targets. **Note**: checker DOES NOT ignore rst comments. Use the
  optional `-d` and `-r` flags to check for `:doc:` and `:ref:` targets, respectively.

