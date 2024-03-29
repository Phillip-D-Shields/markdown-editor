package main

import (
	"testing"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
)

func Test_makeUI(t *testing.T) {
	var testCfg config

	edit, preview := testCfg.makeUI()

	test.Type(edit, "kia ora")

	if preview.String() != "kia ora" {
		t.Error("Preview did not match expected value")
	}
}

func Test_runApp(t *testing.T) {
	var testCfg config

	testApp := test.NewApp()
	testWindow := testApp.NewWindow("Test Window")

	edit, preview := testCfg.makeUI()

	testCfg.createMenuItems(testWindow)

	testWindow.SetContent(container.NewHSplit(edit, preview))

	testApp.Run()

	test.Type(edit, "testing 123")
	if preview.String() != "testing 123" {
		t.Error("Preview did not match expected value")
	}
}
