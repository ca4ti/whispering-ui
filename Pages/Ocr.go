package Pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"strings"
	"whispering-tiger-ui/Fields"
	"whispering-tiger-ui/Settings"
	"whispering-tiger-ui/Utilities"
	"whispering-tiger-ui/Websocket/Messages"
)

func CreateOcrWindow() fyne.CanvasObject {
	defer Utilities.PanicLogger()

	Fields.Field.OcrLanguageCombo.OnSubmitted = func(value string) {
		for i := 0; i < len(Fields.Field.OcrLanguageCombo.Options); i++ {
			if strings.Contains(strings.ToLower(Fields.Field.OcrLanguageCombo.Options[i]), strings.ToLower(value)) {
				Fields.Field.OcrLanguageCombo.SelectItemByValue(Fields.Field.OcrLanguageCombo.Options[i])
				value = Fields.Field.OcrLanguageCombo.Options[i]
				Fields.Field.OcrLanguageCombo.Text = value
				Fields.Field.OcrLanguageCombo.Entry.CursorColumn = len(Fields.Field.OcrLanguageCombo.Text)
				Fields.Field.OcrLanguageCombo.Refresh()
				break
			}
		}

		valueIso := Messages.OcrLanguagesList.GetCodeByName(value)
		if valueIso == "" {
			valueObj := Fields.Field.SourceLanguageTxtTranslateCombo.GetValueOptionEntryByText(value)
			value = valueObj.Value

			valueIso = Messages.OcrLanguagesList.GetCodeByName(value)
		}

		sendMessage := Fields.SendMessageStruct{
			Type:  "setting_change",
			Name:  "ocr_lang",
			Value: valueIso,
		}
		sendMessage.SendMessage()

		log.Println("ocr Select set to", value)
	}

	container.New(layout.NewMaxLayout())
	ocrLanguageWindowForm := container.New(layout.NewFormLayout(), widget.NewLabel("Text in Image Language:"), Fields.Field.OcrLanguageCombo, widget.NewLabel("Window:"), Fields.Field.OcrWindowCombo)

	ocrSettingsRow := container.New(layout.NewGridLayout(1), ocrLanguageWindowForm)

	ocrButton := widget.NewButtonWithIcon("process window with OCR", theme.ConfirmIcon(), func() {

		ocrLanguageCode := Messages.OcrLanguagesList.GetCodeByName(Fields.Field.OcrLanguageCombo.Text)

		fromLang := ""
		if len(Fields.Field.SourceLanguageTxtTranslateCombo.OptionsTextValue) > 0 {
			fromLang = Messages.InstalledLanguages.GetCodeByName(Fields.Field.SourceLanguageTxtTranslateCombo.GetValueOptionEntryByText(Fields.Field.SourceLanguageTxtTranslateCombo.Text).Value)
		} else {
			fromLang = Messages.InstalledLanguages.GetCodeByName(Fields.Field.SourceLanguageTxtTranslateCombo.Text)
		}
		if fromLang == "" || fromLang == "Auto" {
			fromLang = "auto"
		}
		if fromLang == "auto" {
			guessedSrcLangByOCRLang := ""
			// try to guess the language from the OCR language selection if auto to lessen the language guessing
			guessedSrcLangByOCRLang = Messages.InstalledLanguages.GetCodeByName(Fields.Field.OcrLanguageCombo.Text)
			if guessedSrcLangByOCRLang == "" {
				if Utilities.LanguageMapList.GetName(ocrLanguageCode) != "" {
					guessedSrcLangByOCRLang = ocrLanguageCode
				}
			}
			if guessedSrcLangByOCRLang != "" {
				fromLang = guessedSrcLangByOCRLang
				println("guessedSrcLangByOCRLang", guessedSrcLangByOCRLang)
			}
		}

		toLang := Messages.InstalledLanguages.GetCodeByName(Fields.Field.TargetLanguageTxtTranslateCombo.Text)
		//goland:noinspection GoSnakeCaseUsage
		sendMessage := Fields.SendMessageStruct{
			Type: "ocr_req",
			Value: struct {
				Ocr_lang  string `json:"ocr_lang"`
				From_lang string `json:"from_lang"`
				To_lang   string `json:"to_lang"`
			}{
				Ocr_lang:  ocrLanguageCode,
				From_lang: fromLang,
				To_lang:   toLang,
			},
		}
		sendMessage.SendMessage()
	})
	ocrButton.Importance = widget.HighImportance

	buttonRow := container.NewHBox(layout.NewSpacer(),
		ocrButton,
	)

	switchButton := container.NewCenter(widget.NewButton("<==>", func() {
		sourceLanguage := Fields.Field.SourceLanguageTxtTranslateCombo.Text
		// use last detected language when switching between source and target language
		if strings.HasPrefix(strings.ToLower(sourceLanguage), "auto") && Settings.Config.Last_auto_txt_translate_lang != "" {
			sourceLanguage = Utilities.LanguageMapList.GetName(Settings.Config.Last_auto_txt_translate_lang)
		}

		targetLanguage := Fields.Field.TargetLanguageTxtTranslateCombo.Text
		if targetLanguage == "None" {
			targetLanguage = "Auto"
		}

		Fields.Field.SourceLanguageTxtTranslateCombo.Text = targetLanguage
		Fields.Field.SourceLanguageTxtTranslateCombo.Refresh()
		Fields.Field.TargetLanguageTxtTranslateCombo.Text = sourceLanguage
		Fields.Field.TargetLanguageTxtTranslateCombo.Refresh()

		sourceField := Fields.Field.TranscriptionInput.Text
		targetField := Fields.Field.TranscriptionTranslationInput.Text
		Fields.Field.TranscriptionInput.SetText(targetField)
		Fields.Field.TranscriptionTranslationInput.SetText(sourceField)
	}))

	sourceLanguageForm := container.New(layout.NewFormLayout(), widget.NewLabel("Source Language:"), Fields.Field.SourceLanguageTxtTranslateCombo)
	targetLanguageForm := container.New(layout.NewFormLayout(), widget.NewLabel("Target Language:"), Fields.Field.TargetLanguageTxtTranslateCombo)
	languageRow := container.New(layout.NewGridLayout(2), sourceLanguageForm, targetLanguageForm)

	transcriptionRow := container.New(layout.NewGridLayout(2), Fields.Field.TranscriptionInput, Fields.Field.TranscriptionTranslationInput)

	translateOnlyFunction := func() {
		fromLang := Messages.InstalledLanguages.GetCodeByName(Fields.Field.SourceLanguageTxtTranslateCombo.Text)
		if fromLang == "" {
			fromLang = "auto"
		}
		toLang := Messages.InstalledLanguages.GetCodeByName(Fields.Field.TargetLanguageTxtTranslateCombo.Text)
		//goland:noinspection GoSnakeCaseUsage
		sendMessage := Fields.SendMessageStruct{
			Type: "translate_req",
			Value: struct {
				Text                string `json:"text"`
				From_lang           string `json:"from_lang"`
				To_lang             string `json:"to_lang"`
				Ignore_send_options bool   `json:"ignore_send_options"`
			}{
				Text:                Fields.Field.TranscriptionInput.Text,
				From_lang:           fromLang,
				To_lang:             toLang,
				Ignore_send_options: true,
			},
		}
		sendMessage.SendMessage()
	}
	translateOnlyButton := widget.NewButtonWithIcon("Translate Only", theme.MenuExpandIcon(), translateOnlyFunction)

	ocrContent := container.New(layout.NewVBoxLayout(),
		ocrSettingsRow,
		container.New(layout.NewPaddedLayout(), buttonRow),
		widget.NewSeparator(),
		widget.NewLabel("Text-Translation of OCR Result:"),
		languageRow,
		switchButton,
	)

	mainContent := container.NewBorder(
		container.New(layout.NewVBoxLayout(),
			ocrContent,
		),
		nil, nil, nil,
		container.NewVSplit(
			transcriptionRow,
			container.NewBorder(
				container.NewBorder(
					nil, nil, nil, translateOnlyButton,
				),
				nil, nil, nil, Fields.Field.OcrImageContainer,
			),
		),
	)

	return mainContent
}
