#! /bin/bash

# PKG FUNCTION
    function pkg(){
# Including the pkgfile
    source ${BASH_SOURCE%/*}/pkgfile
## Auxiliary Functions
        function PKG_build(){
            echo "Creating base directories..."
            if [[ -f "$1" ]]; then
                echo "*error:* There alrealy exists a file \"$1\". Change the package name in the pkgfile."
            elif [[ -d "$1" ]]; then
                 echo "*error:* There alrealy exists a directory \"$1\". Change the package name in the pkgfile."
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
                touch $1/config/help.txt
                touch $1/config/interactive
                touch $1/$1
                echo "Configuring base files..."
                cat $PKG_install_dir/files/install >> $1/install/install
                cat $PKG_install_dir/files/uninstall >> $1/install/uninstall
                cat $PKG_install_dir/files/configure >> $1/install/configure
                cat $PKG_install_dir/files/config >> $1/config/config
                cat $PKG_install_dir/files/package_manager >> $1/config/package_manager
                cat $PKG_install_dir/files/help.txt >> $1/config/help.txt
                cat $PKG_install_dir/files/base >> $1/$1
                echo "Configuring files to be executable..."
                chmod a+x $1/install/install
                chmod a+x $1/install/uninstall
                chmod a+x $1/install/configure
                chmod a+x $1/config/config
                chmod a+x $1/$1
                if [[ -f "pkgfilecd" ]]; then
                    echo "Configuring pkgfilecd file..."
                    mv pkgfilecd $1/install/pkgfilecd
                fi
                echo "Copying the pkgfile..."
                echo "Confiruring \"$1\" file..."
                sed -i "s/PKG_name/$1/g" $1/$1
                echo "Configuring the help.txt file..."
                sed -i "s/PKG_name/$1/g" $1/config/help.txt
                echo "The application \"$1\" has been created in your working directory."
            fi
        }
## PKG Function Properly   
    if [[ -z $1 ]]; then
        if [[ -f pkgfile ]]; then
            local has_name=$(grep -R "PKG_name" "pkgfile")
            local given_name=$(source pkgfile && echo "$PKG_name")
            if [[ -z $given_name ]] || [[ -z $has_name ]] ; then
                echo "Your \"pkgfile\" is not well constructed. Please provide at least a package name."
            else
                echo "Creating the package \"$given_name\"..."
                PKG_build $given_name
                cp -r pkgfile $given_name
            fi
        else
            echo "\"pkgfile\" not found. Generate it with \"pkg --config\", or try a quick build. See \"pkg --help\"."
        fi
    elif [[ -f $1 ]]; then
        local has_name=$(grep -R "PKG_name" "$1")
        local given_name=$(source $1 && echo "$PKG_name")
        if [[ -z $given_name ]] || [[ -z $has_dir ]]; then
            echo "Your \"pkgfile\" is not well constructed. Please provide at least a package name and dir."
        else
            echo "Creating the package \"$given_name\"..."
            PKG_build $given_name
            cp -r $1 $given_name
        fi
    elif [[ "$1" == "-t" ]] || [[ "$1" == "-tpl" ]] || 
         [[ "$1" == "--template" ]] ||
         [[ "$1" == "-pkg" ]] ||
         [[ "$1" == "--pkgfile" ]]; then
        echo "Copying a template for the \"pkgfile\"..."
            cp -r $PKG_install_dir/pkgfile $PWD
        echo "Done."
    elif [[ "$1" == "-h" ]] || [[ "$1" == "--help" ]]; then
        cat $PKG_install_dir/config/help.txt

    elif [[ "$1" == "-c" ]] || [[ "$1" == "-cfg" ]] || [[ "$1" == "--config" ]] ||
         [[ "$1" == "-pkgfile" ]] || [[ "$1" == "--pkgfile" ]]; then
         sh $PKG_install_dir/config/config
         mv $PKG_install_dir/config/ui/pkgfile $PWD
         if [[ -f $PKG_install_dir/config/ui/pkgfilecd ]]; then
             mv $PKG_install_dir/config/ui/pkgfilecd $PWD
         fi
    else
        echo "Please provide a valid path to a pkgfile."
    fi
       unset -f PKG_build
    }
