#! /bin/bash

# MAIN FUNCTION
function PKG_name(){
## Includes
if [[ -f "${BASH_SOURCE%/*}/pkgfile" ]]; then
    source ${BASH_SOURCE%/*}/pkgfile
else
    echo "error: A pkgfile was not identified."
    return
fi
if [[ -f "${BASH_SOURCE%/*}/.env" ]]; then
    source ${BASH_SOURCE%/*}/.env
fi
## Auxiliary Functions
### (...)

## Main Function Properly
### without options enter in the interactive mode or print help
    if  [[ -z "$1" ]]; then
        if [[ -f "$INSTALL_DIR/src/interactive.sh" ]] &&
           [[ -s "$INSTALL_DIR/src/interactive.sh" ]]; then
            sh $INSTALL_DIR/src/interactive.sh
        else
            cat $INSTALL_DIR/src/help.txt
        fi
### "-c", -"cfg" and "--config" options to enter in the configuration mode
    elif ([[ "$1" == "-c" ]] || 
          [[ "$1" == "-cfg" ]] || 
          [[ "$1" == "--config" ]]) &&
          [[ -z "$2" ]]; then
          if [[ -f "$INSTALL_DIR/config.sh" ]] &&
           [[ -s "$INSTALL_DIR/config.sh" ]]; then
            sh $INSTALL_DIR/config.sh
        else
            echo "error: None configuration mode defined for the \"PKG_name()\" function."
        fi
### "-h" and "--help" options to print help
    elif ([[ "$1" == "-h" ]] || 
          [[ "$1" == "--help" ]]) &&
          [[ -z "$2" ]]; then
          cat $INSTALL_DIR/help.txt
### "-u" and "--uninstall" options to execute the uninstall script
    elif [[ "$1" == "-u" ]] || [[ "$1" == "--uninstall" ]]; then
        cd $INSTALL_DIR/install
        sh uninstall
        cd - > /dev/null

### (your additional code ...)

    else 
        echo "error: Option not defined for the \"PKG_name()\" function."
    fi
}
   
