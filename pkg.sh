#! /bin/bash

install_dir=$HOME/.config/pkg.sh

    function pkg(){
        function pkg_build(){
            echo "Creating base directories..."
            if [[ -f "$1" ]]; then
                echo "*ERROR* There alrealy exists a file \"$1\". Change the package name in the pkgfile."
            elif [[ -d "$1" ]]; then
                 echo "*ERROR* There alrealy exists a directory \"$1\". Change the package name in the pkgfile."
            else
                mkdir $1
                mkdir $1/install
                mkdir $1/config
                mkdir $1/files
                echo "Creating base files..."
                touch $1/install/install
                touch $1/install/uninstall
                touch $1/install/configure
                touch $1/config/config
                touch $1/config/help
                touch $1/files/interactive
                touch $1/$1.sh
                echo "Configuring base files..."
                cat $install_dir/files/install >> $1/install/install
                cat $install_dir/files/uninstall >> $1/install/uninstall
                cat $install_dir/files/configure >> $1/install/configure
                cat $install_dir/config/config >> $1/config/config
                cat $install_dir/config/distros >> $1/config/distros
                cat $install_dir/files/help >> $1/config/help
                cat $install_dir/files/base >> $1/$1.sh
                echo "Configuring files to be executable..."
                chmod a+x $1/install/install
                chmod a+x $1/install/uninstall
                chmod a+x $1/install/configure
                chmod a+x $1/config/config
                chmod a+x $1/$1.sh
                echo "Copying the pkgfile..."
            fi
        }
    
    if [[ -z $1 ]]; then
        if [[ -f pkgfile ]]; then
            has_name=$(grep -R "PKG_name" "pkgfile")
            given_name=$(source pkgfile && echo "$PKG_name")
            if [[ -z $given_name ]] || [[ -z $has_name ]] ; then
                echo "Your \"pkgfile\" is not well constructed. Please provide at least a package name and dir."
            else
                echo "Creating the package \"$given_name\"..."
                pkg_build $given_name
                cp -r pkgfile $given_name
            fi
        else
            echo "\"pkgfile\" not found. Generate it with \"pkg --config\", or try a quick build. See \"pkg --help\"."
        fi
    elif [[ -f $1 ]]; then
        has_name=$(grep -R "PKG_name" "$1")
        given_name=$(source $1 && echo "$PKG_name")
        if [[ -z $given_name ]] || [[ -z $has_dir ]]; then
            echo "Your \"pkgfile\" is not well constructed. Please provide at least a package name and dir."
        else
            echo "Creating the package \"$given_name\"..."
            pkg_build $given_name
            cp -r $1 $given_name
        fi
    elif [[ "$1" == "-t" ]] || [[ "$1" == "-tpl" ]] || [[ "$1" == "--template" ]]; then
        echo "Copying a template for the \"pkgfile\"..."
            cp -r $install_dir/pkgfile $PWD
        echo "Done."
    elif [[ "$1" == "-h" ]] || [[ "$1" == "--help" ]]; then
        echo "displaying help..."

    else
        echo "Please provide a valid path to a pkgfile."
    fi
       unset -f pkg_build
    }
