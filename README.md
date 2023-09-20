About
=======

[pkg.sh](https://github.com/yxm-dev/pkg.sh) is a simple package builder written in pure shell and
aimed to provide standardization to [bash](https://www.gnu.org/software/bash/) projects, in the same lines of [GNU make](https://www.gnu.org/software/make/). 
Here, the `makefile` is replaced with a `pkgfile`, where the new package name and its dependencies are
provided. The building process is organized in the archetypal `./configure` and `./install` steps. A
`./uninstall` script is also provided. 

The `pkgfile` file can be created manually (where a template is provided) or, optionally,
through a TUI interface written in [Go](https://go.dev/) with the help of the
[tview](https://github.com/rivo/tview) library.

For more details, visit the [project website](http://yxm-atip.s3-website-sa-east-1.amazonaws.com/).

Why?
=====

* to replace the use of [GNU make](https://www.gnu.org/software/make/) for the case of simple
projects written in , avoiding additional dependences and still providing standardization;
* to provide a package builders designed for [bash](https://www.gnu.org/software/bash/) applications.

Dependencies
=============

* [bash](https://www.gnu.org/software/bash/) and [sed](https://www.gnu.org/software/sed/) for general use
* [Go](https://go.dev/) for the TUI interface (optional)

Install
========

1. Clone this repository

```bash
    git clone https://github.com/yxm-dev/pkg.sh

```

2. Enter in the `install` directory

```bash
    cd pkg.sh/install
```

3. Execute `./configure` and select the installation directory `PKG_install_dir`
4. Execute `./install`

* To uninstall, enter in the installation directory and execute the `uninstall` script:

```bash
    cd PKG_install_dir/install
    ./uninstall
```

Usage
========

1. create a `pkgfile` containing the new package name `PKG_name` and its dependencies:
    * `pkg --template` to get a template;
    * `pkg --config` to use a TUI interface.
2. execute `pkg /path/to/pkgfile` or `pkg` if the `pkgfile` is in the working directory;
3. enter in the package directory with `cd PKG_name`;
4. configure it as you want.

* To install the `PKG_name` and all its dependencies, follow the same process used to install `pkg.sh`:

```bash

    cd PKG_name/install
    ./configure
    ./install

```
* In the `./configure` step you will fix the install directory `install_dir` of `PKG_name`. To
  uninstall the package (and optionally its dependencies), execute the `uninstall` in `install_dir`:
  
```bash
    cd install_dir/install
    ./uninstall
```

To Do
===========

* improve the error messages in the TUI interface

