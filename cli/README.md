# `certgot/cli`

---

The purpose behind this package is to provide a library to parse any possible valid argument list for certbot. Certbot
uses python [argparse](https://docs.python.org/3/library/argparse.html), however with some modifications to make it work
how they want.

The aim of this library is to replicate the desired behaviours from argparse while also incorporating things from
certbot, such as configuration values that query the user if appropriate.

Other libraries do exist that replicate argparse more faithfully, but considering they would still require modifications
or hacks, it is just easier to write our own while using less (ideally: no) external dependencies.

This package will probably be a bit of a moving target while certbot is being implemented.

There are a few main components to this package. They are,

## Application (`App`)

An `App` is an entry point to a cli application.

## Arguments

Arguments take 3 forms, a [Flag](#flag-argument), a [Command](#command-argument), or [Extra Arguments](#extra-arguments).

Go provides arguments as part of the (Os.Args)[https://golang.org/pkg/os/#Args] field, and takes care of all the
splitting on spaces, and long spaced arguments grouped by quotation marks.

### Flag Argument

A flag is a pretty standard [POSIX/GNU compatible argument](https://www.gnu.org/software/libc/manual/html_node/Argument-Syntax.html).
They come in a few forms that this package accepts.

#### Short flag

A short flag is a single dash '-' followed by a character. The character can optionally be repeated multiple times, and 
also optionally include a value, separated either by a space or an equals sign. The following list shows a few
non-exhaustive examples of valid short flags,

    -a
    -aaaa
    -a=foo
    -a="foo bar"
    -a foo
    -a "foo bar"

Unlike the POSIX/GNU arguments, we do not accept multiple grouped short flags. That is `-abc` is an invalid flag.

#### Long Flag

A long flag is any alphanumeric character following two dash '--' characters, optionally spacing words by another dash.

The following list shows a few non-exhaustive examples of valid long flags,

    --hello
    --hello=world
    --hello world
    --hello "cruel world"
    --foo-bar
    --foo-bar=hello
    --foo-bar hello
    --foo-bar "hello world"

### Command Argument

Command arguments are the first argument encountered that is not associated as a value of a flag.

### Extra Arguments

Extra arguments are any other arguments encountered that are not associated with a flag after a command.

## Configuration

Configuration values are set in a configuration file, or by a flag. They can later be queried by an application and store user,
or system input, immutable-once-program-is-running state. They don't have to be immutable I guess, but are typically 
only set when initialising the program.

## Help Category

Help categories are how commands and flags are grouped when printing help.