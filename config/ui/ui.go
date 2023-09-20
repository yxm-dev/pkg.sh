package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// // INPUT // //

var distros = []string{
	"Debian",
	"Arch",
	"RedHat",
	"openSUSE",
}

// // VARIABLES // //
var app *tview.Application
var err error
var err2 error
var packageName string
var debian string
var basicDependencyName string
var commandDistros string
var customDependencyName string
var customCommands string
var customInstallCommands string

var txtIntro string
var bodyIntro *tview.TextView
var footerIntro *tview.TextView
var pageIntro *tview.Flex

var txtName string
var headerName *tview.TextView
var bodyName *tview.Form
var footerName *tview.TextView
var pageName *tview.Flex

var txtDependenciesIntro string
var bodyDependenciesIntro *tview.TextView
var footerDependenciesIntro *tview.TextView
var pageDependenciesIntro *tview.Flex

var txtBasicDependenciesAdd string
var headerBasicDependenciesAdd *tview.TextView
var bodyBasicDependenciesAdd *tview.Form
var pageBasicDependenciesAdd *tview.Flex
var footerBasicDependenciesAdd *tview.TextView

var txtCustomDependenciesAdd string
var headerCustomDependenciesAdd *tview.TextView
var bodyCustomDependenciesAdd *tview.Form
var pageCustomDependenciesAdd *tview.Flex
var footerCustomDependenciesAdd *tview.TextView

var txtCustomInstallIntro string
var bodyCustomInstallIntro *tview.TextView
var footerCustomInstallIntro *tview.TextView
var pageCustomInstallIntro *tview.Flex

var txtCustomInstallAdd string
var headerCustomInstallAdd *tview.TextView
var bodyCustomInstallAdd *tview.Form
var pageCustomInstallAdd *tview.Flex
var footerCustomInstallAdd *tview.TextView

var pages *tview.Pages

// // FUNCTIONS // //
// check if there already exists a file "pkgfile" in the working directory.
// If don't, create it.
func createPkgfile() {
	_, err := os.Stat("pkgfile")
	if err == nil {
		app.Stop()
		fmt.Println("error: There already exists a pkgfile in the working directory.")
	} else {
		file, err := os.OpenFile("pkgfile", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		shellHeader := []byte("#! /bin/bash \n\n")
		file.Write(shellHeader)
	}
}

// verify if the variable "PKG_name" is defined
// print it to the pkgfile, if not.
func verifyName(packageName string) {
	file, err := os.OpenFile("pkgfile", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	content, err := ioutil.ReadFile("pkgfile")
	if err != nil {
		panic(err)
	}

	fileContent := string(content)
	if strings.Contains(fileContent, "PKG_name=") {

	} else {
		packageNameAdd := []byte("PKG_name=" + "\"" + packageName + "\"\n\n")
		file.Write(packageNameAdd)
	}
}

// create the file "pkgfile_distros" defining the array of distros to include.
// if it already exists, delete it first.
func createDistrosInclude() {
	_, err := os.Stat("pkgfile_distros")
	if err == nil {
		os.Remove("pkgfile_distros")
		if err != nil {
			panic(err)
		}
	} else {
		file, err := os.OpenFile("pkgfile_distros", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		var distroIncludeYes string

		for _, distro := range distros {
			distroIncludeYes = distroIncludeYes +
				"PKG_distro_include[\"" + distro + "\"]=\"yes\"\n"
		}

		distrosIncludeArray := []byte("declare -A PKG_distro_include\n\n" + distroIncludeYes)

		file.Write(distrosIncludeArray)
	}
}

// append the file "pkgfile" with the file "pkgfile_distros"
func appendDistrosIncludePkgfile() {
	file1, err := os.OpenFile("pkgfile", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file1.Close()

	file2, err := os.Open("pkgfile_distros")
	if err != nil {
		panic(err)
	}
	defer file2.Close()

	_, err = io.Copy(file1, file2)
	if err != nil {
		panic(err)
	}

}

// create the file "pkgfile_dependencies" where the array of dependencies to be used will be built.
// if it already exists, delete it first.
func createDependenciesPkgfile() {
	_, err := os.Stat("pkgfile_dependencies")
	if err == nil {
		os.Remove("pkgfile_dependencies")
		if err != nil {
			panic(err)
		}
	} else {
		file, err := os.OpenFile("pkgfile_dependencies", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		dependenciesArray := []byte("declare -a PKG_dependencies\n\nPKG_dependencies=(")
		file.Write(dependenciesArray)
	}
}

// add a dependency to the array of dependencies in the file "pkgfile_dependencies".
func appendDependenciesPkgfile(dependency string) {
	file, err := os.OpenFile("pkgfile_dependencies", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	additionalDependency := []byte(dependency + " ")
	file.Write(additionalDependency)

}

// add the latest dependency to the array of dependencies, concluding it.
func appendLastDependencyPkgfile(dependency string) {
	file, err := os.OpenFile("pkgfile_dependencies", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lastDependency := []byte(dependency + ")")
	file.Write(lastDependency)
}

// create the file "pkgfile_dependencies_distros" where the bidimensional array
// of commands of packages in distros will be build.
// If it already exists, delete it first.
func createDependenciesDistrosPkgfile() {
	_, err := os.Stat("pkgfile_dependencies_distros")
	if err == nil {
		os.Remove("pkgfile_dependencies_distros")
		if err != nil {
			panic(err)
		}
	} else {
		file, err := os.OpenFile("pkgfile_dependencies_distros", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		dependenciesDistroArray := []byte("declare -A PKG_distro_package_name\n\n")
		file.Write(dependenciesDistroArray)
	}
}

// append "pkgfile_dependencies_distros" with the command of the package in the distros:
func appendDependenciesDistrosPkgfile(name, command string) {
	file, err := os.OpenFile("pkgfile_dependencies_distros", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var commandDistros string

	for _, distro := range distros {
		commandDistros = commandDistros +
			"PKG_distro_package_name[\"" + distro + "\",\"" + name + "\"]=\"" + command + "\"\n"
	}

	commandByte := []byte(commandDistros)
	file.Write(commandByte)
}

// create a form to add the package name and then move to dependencies or custom install script
func addName() {
	bodyName.AddInputField("name:", "", 30, nil, nil).
		SetFieldBackgroundColor(tcell.ColorGray).
		SetFieldTextColor(tcell.ColorWhite).
		SetLabelColor(tcell.ColorWhite)

	bodyName.AddButton("add dependencies", func() {
		packageName = bodyName.GetFormItemByLabel("name:").(*tview.InputField).GetText()
		verifyName(packageName)
		createDistrosInclude()
		createDependenciesPkgfile()
		createCustomDependencies()
		pages.SwitchToPage("Dependencies Intro")
	}).
		SetLabelColor(tcell.ColorWhite).
		SetBorder(false)

	bodyName.AddButton("add install script", func() {
		packageName = bodyName.GetFormItemByLabel("name:").(*tview.InputField).GetText()
		verifyName(packageName)
		pages.SwitchToPage("Install Intro")
	}).
		SetLabelColor(tcell.ColorWhite).
		SetBorder(false)

	bodyName.AddButton("conclude", func() {
		concludePkgfile()
		deleteConcludePkgfile()
		app.Stop()
		fmt.Println("Pkgfile has been created in your working directory. Type \"pkg\" to create your package.")
	}).
		SetLabelColor(tcell.ColorWhite).
		SetBorder(false)

}

// create a form of the dependencies names and commands in distros.
// print the result in the files "pkgfile_dependencies" and "pkgfile_dependencies_distros".
func addBasicDependency() {

	var labels []string

	labels = append(labels, "dependency name:")
	labels = append(labels, "Debian based:")
	labels = append(labels, "Arch Linux based:")
	labels = append(labels, "Red Hat based:")
	labels = append(labels, "SUSE based:")

	for i := 0; i < len(labels); i++ {
		bodyBasicDependenciesAdd.AddInputField(labels[i], "", 30, nil, nil).
			SetFieldBackgroundColor(tcell.ColorGray).
			SetFieldTextColor(tcell.ColorWhite).
			SetLabelColor(tcell.ColorWhite)
	}

	bodyBasicDependenciesAdd.AddButton("add other", func() {
		basicDependencyName = bodyBasicDependenciesAdd.GetFormItemByLabel(labels[0]).(*tview.InputField).GetText()
		commandDistros = bodyBasicDependenciesAdd.GetFormItemByLabel(labels[1]).(*tview.InputField).GetText()

		appendDependenciesPkgfile(basicDependencyName)
		appendDependenciesDistrosPkgfile(basicDependencyName, commandDistros)

		pages.SwitchToPage("Dependencies Intro")
	}).
		SetLabelColor(tcell.ColorWhite).
		SetBorder(false)

	bodyBasicDependenciesAdd.AddButton("return", func() {
		bodyName.Clear(true)
		addName()
		pages.SwitchToPage("Name")
	}).
		SetLabelColor(tcell.ColorWhite).
		SetBorder(false)

	bodyBasicDependenciesAdd.AddButton("conclude", func() {
		basicDependencyName = bodyBasicDependenciesAdd.GetFormItemByLabel(labels[0]).(*tview.InputField).GetText()
		commandDistros = bodyBasicDependenciesAdd.GetFormItemByLabel(labels[1]).(*tview.InputField).GetText()

		appendLastDependencyPkgfile(basicDependencyName)
		appendDependenciesDistrosPkgfile(basicDependencyName, commandDistros)

		concludePkgfile()
		deleteConcludePkgfile()
		app.Stop()
		fmt.Println("Pkgfile has been created in your working directory. Type \"pkg\" to create your package.")
	}).
		SetLabelColor(tcell.ColorWhite).
		SetBorder(false)
}

// create the file "pkgfilecd" of custom dependencies
// if it already exists, do nothing.
func createCustomDependencies() {
	_, err := os.Stat("pkgfilecd")
	if err == nil {
		os.Remove("pkgfilecd")
		if err != nil {
			panic(err)
		}
	} else {
		file, err := os.OpenFile("pkgfilecd", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		txtCustomDependenciesHeader := "#! /bin/bash \n\n"
		customDependenciesHeader := []byte(txtCustomDependenciesHeader)
		file.Write(customDependenciesHeader)
	}
}

func appendCustomDependencyCheck(dependency string) {
	file, err := os.OpenFile("pkgfilecd", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	customDependencyNameCheck := []byte("echo \"Checking for custom dependency " + dependency + "...\"\n")
	customDependencyNameCheckIf := "if [[ -x " + dependency + " ]]; then \n    echo \"ok...\"\n"
	customDependencyNameCheckElse := "else\n    echo echo \"Installing custom dependency " + dependency + "...\"\nfi\n\n"
	customDependencyNameCheckIfElse := []byte(customDependencyNameCheckIf + customDependencyNameCheckElse)

	file.Write(customDependencyNameCheck)
	file.Write(customDependencyNameCheckIfElse)
}

func appendCustomDependencyCommands(commands string) {
	file, err := os.OpenFile("pkgfilecd", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	customDependencyCommands := []byte(commands + "\n\n")

	file.Write(customDependencyCommands)
}

// add custom dependency to the file "pkgfilecd"
func addCustomDependency() {
	bodyCustomDependenciesAdd.AddInputField("dependency name:", "", 30, nil, nil).
		SetFieldBackgroundColor(tcell.ColorGray).
		SetFieldTextColor(tcell.ColorWhite).
		SetLabelColor(tcell.ColorWhite)

	bodyCustomDependenciesAdd.AddTextArea("dependency script:", "", 60, 0, 0, nil).
		SetFieldBackgroundColor(tcell.ColorGray).
		SetFieldTextColor(tcell.ColorWhite).
		SetLabelColor(tcell.ColorWhite)

	bodyCustomDependenciesAdd.AddButton("add other", func() {
		customDependencyName = bodyCustomDependenciesAdd.GetFormItemByLabel("dependency name:").(*tview.InputField).GetText()
		customCommands = bodyCustomDependenciesAdd.GetFormItemByLabel("dependency script:").(*tview.TextArea).GetText()

		appendCustomDependencyCheck(customDependencyName)
		appendCustomDependencyCommands(customCommands)

		pages.SwitchToPage("Dependencies Intro")
	}).
		SetLabelColor(tcell.ColorWhite).
		SetBorder(false)

	bodyCustomDependenciesAdd.AddButton("return", func() {
		bodyName.Clear(true)
		addName()
		pages.SwitchToPage("Name")
	}).
		SetLabelColor(tcell.ColorWhite).
		SetBorder(false)

	bodyCustomDependenciesAdd.AddButton("conclude", func() {
		customDependencyName = bodyCustomDependenciesAdd.GetFormItemByLabel("dependency name:").(*tview.InputField).GetText()
		customCommands = bodyCustomDependenciesAdd.GetFormItemByLabel("install script:").(*tview.TextArea).GetText()

		appendCustomDependencyCheck(customDependencyName)
		appendCustomDependencyCommands(customCommands)
		appendLastDependencyPkgfile("")
		concludePkgfile()
		deleteConcludePkgfile()
		app.Stop()
		fmt.Println("Pkgfile has been created in your working directory. Type \"pkg\" to create your package.")
	}).
		SetLabelColor(tcell.ColorWhite).
		SetBorder(false)

}

// create the file "pkgfileci" of custom install steps
// if it already exists, do nothing.
func createCustomInstall() {
	_, err := os.Stat("pkgfileci")
	if err == nil {
		os.Remove("pkgfileci")
		if err != nil {
			panic(err)
		}
	} else {
		file, err := os.OpenFile("pkgfileci", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		txtCustomInstallHeader := "#! /bin/bash \n\n"
		customInstallHeader := []byte(txtCustomInstallHeader)
		txtCustomInstallBase := "echo \"Executing custom installation steps for the package \\\"$PKG_name\\\"...\"\n\n"
		customInstallBase := []byte(txtCustomInstallBase)
		file.Write(customInstallHeader)
		file.Write(customInstallBase)
	}
}

func appendCustomInstallCommands(commands string) {
	file, err := os.OpenFile("pkgfileci", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	customInstallCommands := []byte(commands + "\n\n")

	file.Write(customInstallCommands)
}

// verify if some basic dependency was defined and, in an affirmative case
// call lastDependency()
func verifyBasicDependency(packageName string) {
	file, err := os.OpenFile("pkgfile", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "PKG_dependencies") {
			line += ")"

		}

	}

}

// add custom install script to the file "pkgfileci"
func addCustomInstall() {
	bodyCustomInstallAdd.AddTextArea("install script:", "", 60, 0, 0, nil).
		SetFieldBackgroundColor(tcell.ColorGray).
		SetFieldTextColor(tcell.ColorWhite).
		SetLabelColor(tcell.ColorWhite)

	bodyCustomInstallAdd.AddButton("conclude", func() {
		customInstallCommands = bodyCustomInstallAdd.GetFormItemByLabel("install script:").(*tview.TextArea).GetText()

		appendCustomInstallCommands(customInstallCommands)
		concludePkgfile()
		deleteConcludePkgfile()
		app.Stop()
		fmt.Println("Pkgfile has been created in your working directory. Type \"pkg\" to create your package.")
	}).
		SetLabelColor(tcell.ColorWhite).
		SetBorder(false)

	bodyCustomInstallAdd.AddButton("return", func() {
		bodyName.Clear(true)
		addName()
		pages.SwitchToPage("Name")
	}).
		SetLabelColor(tcell.ColorWhite).
		SetBorder(false)

}

// function to pass when the conclude button is pressed.
// append the file "pkgfile" with the other files.
func concludePkgfile() {
	var pkgfiles = []string{
		"pkgfile_distros",
		"pkgfile_dependencies",
		"pkgfile_dependencies_distros",
	}

	target, err := os.OpenFile("pkgfile", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer target.Close()
	if err != nil {
		panic(err)
	}

	for _, file := range pkgfiles {
		_, err := os.Stat(file)
		if err == nil {
			source, err := os.Open(file)
			if err != nil {
				panic(err)
			} else {
				_, err = io.Copy(target, source)
				if err != nil {
					panic(err)
				}
			}
			defer source.Close()
			if err != nil {
				panic(err)
			}
		}

		newLine := []byte("\n\n")
		_, err = target.Write(newLine)
		if err != nil {
			panic(err)
		}

	}
}

// delete the additional files.
// delete "pkgfilecd" and "pkgfileci" if empty.
// used when concluding the configuration
func deleteConcludePkgfile() {
	var pkgfiles = []string{
		"pkgfile_distros",
		"pkgfile_dependencies",
		"pkgfile_dependencies_distros",
	}

	for _, file := range pkgfiles {
		_, err := os.Stat(file)
		if err == nil {
			os.Remove(file)
		}
	}

	file, err := os.Stat("pkgfilecd")
	if err == nil {
		if file.Size() == 0 {
			os.Remove("pkgfilecd")
		}
	}
	file, err = os.Stat("pkgfileci")
	if err == nil {
		if file.Size() == 0 {
			os.Remove("pkgfileci")
		}
	}

}

// delete the additional files and the "pkgfile".
// used when exiting without concluding.

func deletePkgfile() {
	deleteConcludePkgfile()
	_, err := os.Stat("pkgfile")
	if err == nil {
		os.Remove("pkgfile")
	}

	_, err2 := os.Stat("pkgfilecd")
	if err2 == nil {
		os.Remove("pkgfilecd")
	}

	_, err3 := os.Stat("pkgfileci")
	if err3 == nil {
		os.Remove("pkgfileci")
	}

}

// // MAIN // //
func main() {

	// initializing the app
	app = tview.NewApplication()

	// PAGE INTRO - text
	txtIntro = `
  Welcome to the configurarion mode of pkg.sh.

  In the following you will create a "pkgfile" by providing:
	1. the new package name;
	2. its dependencies;
	3. the corresponding commands in each class of distributions.

  You will also be able to create a "pkgfilecd", by providing
  customized scripts ot install custom dependencies. This is
  particularly useful to install dependencies from the source.

  Finally, if needed, you will also be able to create a "pkgfileci"
  containing customized steps in the installation process of your
  Bash package.
`
	// PAGE INTRO - body
	bodyIntro = tview.NewTextView().
		SetText(txtIntro)

	// PAGE INTRO - footer
	footerIntro = tview.NewTextView().
		SetText(" (enter) to continue\n (q) to quit\n (ctrl+c) to kill")

	// PAGE INTRO - building
	pageIntro = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(bodyIntro, 0, 4, false).
		AddItem(footerIntro, 0, 1, false)

	// PAGE INTRO - keybind
	pageIntro.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 0 {

		} else if event.Rune() == 13 { // (enter)
			createPkgfile()
			addName()
			pages.SwitchToPage("Name")
		} else if event.Rune() == 113 { // (q)
			deletePkgfile()
			app.Stop()
			fmt.Println("Exiting config mode...")
		}
		return event
	})

	// PAGE NAME - header
	txtName = `
  Enter the package name and then select if you want to manage
  dependencies or the installation process.
`
	headerName = tview.NewTextView().
		SetText(txtName)

	// PAGE NAME - body
	bodyName = tview.NewForm()
	// PAGE NAME - footer
	footerName = tview.NewTextView().
		SetText(" (ctrl+c) to quit")

	// PAGE NAME - building
	pageName = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerName, 0, 1, false).
		AddItem(bodyName, 0, 4, true).
		AddItem(footerName, 0, 1, false)

	// PAGE NAME - keybind
	pageName.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 0 {
		} else {

		}
		return event
	})

	// PAGE DEPENDENCIES INTRO - body
	txtDependenciesIntro = `
  A *basic* dependency is one which is in the default package managers of the most
  used classes of Linux distributions. 

  A *custom* dependency is one in which a custom command need to be provided for 
  installation. This include installing packages for modular programming languages 
  as Python, Go, Node.js, Ruby, etc, as well as building and installing from the source.

	1. Hit (b) or (c) to add a dependency and complete the informations.
	2. Act the button add new to repeat the process.
	3. Proceed until adding all dependencies. 
	4. Conclude adding the last dependency and acting the button conclude.

  If you want to return, hit (r).
`
	bodyDependenciesIntro = tview.NewTextView().
		SetText(txtDependenciesIntro)

	// PAGE DEPENDENCIES INTRO - footer
	footerDependenciesIntro = tview.NewTextView().
		SetText(" (r) to return\n (b) to add a basic dependency\n (c) to add a custom dependency\n (q) to quit\n (ctrl+c) to kill ")

	// PAGE DEPENDENCIES INTRO - building
	pageDependenciesIntro = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(bodyDependenciesIntro, 0, 4, true).
		AddItem(footerDependenciesIntro, 0, 1, false)

	// PAGE DEPENDENCIES INTRO - keybind
	pageDependenciesIntro.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 0 { // (esc)

		} else if event.Rune() == 98 { // (a)
			bodyBasicDependenciesAdd.Clear(true)
			addBasicDependency()
			pages.SwitchToPage("Basic Dependencies Add")
		} else if event.Rune() == 99 { // (c)
			bodyCustomDependenciesAdd.Clear(true)
			addCustomDependency()
			pages.SwitchToPage("Custom Dependencies Add")
		} else if event.Rune() == 114 { // (r)
			bodyName.Clear(true)
			addName()
			pages.SwitchToPage("Name")
		} else if event.Rune() == 113 { // (q)
			deletePkgfile()
			app.Stop()
			fmt.Println("Exiting config mode...")
		}

		return event
	})

	// PAGE BASIC DEPENDENCIES ADD - header
	txtBasicDependenciesAdd = ` 
  Enter the dependency name and its command in each class of Linux 
  distributions...
`
	headerBasicDependenciesAdd = tview.NewTextView().
		SetText(txtBasicDependenciesAdd)

	// PAGE BASIC DEPENDENCIES ADD - body
	bodyBasicDependenciesAdd = tview.NewForm()

	// PAGE BASIC DEPENDENCIES ADD - footer
	footerBasicDependenciesAdd = tview.NewTextView().
		SetText(" (ctrl+c) to kill")

	// PAGE BASIC DEPENDENCIES ADD - building
	pageBasicDependenciesAdd = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerBasicDependenciesAdd, 0, 1, false).
		AddItem(bodyBasicDependenciesAdd, 0, 4, true).
		AddItem(footerBasicDependenciesAdd, 0, 1, false)

	// PAGE BASIC DEPENDENCIES ADD - keybind
	pageBasicDependenciesAdd.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 0 {
		} else {

		}
		return event
	})

	// PAGE CUSTOM DEPENDENCIES ADD - header
	txtCustomDependenciesAdd = ` 
  Enter the dependency name and then the script to install it.

	* the textarea will be read as a Bash script. Thus, be
	  sure to add one command per line.
`
	headerCustomDependenciesAdd = tview.NewTextView().
		SetText(txtCustomDependenciesAdd)

	// PAGE CUSTOM DEPENDENCIES ADD - body
	bodyCustomDependenciesAdd = tview.NewForm()

	// PAGE CUSTOM DEPENDENCIES ADD - footer
	footerCustomDependenciesAdd = tview.NewTextView().
		SetText(" (ctrl+c) to kill")

	// PAGE CUSTOM DEPENDENCIES ADD - building
	pageCustomDependenciesAdd = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerCustomDependenciesAdd, 0, 1, false).
		AddItem(bodyCustomDependenciesAdd, 0, 4, true).
		AddItem(footerCustomDependenciesAdd, 0, 1, false)

	// PAGE CUSTOM INSTALL INTRO - body
	txtCustomInstallIntro = `
  The default install process consists in

	1. copying the contents of the directory where the install script is runned
	   to the path specified when ./configure is runned;
	2. fixing some variables with that defined in the pkgfile;
	3. calling the main function in the .bashrc file.

  In the following you will be able to add custom installation steps to be executed
  before the default installation process.
`
	bodyCustomInstallIntro = tview.NewTextView().
		SetText(txtCustomInstallIntro)

	// PAGE CUSTOM INSTALL INTRO - footer
	footerCustomInstallIntro = tview.NewTextView().
		SetText(" (r) to return\n (a) to add a custom install script\n (q) to quit\n (ctrl+c) to kill")

	// PAGE CUSTOM INSTALL INTRO - building
	pageCustomInstallIntro = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(bodyCustomInstallIntro, 0, 4, true).
		AddItem(footerCustomInstallIntro, 0, 1, false)

	// PAGE CUSTOM INSTALL INTRO - keybind
	pageCustomInstallIntro.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 0 {

		} else if event.Rune() == 97 { // (a)
			bodyCustomInstallAdd.Clear(true)
			createCustomInstall()
			addCustomInstall()
			pages.SwitchToPage("Install Add")
		} else if event.Rune() == 114 { // (r)
			bodyName.Clear(true)
			addName()
			pages.SwitchToPage("Name")
		} else if event.Rune() == 113 { // (q)
			deletePkgfile()
			app.Stop()
			fmt.Println("Exiting config mode...")
		}

		return event
	})

	// PAGE CUSTOM INSTALL ADD - header
	txtCustomInstallAdd = `
  Enter the script containing the custom installation steps.

	* the textarea will be read as a Bash script. Thus, be
	  sure to add one command per line.
`
	headerCustomInstallAdd = tview.NewTextView().
		SetText(txtCustomInstallAdd)

	// PAGE CUSTOM INSTALL ADD - body
	bodyCustomInstallAdd = tview.NewForm()

	// PAGE CUSTOM INSTALL ADD - footer
	footerCustomInstallAdd = tview.NewTextView().
		SetText(" (ctrl+c) to kill")

	// PAGE CUSTOM INSTALL ADD - building
	pageCustomInstallAdd = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerCustomInstallAdd, 0, 1, false).
		AddItem(bodyCustomInstallAdd, 0, 4, true).
		AddItem(footerCustomInstallAdd, 0, 1, false)

	// PAGE CUSTOM INSTALL ADD - keybind
	pageCustomInstallAdd.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 0 {
		} else {

		}
		return event
	})

	// PAGES
	pages = tview.NewPages().
		AddPage("Intro", pageIntro, true, true).
		AddPage("Name", pageName, true, false).
		AddPage("Dependencies Intro", pageDependenciesIntro, true, false).
		AddPage("Basic Dependencies Add", pageBasicDependenciesAdd, true, false).
		AddPage("Custom Dependencies Add", pageCustomDependenciesAdd, true, false).
		AddPage("Install Intro", pageCustomInstallIntro, true, false).
		AddPage("Install Add", pageCustomInstallAdd, true, false)

	// setting the widget "pages" as the root. Enable mouse.
	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}
