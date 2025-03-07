package Pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"os"
	"path/filepath"
	"time"
	"whispering-tiger-ui/Fields"
	"whispering-tiger-ui/Settings"
	"whispering-tiger-ui/Utilities"
	"whispering-tiger-ui/Websocket/Messages"
)

func CreateSpeechToTextWindow() fyne.CanvasObject {
	defer Utilities.PanicLogger()

	speechLanguageLabel := widget.NewLabel("Speech Language:")

	speechTaskWidgetLabel := widget.NewLabel("Speech Task:")
	var speechTaskWidget fyne.CanvasObject = Fields.Field.TranscriptionTaskCombo
	// disable task for seamless_m4t model, as it always translates to target language (Speech Language)
	if Settings.Config.Stt_type == "seamless_m4t" || Settings.Config.Stt_type == "nemo_canary" {
		speechTaskWidgetLabel.SetText("Target Language:")
		speechTaskWidget = Fields.Field.TranscriptionTargetLanguageCombo
	}
	if Settings.Config.Stt_type == "wav2vec_bert" {
		speechTaskWidgetLabel.SetText("")
		speechTaskWidget.(*widget.Select).Hide()
	}

	languageRow := container.New(layout.NewVBoxLayout(),
		container.New(layout.NewFormLayout(),
			speechLanguageLabel,
			Fields.Field.TranscriptionSpeakerLanguageCombo,

			speechTaskWidgetLabel,
			speechTaskWidget,
		),
	)

	transcriptionRow := container.New(
		layout.NewGridLayout(2),
		container.NewBorder(nil, Fields.Field.TranscriptionInputHint, nil, nil, Fields.Field.TranscriptionInput),
		container.NewBorder(nil, Fields.Field.TranscriptionTranslationInputHint, nil, nil, Fields.Field.TranscriptionTranslationInput),
	)

	beginLine := canvas.NewHorizontalGradient(&color.NRGBA{R: 198, G: 123, B: 0, A: 255}, &color.NRGBA{R: 198, G: 123, B: 0, A: 0})
	beginLine.Resize(fyne.NewSize(Fields.Field.SttEnabled.Size().Width, 2))

	// quick options row
	quickOptionsRow := container.New(
		layout.NewVBoxLayout(),
		Fields.Field.SttEnabled,
		container.NewGridWithColumns(2, beginLine),
		Fields.Field.TextTranslateEnabled,
		Fields.Field.TtsEnabled,
		container.NewHBox(
			container.NewBorder(nil, nil, nil, Fields.Field.OscLimitHint, Fields.Field.OscEnabled),
		),
	)

	// main layout
	leftVerticalLayout := container.NewBorder(
		container.New(layout.NewVBoxLayout(),
			languageRow,
		),
		nil, nil, nil,
		container.NewVSplit(
			transcriptionRow,
			container.New(layout.NewHBoxLayout(), quickOptionsRow),
		),
	)

	Fields.Field.ProcessingStatus = widget.NewProgressBarInfinite()

	// whisper Result list
	Fields.Field.WhisperResultList = widget.NewListWithData(Fields.DataBindings.WhisperResultsDataBinding,
		func() fyne.CanvasObject {
			return container.New(layout.NewGridLayout(1),
				container.NewBorder(
					nil,
					nil,
					nil,
					widget.NewLabelWithStyle("[ResultLang]", fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
					widget.NewLabelWithStyle("TranslateResult", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				),
				container.NewBorder(
					nil,
					nil,
					nil,
					widget.NewLabelWithStyle("[ResultLang]", fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
					widget.NewLabel("Transcription"),
				),
			)
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			value := i.(binding.Untyped)
			whisperMessage, _ := value.Get()

			result := whisperMessage.(Messages.WhisperResult)

			translateResultBind := binding.NewString()
			translateResultBind.Set(result.TxtTranslation)

			translateResultLanguageBind := binding.NewString()
			translateResultLanguageBind.Set("[" + result.TxtTranslationTarget + "]")

			originalTranscriptBind := binding.NewString()
			originalTranscriptBind.Set(result.Text)

			originalTranscriptLanguageBind := binding.NewString()
			originalTranscriptLanguageBind.Set("[" + result.Language + "]")

			// get all template elements
			mainContainer := o.(*fyne.Container)
			finalTranslationContainer := mainContainer.Objects[0].(*fyne.Container)
			originalTranscriptionContainer := mainContainer.Objects[1].(*fyne.Container)

			translateResultLabel := finalTranslationContainer.Objects[0].(*widget.Label)
			translateResultLabel.Wrapping = fyne.TextWrapWord
			translateResultLanguageLabel := finalTranslationContainer.Objects[1].(*widget.Label)

			originalTranscriptionLabel := originalTranscriptionContainer.Objects[0].(*widget.Label)
			originalTranscriptionLabel.Wrapping = fyne.TextWrapWord
			originalTranscriptionLanguageLabel := originalTranscriptionContainer.Objects[1].(*widget.Label)

			// bind data to elements if no translation is generated (sets transcription to top label)
			if result.TxtTranslation == "" {
				translateResultLabel.Bind(originalTranscriptBind)
				translateResultLanguageLabel.Bind(originalTranscriptLanguageBind)

				originalTranscriptionLabel.SetText("")
				originalTranscriptionLanguageLabel.SetText("")
			} else { // bind data to elements if translation was generated
				translateResultLabel.Bind(translateResultBind)
				translateResultLanguageLabel.Bind(translateResultLanguageBind)

				originalTranscriptionLabel.Bind(originalTranscriptBind)
				originalTranscriptionLanguageLabel.Bind(originalTranscriptLanguageBind)
			}

			// resize
			//mainContainer.Resize(fyne.NewSize(mainContainer.Size().Width, translateResultLabel.Size().Height+originalTranscriptionLabel.Size().Height+10))
		})

	Fields.Field.WhisperResultList.OnSelected = func(id widget.ListItemID) {
		whisperMessage, _ := Fields.DataBindings.WhisperResultsDataBinding.GetValue(id)

		result := whisperMessage.(Messages.WhisperResult)

		Fields.Field.TranscriptionInput.SetText(result.Text)
		if result.TxtTranslation != "" {
			Fields.Field.TranscriptionTranslationInput.SetText(result.TxtTranslation)
		} else {
			Fields.Field.TranscriptionTranslationInput.SetText(result.Text)
		}

		go func() {
			time.Sleep(200 * time.Millisecond)
			Fields.Field.WhisperResultList.Unselect(id)
		}()
	}

	if !Settings.Config.Realtime {
		Fields.Field.RealtimeResultLabel.Hide()
	}
	realtimeWhisperResultBlock := container.NewBorder(
		nil, container.NewVBox(widget.NewSeparator(), widget.NewSeparator()), nil, nil,
		Fields.Field.RealtimeResultLabel,
	)

	saveCsvButton := widget.NewButton("Save CSV", func() {
		dialogSize := fyne.CurrentApp().Driver().AllWindows()[0].Canvas().Size()
		dialogSize.Height = dialogSize.Height - 80
		dialogSize.Width = dialogSize.Width - 80

		saveStartingPath := fyne.CurrentApp().Preferences().StringWithFallback("LastCSVTranscriptionSavePath", "")

		fileDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err == nil && writer != nil {
				filePath := writer.URI().Path()
				sendMessage := Fields.SendMessageStruct{
					Type:  "save_transcription",
					Value: filePath,
				}
				sendMessage.SendMessage()
				fyne.CurrentApp().Preferences().SetString("LastCSVTranscriptionSavePath", filepath.Dir(filePath))
			}
		}, fyne.CurrentApp().Driver().AllWindows()[0])

		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
		fileDialog.Resize(dialogSize)

		if saveStartingPath != "" {
			// check if folder exists
			folderExists := false
			if _, err := os.Stat(saveStartingPath); !os.IsNotExist(err) {
				folderExists = true
			}
			if folderExists {
				fileURI := storage.NewFileURI(saveStartingPath)
				fileLister, _ := storage.ListerForURI(fileURI)

				fileDialog.SetLocation(fileLister)
			}
		}

		fileDialog.SetFileName("transcription_" + time.Now().Format("2006-01-02_15-04-05") + ".csv")

		fileDialog.Show()
	})
	lastResultLine := container.NewBorder(nil, nil, saveCsvButton, nil, Fields.Field.ProcessingStatus)

	whisperResultContainer := container.NewMax(
		container.NewBorder(
			realtimeWhisperResultBlock, lastResultLine, nil, nil,
			Fields.Field.WhisperResultList,
		),
	)

	mainContent := container.NewHSplit(
		leftVerticalLayout,
		container.NewMax(
			whisperResultContainer,
		),
	)

	Fields.Field.ProcessingStatus.Stop()

	mainContent.SetOffset(0.6)

	return mainContent
}
