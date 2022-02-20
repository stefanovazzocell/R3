package main

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/stefanovazzocell/R3/clientlib"
)

var APIEndpoint string = "https://api.rkt.one"

var ttlTOint map[string]int = map[string]int{
	"10min": 600, "1h": 3600, "12h": 43200, "1d": 86400, "2d": 172800, "1w": 604800,
}

type PickedFile struct {
	Name string
	Path string
	Data []byte
}

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Rocket Share")

	ShareID := widget.NewEntry()
	ShareID.SetText(clientlib.GenID())
	ShareID.Validator = validation.NewRegexp("^[A-Za-z0-9-._~:/?#\\[\\]@!\\$&'\\(\\)\\*\\+,;=]{3,}$", "Please enter a string with letters and numbers, only a limited number of special characters are supported")

	ShareTTL := widget.NewSelectEntry([]string{"10min", "1h", "12h", "1d", "2d", "1w"})
	ShareTTL.SetText("10min")

	ShareViews := widget.NewEntry()
	ShareViews.Validator = validation.NewRegexp("^1?[0-9]{1,6}$", "Please enter a valid number from 1 to 1'000'000")
	ShareViews.SetText("10")

	SharePassword := widget.NewEntry()

	ShareCreationMessage := widget.NewLabel("")
	ShareCreationProgress := widget.NewProgressBar()
	ShareCreationDialog := dialog.NewCustom("Creating Share", "", container.NewGridWithRows(
		2,
		ShareCreationMessage,
		ShareCreationProgress,
	), w)

	ShareURL := widget.NewEntry()
	ShareURL.Validator = validation.NewRegexp("^(([^:/?#]+):)?(//([^/?#]*))?([^?#]*)(\\?([^#]*))?(#(.*))?", "Please enter a valid URL")

	ShareMessage := widget.NewMultiLineEntry()
	ShareMessage.Validator = validation.NewRegexp(".{1,}$", "Please enter a message")
	ShareMessage.SetPlaceHolder("Your message")

	ShareFiles := []PickedFile{}
	var ShareFilesList *widget.List
	ShareFilesList = widget.NewList(func() int {
		return len(ShareFiles)
	}, func() fyne.CanvasObject {
		return container.NewGridWithColumns(2, widget.NewLabel(""), widget.NewButton("Remove", func() {}))
	}, func(lii widget.ListItemID, co fyne.CanvasObject) {
		// co.(*widget.Label).SetText(ShareFiles[lii].Name)
		co.(*fyne.Container).Objects[0].(*widget.Label).SetText(fmt.Sprintf("%s (%d MB)", ShareFiles[lii].Name, len(ShareFiles[lii].Data)/1024/1024))
		co.(*fyne.Container).Objects[0].(*widget.Label).Wrapping = fyne.TextTruncate
		co.(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
			ShareFiles = append(ShareFiles[:lii], ShareFiles[lii+1:]...)
			ShareFilesList.Refresh()
		}
	})
	ShareFilesAdd := widget.NewButton("Add", func() {
		fileOpen(w, func(name string, b []byte, path string, err error) {
			if err == nil {
				ShareFiles = append(ShareFiles, PickedFile{
					Name: name,
					Path: path,
					Data: b,
				})
				ShareFilesList.Refresh()
			}
		})
	})

	ShareDataTab := container.NewAppTabs(
		container.NewTabItem("URL", container.NewVBox(
			widget.NewLabel("Enter your URL here:"),
			ShareURL,
		)),
		container.NewTabItem("Message", ShareMessage),
		container.NewTabItem("Files", container.NewBorder(nil, ShareFilesAdd, nil, nil, ShareFilesList)),
	)

	ShareCreateBtn := widget.NewButton("Create", func() {
		ShareCreationMessage.SetText("Validating Data...")
		ShareCreationDialog.SetDismissText("Loading...")
		ShareCreationProgress.SetValue(0)
		ShareCreationDialog.Show()
		// Prepare variables for encryption and prepare hash/key
		var kc chan []byte = make(chan []byte, 1)
		var sc chan []byte = make(chan []byte, 1)
		var dc chan string = make(chan string, 1)
		var hc chan string = make(chan string, 1)
		var phc chan string = make(chan string, 1)
		var ehc chan string = make(chan string, 1)
		go clientlib.GenHash(ShareID.Text, hc)
		go clientlib.KeyDerivation(ShareID.Text, kc, sc)
		defer close(kc)
		defer close(sc)
		defer close(dc)
		defer close(hc)
		defer close(phc)
		defer close(ehc)
		ehc <- ""
		phc <- ""
		ShareCreationProgress.SetValue(0.05)
		// Perform initial validation
		ttl := ttlTOint[ShareTTL.Text]
		if ShareID.Validate() != nil {
			ShareCreationDialog.Hide()
			dialog.NewInformation("Error creating share", "Check your Share ID", w).Show()
			<-hc
			<-kc
			<-sc
			return
		}
		if ttl < 600 {
			ShareCreationDialog.Hide()
			dialog.NewInformation("Error creating share", "Check your Share expiration", w).Show()
			<-hc
			<-kc
			<-sc
			return
		}
		hits, err := strconv.ParseInt(ShareViews.Text, 10, 0)
		if ShareViews.Validate() != nil || err != nil {
			ShareCreationDialog.Hide()
			dialog.NewInformation("Error creating share", "Check your Share TTL", w).Show()
			<-hc
			<-kc
			<-sc
			return
		}
		if ShareDataTab.Selected().Text == "URL" && ShareURL.Validate() != nil {
			ShareCreationDialog.Hide()
			dialog.NewInformation("Error creating share", "Check your URL", w).Show()
			<-hc
			<-kc
			<-sc
			return
		}
		if ShareDataTab.Selected().Text == "Message" && ShareMessage.Validate() != nil {
			ShareCreationDialog.Hide()
			dialog.NewInformation("Error creating share", "Check your Message", w).Show()
			<-hc
			<-kc
			<-sc
			return
		}
		if ShareDataTab.Selected().Text == "Files" && len(ShareFiles) <= 0 {
			ShareCreationDialog.Hide()
			dialog.NewInformation("Error creating share", "Please add a file", w).Show()
			<-hc
			<-kc
			<-sc
			return
		}
		// Setup encryption
		ShareCreationProgress.SetValue(0.15)
		ShareCreationMessage.SetText("Packaging Share...")
		var data []byte
		if ShareDataTab.Selected().Text == "URL" {
			data = clientlib.ShareEncode("url", []byte(ShareURL.Text))
		} else if ShareDataTab.Selected().Text == "Message" {
			data = clientlib.ShareEncode("text", []byte(ShareMessage.Text))
		} else {
			for i := 0; i < len(ShareFiles); i++ {
				var file clientlib.File
				file.Load(ShareFiles[i].Path, ShareFiles[i].Data)
				fmt.Println(len(file.Encode()))
				data = append(data, file.Encode()...)
			}
			fmt.Println(len(data))
			data = clientlib.ShareEncode("file", data)
			fmt.Println(len(data))
		}
		if len(data) > 10485760 {
			ShareCreationDialog.Hide()
			dialog.NewInformation("Error creating share", "Your share is too large!", w).Show()
			<-hc
			<-kc
			<-sc
			return
		} else if len(data) > 10240 && ttl > 3600 {
			ttl = 3600
		}
		clientlib.Encrypt(data, kc, sc, dc)
		ShareCreationProgress.SetValue(0.50)
		ShareCreationMessage.SetText("Encrypting and uploading...")
		apiResp, err := clientlib.EditLink(APIEndpoint, hc, phc, false, dc, ttl, int(hits), ehc)
		if err != nil {
			ShareCreationDialog.Hide()
			dialog.NewInformation("Error creating share", fmt.Sprintf("Client Error: %v", err), w).Show()
		} else if !apiResp.Success {
			ShareCreationDialog.Hide()
			dialog.NewInformation("Error creating share", fmt.Sprintf("Server Error: %v", apiResp.Err), w).Show()
		} else {
			w.Clipboard().SetContent(fmt.Sprintf("https://rkt.one/#%s", ShareID.Text))
			ShareCreationProgress.SetValue(1)
			ShareCreationMessage.SetText("Success!\nYour link is copied to the clipboard")
			time.Sleep(time.Second * 2)
			ShareID.SetText(clientlib.GenID())
			ShareCreationDialog.SetDismissText("Done")
			ShareCreationDialog.Hide()
		}
	})

	w.SetContent(
		container.NewGridWithColumns(2,
			ShareDataTab,
			container.NewBorder(nil, ShareCreateBtn, nil, nil, container.New(
				layout.NewFormLayout(),
				widget.NewLabel("https://rkt.one/#"),
				ShareID,
				layout.NewSpacer(),
				widget.NewButton("Random", func() { ShareID.SetText(clientlib.GenID()) }),
				widget.NewLabel("Expiration timer"),
				ShareTTL,
				widget.NewLabel("Views limit"),
				ShareViews,
				widget.NewLabel("Share Password"),
				SharePassword,
			)),
		),
	)

	if !fyne.CurrentDevice().IsMobile() {
		w.CenterOnScreen()
		w.Resize(fyne.Size{Width: 700, Height: 500})
	}

	w.ShowAndRun()
}

func fileOpen(w fyne.Window, callback func(name string, b []byte, path string, err error)) {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if reader == nil {
			return
		}

		b, err := io.ReadAll(reader)
		callback(reader.URI().Name(), b, reader.URI().Path(), err)
	}, w)
	fd.Show()
}
