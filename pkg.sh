#! /bin/bash

installdir=$HOME/.config/pkg.sh
    function pkg(){
        function pkg_build(){
            echo "Creating base directories..."
            mkdir $2
            mkdir $2/install
            mkdir $2/config
            mkdir $2/files
            echo "Creating base files..."
            touch $2/install/install
            touch $2/install/uninstall
            touch $2/install/configure
            touch $2/config/config
            touch $2/config/help
            touch $2/$1
            echo "Configuring the base files..."
            cat $installdir/files/install >> $2/install/install
            cat $installdir/files/uninstall >> $2/install/uninstall
            cat $installdir/files/configure >> $2/install/configure
            cat $installdir/files/config >> $2/config/config
            cat $installdir/files/help >> $2/config/help
            cat $installdir/files/base >> $2/$1
            echo "Configuring files to be executable..."
            chmod a+x $2/install/install
            chmod a+x $2/install/uninstall
            chmod a+x $2/install/configure
            chmod a+x $2/$1
            echo "Copying the pkgfile..."
        }
    
    if [[ -z $1 ]]; then
        if [[ -f pkgfile ]]; then
            has_name=$(grep -R "PKG_name" "pkgfile")
            has_dir=$(grep -R "PKG_dir" "pkgfile")
            given_name=$(source pkgfile && echo "$PKG_name")
            given_dir=$(source pkgfile && echo "$PKG_dir")
            if [[ -n $has_name ]] && [[ -n $has_dir ]]; then
                echo "Creating the package \"$given_name\" in the directory \"$given_dir\"..."
                pkg_build $given_name $given_dir
                cp -r pkgfile $given_dir
            else
                echo "Your \"pkgfile\" is not well constructed. Please provide at least a package name and dir."
            fi
        else
            echo "\"pkgfile\" not found. Generate it with \"pkg --config\", or try a quick build. See \"pkg --help\"."
        fi
    elif [[ -f $1 ]]; then
        has_name=$(grep -R "PKG_name" "$1")
        has_dir=$(grep -R "PKG_dir" "$1")
        given_name=$(source $1 && echo "$PKG_name")
        given_dir=$(source $1 && echo "$PKG_dir")
        if [[ -n $has_name ]] && [[ -n $has_dir ]]; then
            echo "Creating the package \"$given_name\" in the directory \"$given_dir\"..."
            pkg_build $given_name $given_dir
            cp -r $1 $given_dir
        else
            echo "Your \"pkgfile\" is not well constructed. Please provide at least a package name and dir."
        fi
    elif [[ "$1" == "-t" ]] || [[ "$1" == "-tpl" ]] || [[ "$1" == "--template" ]]; then
        echo "Copying a template for the \"pkgfile\"..."
            cp -r $installdir/pkgfile $PWD
        echo "Done."
    elif [[ "$1" == "-h" ]] || [[ "$1" == "--help" ]]; then
        echo "displaying help..."

    else
        echo "Please provide a valid path to a pkgfile."
    fi
       unset -f pkg_build
    }
