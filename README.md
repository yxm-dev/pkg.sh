About
=======

[pkg.sh](https://github.com/yxm-dev/pkg.sh) is a simple package builder written in pure shell and aimed to
provide standardization to pure bash projects, with configurations prescribed in a `pkgfile` in the same
spirit of `GNU make`.

The packages created with [pkg.sh](https://github.com/yxm-dev/pkg.sh) comes with `install` and`unistall`
scripts that install/uninstall the needed dependencies (defined in the `pkgfile`). The uninstall process fully
remove the package and the installed dependencies.

The `pkgfile` file can be created manually (a template is provided with `pkg --template`) or, optionally,
through a TUI interface written in `Go`, accessed with `pkg --pkgfile`. 

Dependencies
=============

* `Bash` and `sed` for general usage
* `Go` for the TUI interface (it can be automatically installed from `pkg --pkgfile`).

Install
========

1. Clone this repository

```bash
    git clone https://github.com/yxm-dev/pkg.sh

```

2. enter in the `install` directory

```bash
    cd pkg.sh/install
```

3. execute `./configure` and select the installation directory
4. execute `./install`

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






