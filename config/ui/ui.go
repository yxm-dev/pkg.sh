package main

import (
	"fmt"
	"io"
	"os"

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
var packageName string
var debian string
var dependencyName string
var commandDistros string

var txtIntro string
var bodyIntro *tview.TextView
var footerIntro *tview.TextView
var pageIntro *tview.Flex

var headerName *tview.TextView
var bodyName *tview.InputField
var footerName *tview.TextView
var pageName *tview.Flex

var txtDependenciesIntro string
var bodyDependenciesIntro *tview.Form
var footerDependenciesIntro *tview.TextView
var pageDependenciesIntro *tview.Flex

var headerDependenciesAdd *tview.TextView
var bodyDependenciesAdd *tview.Form
var pageDependenciesAdd *tview.Flex
var footerDependenciesAdd *tview.TextView

var pages *tview.Pages

// // FUNCTIONS // //
// check if there already exists a file "pkgfile" in the working directory.
// If don't, create it.
func createPkgfile() {
	_, err := os.Stat("pkgfile")
	if err == nil {
		app.Stop()
		fmt.Println("ERROR. There already exists a pkgfile in the working directory.")
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

// print the package name in the file "pkgfile"
func printPackageNamePkgfile() {
	file, err := os.OpenFile("pkgfile", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	packageNamePkgfile := []byte("PKG_name=" + "\"" + packageName + "\"\n\n")
	file.Write(packageNamePkgfile)
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

// create a form of the dependencies names and commands in distros.
// print the result in the files "pkgfile_dependencies" and "pkgfile_dependencies_distros".
func addDependency() {

	var labels []string

	labels = append(labels, "dependency name:")
	labels = append(labels, "Debian based:")
	labels = append(labels, "Arch Linux based:")
	labels = append(labels, "Red Hat based:")
	labels = append(labels, "SUSE based:")

	for i := 0; i < len(labels); i++ {
		bodyDependenciesAdd.AddInputField(labels[i], "", 30, nil, nil).
			SetFieldBackgroundColor(tcell.ColorGray).
			SetFieldTextColor(tcell.ColorWhite).
			SetLabelColor(tcell.ColorWhite)
	}

	bodyDependenciesAdd.AddButton("add other", func() {
		dependencyName = bodyDependenciesAdd.GetFormItemByLabel(labels[0]).(*tview.InputField).GetText()
		commandDistros = bodyDependenciesAdd.GetFormItemByLabel(labels[1]).(*tview.InputField).GetText()

		appendDependenciesPkgfile(dependencyName)
		appendDependenciesDistrosPkgfile(dependencyName, commandDistros)

		pages.SwitchToPage("Dependencies Intro")
	}).
		SetLabelColor(tcell.ColorWhite).
		SetBorder(false)

	bodyDependenciesAdd.AddButton("conclude", func() {
		dependencyName = bodyDependenciesAdd.GetFormItemByLabel(labels[0]).(*tview.InputField).GetText()
		commandDistros = bodyDependenciesAdd.GetFormItemByLabel(labels[1]).(*tview.InputField).GetText()

		appendLastDependencyPkgfile(dependencyName)
		appendDependenciesDistrosPkgfile(dependencyName, commandDistros)

		concludePkgfile()
		deleteConcludePkgfile()
		app.Stop()
		fmt.Println("Pkgfile has been created.")
	}).
		SetLabelColor(tcell.ColorWhite).
		SetBorder(false)
}

// function to pass when the conclude button is pressed.
// append the file "pkgfile" with the other files.
// move the "pkgfile" to the working directory
func concludePkgfile() {
	var pkgfiles = []string{
		"pkgfile_distros",
		"pkgfile_dependencies",
		"pkgfile_dependencies_distros",
	}
	var pkgfile string = "pkgfile"

	target, err := os.OpenFile(pkgfile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	for _, file := range pkgfiles {
		source, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(target, source)
		if err != nil {
			panic(err)
		}

		newLine := []byte("\n\n")
		_, err = target.Write(newLine)
		if err != nil {
			panic(err)
		}

		err = source.Close()
		if err != nil {
			panic(err)
		}
	}

	defer target.Close()
	if err != nil {
		panic(err)
	}

}

// delete just the additional files.
// used when concluding the configuration
func deleteConcludePkgfile() {
	os.Remove("pkgfile_distros")
	if err != nil {
		panic(err)
	}

	_, err := os.Stat("pkgfile_dependencies")
	if err == nil {
		os.Remove("pkgfile_dependencies")
	} else {
		panic(err)
	}

	_, err = os.Stat("pkgfile_dependencies_distros")
	if err == nil {
		os.Remove("pkgfile_dependencies_distros")
	} else {
		panic(err)
	}

}

// delete the additional files and the "pkgfile".
// used when exiting without concluding.

func deletePkgfile() {

	os.Remove("pkgfile")
	if err != nil {
		panic(err)
	}

	deleteConcludePkgfile()
}

// // MAIN // //
func main() {

	// initializing the app
	app = tview.NewApplication()

	// PAGE INTRO - text
	txtIntro = `
  Welcome to the configurarion mode of pkg.sh.

  In the following you will create a pkgfile by providing:
	1. the new package name;
	2. its dependencies;
	3. the corresponding commands in each class of distributions.
	`
	// PAGE INTRO - body
	bodyIntro = tview.NewTextView().
		SetText(txtIntro)

	// PAGE INTRO - footer
	footerIntro = tview.NewTextView().
		SetText(" (enter) to continue\n (esc) to quit\n (ctrl+c) to kill")

	// PAGE INTRO - building
	pageIntro = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(bodyIntro, 0, 4, false).
		AddItem(footerIntro, 0, 1, false)

	// PAGE INTRO - keybind
	pageIntro.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 0 { // (esc)
			app.Stop()
			fmt.Println("Exiting config mode...")
		} else if event.Rune() == 13 { // (enter)
			createPkgfile()
			pages.SwitchToPage("Name")
		}
		return event
	})

	// PAGE NAME - header
	headerName = tview.NewTextView().
		SetText("\n Enter the package name...")

	// PAGE NAME - body
	bodyName = tview.NewInputField().
		SetLabel(" name: ").
		SetLabelColor(tcell.ColorWhite).
		SetFieldWidth(30).
		SetFieldBackgroundColor(tcell.ColorGray).
		SetFieldTextColor(tcell.ColorWhite).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				packageName = bodyName.GetText()
				printPackageNamePkgfile()
				createDistrosInclude()
				createDependenciesPkgfile()
				pages.SwitchToPage("Dependencies Intro")
			}
		})

	// PAGE NAME - footer
	footerName = tview.NewTextView().
		SetText(" (enter) to continue\n (esc) to quit\n (ctrl+c) to kill")

	// PAGE NAME - building
	pageName = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerName, 0, 1, false).
		AddItem(bodyName, 0, 4, true).
		AddItem(footerName, 0, 1, false)

	// PAGE NAME - keybind
	pageName.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 0 { // (esc)
			deletePkgfile()
			app.Stop()
			fmt.Println("Exiting config mode...")
		}
		return event
	})

	// PAGE DEPENDENCIES INTRO - body
	txtDependenciesIntro = `
 1. Hit (a) to add a dependency and complete the informations.
 2. Act the button "add new" to repeat the process.
 3. Proceed until adding all dependencies. 
 4. Conclude adding the last dependency and acting the button "conclude".
`
	bodyDependenciesIntro = tview.NewForm().
		AddTextView("", txtDependenciesIntro, 0, 10, true, true)

	// PAGE DEPENDENCIES INTRO - footer
	footerDependenciesIntro = tview.NewTextView().
		SetText(" (a) to add new dependency\n (esc) to quit\n (ctrl+c) to kill")

	// PAGE DEPENDENCIES INTRO - building
	pageDependenciesIntro = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(bodyDependenciesIntro, 0, 4, true).
		AddItem(footerDependenciesIntro, 0, 1, false)

	// PAGE DEPENDENCIES INTRO - keybind
	pageDependenciesIntro.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 0 { // (esc)
			deletePkgfile()
			app.Stop()
			fmt.Println("Exiting config mode...")
		} else if event.Rune() == 97 { // (a)
			bodyDependenciesAdd.Clear(true)
			addDependency()
			pages.SwitchToPage("Dependencies Add")
		}
		return event
	})

	// PAGE DEPENDENCIES ADD - header
	headerDependenciesAdd = tview.NewTextView().
		SetText("\n Enter the dependency name and its command in the distros...")

	// PAGE DEPENDENCIES ADD - body
	bodyDependenciesAdd = tview.NewForm()

	// PAGE DEPENDENCIES ADD- footer
	footerDependenciesAdd = tview.NewTextView().
		SetText(" (ctrl+c) to kill")

	// PAGE DEPENDENCIES ADD - building
	pageDependenciesAdd = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerDependenciesAdd, 0, 1, false).
		AddItem(bodyDependenciesAdd, 0, 4, true).
		AddItem(footerDependenciesAdd, 0, 1, false)

	// PAGES
	pages = tview.NewPages().
		AddPage("Intro", pageIntro, true, true).
		AddPage("Name", pageName, true, false).
		AddPage("Dependencies Intro", pageDependenciesIntro, true, false).
		AddPage("Dependencies Add", pageDependenciesAdd, true, false)

	// setting the widget "pages" as the root. Enable mouse.
	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}
