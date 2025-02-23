package Fields

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/fyne-io/terminal"
	"image/color"
	"log"
	"strings"
	"whispering-tiger-ui/CustomWidget"
	"whispering-tiger-ui/Utilities"
)

const SttTextTranslateLabelConst = "Automatic Text-Translate from %s to %s"
const OscLimitLabelConst = "[%d / %d]"

var OscLimitHintUpdateFunc = func() {}

var Field = struct {
	RealtimeResultLabel               *widget.Label // only displayed if realtime is enabled
	ProcessingStatus                  *widget.ProgressBarInfinite
	WhisperResultList                 *widget.List
	TranscriptionTaskCombo            *widget.Select
	TranscriptionSpeakerLanguageCombo *CustomWidget.CompletionEntry
	TranscriptionTargetLanguageCombo  *CustomWidget.CompletionEntry
	TranscriptionInput                *CustomWidget.EntryWithPopupMenu
	TranscriptionInputHint            *canvas.Text
	TranscriptionTranslationInput     *CustomWidget.EntryWithPopupMenu
	TranscriptionTranslationInputHint *canvas.Text
	SourceLanguageCombo               *CustomWidget.CompletionEntry
	TargetLanguageCombo               *CustomWidget.CompletionEntry
	SourceLanguageTxtTranslateCombo   *CustomWidget.CompletionEntry // used in OCR tab
	TargetLanguageTxtTranslateCombo   *CustomWidget.CompletionEntry // used in OCR tab
	TtsModelCombo                     *widget.Select
	TtsVoiceCombo                     *widget.Select
	TextTranslateEnabled              *widget.Check
	SttEnabled                        *widget.Check
	TtsEnabled                        *widget.Check
	OscEnabled                        *widget.Check
	OscLimitHint                      *canvas.Text
	OcrLanguageCombo                  *CustomWidget.CompletionEntry
	OcrWindowCombo                    *CustomWidget.TappableSelect
	OcrImageContainer                 *fyne.Container
	LogText                           *terminal.Terminal
}{
	RealtimeResultLabel: widget.NewLabelWithData(DataBindings.WhisperResultIntermediateResult),
	ProcessingStatus:    nil,
	WhisperResultList:   nil,
	TranscriptionTaskCombo: widget.NewSelect([]string{"transcribe", "translate (to English)"}, func(value string) {
		switch value {
		case "transcribe":
			value = "transcribe"
		case "translate (to English)":
			value = "translate"
		}
		sendMessage := SendMessageStruct{
			Type:  "setting_change",
			Name:  "whisper_task",
			Value: value,
		}
		sendMessage.SendMessage()
	}),
	TranscriptionSpeakerLanguageCombo: CustomWidget.NewCompletionEntry([]string{"Auto"}),
	TranscriptionTargetLanguageCombo:  CustomWidget.NewCompletionEntry([]string{}),
	TranscriptionInput: func() *CustomWidget.EntryWithPopupMenu {
		entry := CustomWidget.NewMultiLineEntry()
		entry.Wrapping = fyne.TextWrapWord

		entry.AddAdditionalMenuItem(fyne.NewMenuItem("Send to Text-to-Speech", func() {
			valueData := struct {
				Text     string `json:"text"`
				ToDevice bool   `json:"to_device"`
				Download bool   `json:"download"`
			}{
				Text:     entry.Text,
				ToDevice: true,
				Download: false,
			}
			sendMessage := SendMessageStruct{
				Type:  "tts_req",
				Value: valueData,
			}
			sendMessage.SendMessage()
		}))
		entry.AddAdditionalMenuItem(fyne.NewMenuItem("Send to OSC (VRChat)", func() {
			sendMessage := SendMessageStruct{
				Type: "send_osc",
				//Value: &entry.Text,
				Value: struct {
					Text *string `json:"text"`
				}{
					Text: &entry.Text,
				},
			}
			sendMessage.SendMessage()
		}))
		return entry
	}(),
	TranscriptionInputHint: canvas.NewText("0", color.NRGBA{R: 0xb2, G: 0xb2, B: 0xb2, A: 0xff}),
	TranscriptionTranslationInput: func() *CustomWidget.EntryWithPopupMenu {
		entry := CustomWidget.NewMultiLineEntry()
		entry.Wrapping = fyne.TextWrapWord

		entry.AddAdditionalMenuItem(fyne.NewMenuItem("Send to Text-to-Speech", func() {
			valueData := struct {
				Text     string `json:"text"`
				ToDevice bool   `json:"to_device"`
				Download bool   `json:"download"`
			}{
				Text:     entry.Text,
				ToDevice: true,
				Download: false,
			}
			sendMessage := SendMessageStruct{
				Type:  "tts_req",
				Value: valueData,
			}
			sendMessage.SendMessage()
		}))
		entry.AddAdditionalMenuItem(fyne.NewMenuItem("Send to OSC (VRChat)", func() {

			sendMessage := SendMessageStruct{
				Type: "send_osc",
				//Value: &entry.Text,
				Value: struct {
					Text *string `json:"text"`
				}{
					Text: &entry.Text,
				},
			}
			sendMessage.SendMessage()
		}))
		return entry
	}(),
	TranscriptionTranslationInputHint: canvas.NewText("0", color.NRGBA{R: 0xb2, G: 0xb2, B: 0xb2, A: 0xff}),
	SourceLanguageCombo:               CustomWidget.NewCompletionEntry([]string{"Auto"}),
	TargetLanguageCombo:               CustomWidget.NewCompletionEntry([]string{"None"}),
	SourceLanguageTxtTranslateCombo:   CustomWidget.NewCompletionEntry([]string{"Auto"}),
	TargetLanguageTxtTranslateCombo:   CustomWidget.NewCompletionEntry([]string{"None"}),
	TtsModelCombo: widget.NewSelect([]string{}, func(value string) {
		sendMessage := SendMessageStruct{
			Type:  "setting_change",
			Name:  "tts_model",
			Value: value,
		}
		sendMessage.SendMessage()

		log.Println("Select set to", value)
	}),
	TtsVoiceCombo: widget.NewSelect([]string{}, func(value string) {

		sendMessage := SendMessageStruct{
			Type:  "setting_change",
			Name:  "tts_voice",
			Value: value,
		}
		sendMessage.SendMessage()

		log.Println("Select set to", value)
	}),
	TextTranslateEnabled: widget.NewCheckWithData(fmt.Sprintf(SttTextTranslateLabelConst, "?", "?"), DataBindings.TextTranslateEnabledDataBinding),
	SttEnabled:           widget.NewCheckWithData("Speech-to-Text Enabled", DataBindings.SpeechToTextEnabledDataBinding),
	TtsEnabled:           widget.NewCheckWithData("Automatic Text-to-Speech", DataBindings.TextToSpeechEnabledDataBinding),
	OscEnabled:           widget.NewCheckWithData("Automatic OSC (VRChat)", DataBindings.OSCEnabledDataBinding),
	OscLimitHint:         canvas.NewText(fmt.Sprintf(OscLimitLabelConst, 0, 0), color.NRGBA{R: 0xb2, G: 0xb2, B: 0xb2, A: 0xff}),

	OcrLanguageCombo: CustomWidget.NewCompletionEntry([]string{}),
	OcrWindowCombo: CustomWidget.NewSelect([]string{}, func(value string) {
		sendMessage := SendMessageStruct{
			Type:  "setting_change",
			Name:  "ocr_window_name",
			Value: value,
		}
		sendMessage.SendMessage()
	}),
	OcrImageContainer: container.NewMax(),
	LogText:           terminal.New(),
}

func updateCompletionEntryBasedOnValue(completionEntryWidget *CustomWidget.CompletionEntry, value string) string {
	foundEntry := false
	for _, option := range completionEntryWidget.Options {
		if strings.HasPrefix(strings.ToLower(option), strings.ToLower(value)) {
			completionEntryWidget.SelectItemByValue(option)
			value = option
			completionEntryWidget.Text = value
			completionEntryWidget.Entry.CursorColumn = len(completionEntryWidget.Text)
			completionEntryWidget.Refresh()
			foundEntry = true
			break
		}
	}
	if !foundEntry {
		for _, option := range completionEntryWidget.Options {
			if strings.Contains(strings.ToLower(option), strings.ToLower(value)) {
				completionEntryWidget.SelectItemByValue(option)
				value = option
				completionEntryWidget.Text = value
				completionEntryWidget.Entry.CursorColumn = len(completionEntryWidget.Text)
				completionEntryWidget.Refresh()
				foundEntry = true
				break
			}
		}
	}
	return value
}

func init() {
	defer Utilities.PanicLogger()

	Field.RealtimeResultLabel.Wrapping = fyne.TextWrapWord
	Field.RealtimeResultLabel.TextStyle = fyne.TextStyle{Italic: true}

	Field.SourceLanguageCombo.OptionsTextValue = []CustomWidget.TextValueOption{{
		Text:  "Auto",
		Value: "auto",
	}}
	Field.SourceLanguageCombo.ShowAllEntryText = "... show all"
	Field.SourceLanguageCombo.Entry.PlaceHolder = "Select source language"
	Field.SourceLanguageCombo.OnChanged = func(value string) {
		// filter out the values of Options that do not contain the value
		var filteredValues []string
		for i := 0; i < len(Field.SourceLanguageCombo.Options); i++ {
			if len(Field.SourceLanguageCombo.Options) > i && strings.Contains(strings.ToLower(Field.SourceLanguageCombo.Options[i]), strings.ToLower(value)) {
				filteredValues = append(filteredValues, Field.SourceLanguageCombo.Options[i])
			}
		}

		Field.SourceLanguageCombo.SetOptionsFilter(filteredValues)
		Field.SourceLanguageCombo.ShowCompletion()
	}
	Field.SourceLanguageCombo.OnSubmitted = func(value string) {
		value = updateCompletionEntryBasedOnValue(Field.SourceLanguageCombo, value)

		valueObj := Field.SourceLanguageCombo.GetValueOptionEntryByText(value)

		println("Submitted TargetLanguageCombo", valueObj.Value)

		sendMessage := SendMessageStruct{
			Type:  "setting_change",
			Name:  "src_lang",
			Value: valueObj.Value,
		}
		sendMessage.SendMessage()

		Field.TextTranslateEnabled.Text = fmt.Sprintf(SttTextTranslateLabelConst, valueObj.Value, Field.TargetLanguageCombo.Text)
		Field.TextTranslateEnabled.Refresh()

		log.Println("Select set to", valueObj.Value)
	}

	Field.TargetLanguageCombo.ShowAllEntryText = "... show all"
	Field.TargetLanguageCombo.Entry.PlaceHolder = "Select target language"
	Field.TargetLanguageCombo.OnChanged = func(value string) {
		// filter out the values of Options that do not contain the value
		var filteredValues []string
		for i := 0; i < len(Field.TargetLanguageCombo.Options); i++ {
			if len(Field.TargetLanguageCombo.Options) > i && strings.Contains(strings.ToLower(Field.TargetLanguageCombo.Options[i]), strings.ToLower(value)) {
				filteredValues = append(filteredValues, Field.TargetLanguageCombo.Options[i])
			}
		}

		Field.TargetLanguageCombo.SetOptionsFilter(filteredValues)
		Field.TargetLanguageCombo.ShowCompletion()
	}
	Field.TargetLanguageCombo.OnSubmitted = func(value string) {
		// check if value is not in Options
		value = updateCompletionEntryBasedOnValue(Field.TargetLanguageCombo, value)

		sendMessage := SendMessageStruct{
			Type:  "setting_change",
			Name:  "trg_lang",
			Value: value,
		}
		sendMessage.SendMessage()

		Field.TextTranslateEnabled.Text = fmt.Sprintf(SttTextTranslateLabelConst, Field.SourceLanguageCombo.GetValueOptionEntryByText(Field.SourceLanguageCombo.Text).Value, value)
		Field.TextTranslateEnabled.Refresh()

		log.Println("Select set to", value)
	}

	Field.SourceLanguageTxtTranslateCombo.ShowAllEntryText = "... show all"
	Field.SourceLanguageTxtTranslateCombo.Entry.PlaceHolder = "Select source language"
	Field.SourceLanguageTxtTranslateCombo.OnChanged = func(value string) {
		// filter out the values of Options that do not contain the value
		var filteredValues []string
		for i := 0; i < len(Field.SourceLanguageTxtTranslateCombo.Options); i++ {
			if len(Field.SourceLanguageTxtTranslateCombo.Options) > i && strings.Contains(strings.ToLower(Field.SourceLanguageTxtTranslateCombo.Options[i]), strings.ToLower(value)) {
				filteredValues = append(filteredValues, Field.SourceLanguageTxtTranslateCombo.Options[i])
			}
		}

		Field.SourceLanguageTxtTranslateCombo.SetOptionsFilter(filteredValues)
		Field.SourceLanguageTxtTranslateCombo.ShowCompletion()
	}
	Field.SourceLanguageTxtTranslateCombo.OnSubmitted = func(value string) {
		// check if value is not in Options
		value = updateCompletionEntryBasedOnValue(Field.SourceLanguageTxtTranslateCombo, value)

		sendMessage := SendMessageStruct{
			Type:  "setting_change",
			Name:  "ocr_txt_src_lang",
			Value: value,
		}
		sendMessage.SendMessage()

		log.Println("Select set to", value)
	}

	Field.TargetLanguageTxtTranslateCombo.ShowAllEntryText = "... show all"
	Field.TargetLanguageTxtTranslateCombo.Entry.PlaceHolder = "Select target language"
	Field.TargetLanguageTxtTranslateCombo.OnChanged = func(value string) {
		// filter out the values of Options that do not contain the value
		var filteredValues []string
		for i := 0; i < len(Field.TargetLanguageTxtTranslateCombo.Options); i++ {
			if len(Field.TargetLanguageTxtTranslateCombo.Options) > i && strings.Contains(strings.ToLower(Field.TargetLanguageTxtTranslateCombo.Options[i]), strings.ToLower(value)) {
				filteredValues = append(filteredValues, Field.TargetLanguageTxtTranslateCombo.Options[i])
			}
		}
		Field.TargetLanguageTxtTranslateCombo.SetOptionsFilter(filteredValues)
		Field.TargetLanguageTxtTranslateCombo.ShowCompletion()
	}
	Field.TargetLanguageTxtTranslateCombo.OnSubmitted = func(value string) {
		// check if value is not in Options
		value = updateCompletionEntryBasedOnValue(Field.TargetLanguageTxtTranslateCombo, value)

		sendMessage := SendMessageStruct{
			Type:  "setting_change",
			Name:  "ocr_txt_trg_lang",
			Value: value,
		}
		sendMessage.SendMessage()

		log.Println("Select set to", value)
	}

	Field.TranscriptionSpeakerLanguageCombo.ShowAllEntryText = "... show all"
	Field.TranscriptionSpeakerLanguageCombo.Entry.PlaceHolder = "Select a language"
	Field.TranscriptionSpeakerLanguageCombo.OnChanged = func(value string) {
		// filter out the values of Options that do not contain the value
		var filteredValues []string
		for i := 0; i < len(Field.TranscriptionSpeakerLanguageCombo.Options); i++ {
			if len(Field.TranscriptionSpeakerLanguageCombo.Options) > i && strings.Contains(strings.ToLower(Field.TranscriptionSpeakerLanguageCombo.Options[i]), strings.ToLower(value)) {
				filteredValues = append(filteredValues, Field.TranscriptionSpeakerLanguageCombo.Options[i])
			}
		}

		Field.TranscriptionSpeakerLanguageCombo.SetOptionsFilter(filteredValues)
		Field.TranscriptionSpeakerLanguageCombo.ShowCompletion()
	}
	Field.TranscriptionSpeakerLanguageCombo.OnSubmitted = func(value string) {
		// check if value is not in Options
		value = updateCompletionEntryBasedOnValue(Field.TranscriptionSpeakerLanguageCombo, value)

		sendMessage := SendMessageStruct{
			Type:  "setting_change",
			Name:  "current_language",
			Value: value,
		}
		sendMessage.SendMessage()
	}

	Field.TranscriptionTargetLanguageCombo.ShowAllEntryText = "... show all"
	Field.TranscriptionTargetLanguageCombo.Entry.PlaceHolder = "Select a language"
	Field.TranscriptionTargetLanguageCombo.OnChanged = func(value string) {
		// filter out the values of Options that do not contain the value
		var filteredValues []string
		for i := 0; i < len(Field.TranscriptionTargetLanguageCombo.Options); i++ {
			if len(Field.TranscriptionTargetLanguageCombo.Options) > i && strings.Contains(strings.ToLower(Field.TranscriptionTargetLanguageCombo.Options[i]), strings.ToLower(value)) {
				filteredValues = append(filteredValues, Field.TranscriptionTargetLanguageCombo.Options[i])
			}
		}

		Field.TranscriptionTargetLanguageCombo.SetOptionsFilter(filteredValues)
		Field.TranscriptionTargetLanguageCombo.ShowCompletion()
	}
	Field.TranscriptionTargetLanguageCombo.OnSubmitted = func(value string) {
		// check if value is not in Options
		value = updateCompletionEntryBasedOnValue(Field.TranscriptionTargetLanguageCombo, value)

		sendMessage := SendMessageStruct{
			Type:  "setting_change",
			Name:  "target_language",
			Value: value,
		}
		sendMessage.SendMessage()
	}

	Field.TextTranslateEnabled.OnChanged = func(value bool) {
		sendMessage := SendMessageStruct{
			Type:  "setting_change",
			Name:  "txt_translate",
			Value: value,
		}
		sendMessage.SendMessage()
	}

	Field.SttEnabled.OnChanged = func(value bool) {
		sendMessage := SendMessageStruct{
			Type:  "setting_change",
			Name:  "stt_enabled",
			Value: value,
		}
		sendMessage.SendMessage()
	}
	Field.TtsEnabled.OnChanged = func(value bool) {
		sendMessage := SendMessageStruct{
			Type:  "setting_change",
			Name:  "tts_answer",
			Value: value,
		}
		sendMessage.SendMessage()
	}
	Field.OscEnabled.OnChanged = func(value bool) {
		if value {
			Field.OscLimitHint.Show()
		} else {
			Field.OscLimitHint.Hide()
		}
		sendMessage := SendMessageStruct{
			Type:  "setting_change",
			Name:  "osc_auto_processing_enabled",
			Value: value,
		}
		sendMessage.SendMessage()
	}

	Field.OcrLanguageCombo.ShowAllEntryText = "... show all"
	Field.OcrLanguageCombo.Entry.PlaceHolder = "Select language in image language"
	Field.OcrLanguageCombo.OnChanged = func(value string) {
		// filter out the values of Field.TargetLanguageTxtTranslateCombo.Options that do not contain the value
		var filteredValues []string
		for i := 0; i < len(Field.OcrLanguageCombo.Options); i++ {
			if len(Field.OcrLanguageCombo.Options) > i && strings.Contains(strings.ToLower(Field.OcrLanguageCombo.Options[i]), strings.ToLower(value)) {
				filteredValues = append(filteredValues, Field.OcrLanguageCombo.Options[i])
			}
		}
		Field.OcrLanguageCombo.SetOptionsFilter(filteredValues)
		Field.OcrLanguageCombo.ShowCompletion()
	}
	Field.OcrWindowCombo.UpdateBeforeOpenFunc = func() {
		sendMessage := SendMessageStruct{
			Type: "get_windows_list",
		}
		sendMessage.SendMessage()
	}

	Field.OscLimitHint.TextSize = theme.TextSize()

	Field.TranscriptionInputHint.TextSize = theme.CaptionTextSize()
	Field.TranscriptionInputHint.Alignment = fyne.TextAlignLeading
	Field.TranscriptionInput.OnChanged = func(value string) {
		Field.TranscriptionInputHint.Text = fmt.Sprintf("%d", len([]rune(value)))
		Field.TranscriptionInputHint.Refresh()

		OscLimitHintUpdateFunc()
	}
	Field.TranscriptionTranslationInputHint.TextSize = theme.CaptionTextSize()
	Field.TranscriptionTranslationInputHint.Alignment = fyne.TextAlignLeading
	Field.TranscriptionTranslationInput.OnChanged = func(value string) {
		Field.TranscriptionTranslationInputHint.Text = fmt.Sprintf("%d", len([]rune(value)))
		Field.TranscriptionTranslationInputHint.Refresh()

		OscLimitHintUpdateFunc()
	}

}
