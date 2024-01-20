package service

import (
	"fmt"
	"os"
	"speech/config"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

type SpeechRequest struct {
	Content      string
	VoiceName    string
	FileName     string
	ResponseChan chan SpeechResponse
}

type SpeechResponse struct {
	Error    string
	FileName string
}

var WaitChan chan SpeechRequest

func init() {
	WaitChan = make(chan SpeechRequest, 20)
}

func synthesizeStartedHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	fmt.Println("Synthesis started.")
}

func synthesizingHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	fmt.Printf("Synthesizing, audio chunk size %d.\n", len(event.Result.AudioData))
}

func synthesizedHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	fmt.Printf("Synthesized, audio length %d.\n", len(event.Result.AudioData))
}

func cancelledHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	fmt.Println("Received a cancellation.")
}

func Speek() {
	// This example requires environment variables named "SPEECH_KEY" and "SPEECH_REGION"
	audioConfig, err := audio.NewAudioConfigFromDefaultSpeakerOutput()
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}

	defer audioConfig.Close()
	speechConfig, err := speech.NewSpeechConfigFromSubscription(config.AppConfig.Azure.Key, config.AppConfig.Azure.Region)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}

	defer speechConfig.Close()

	speechConfig.SetSpeechSynthesisVoiceName(config.AppConfig.Azure.DefaultVoice)

	speechSynthesizer, err := speech.NewSpeechSynthesizerFromConfig(speechConfig, audioConfig)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer speechSynthesizer.Close()

	speechSynthesizer.SynthesisStarted(synthesizeStartedHandler)
	speechSynthesizer.Synthesizing(synthesizingHandler)
	speechSynthesizer.SynthesisCompleted(synthesizedHandler)
	speechSynthesizer.SynthesisCanceled(cancelledHandler)

	for req := range WaitChan {
		go func(req SpeechRequest) {
			if req.VoiceName != "" {
				speechConfig.SetSpeechSynthesisVoiceName(req.VoiceName)
				speechSynthesizer, err = speech.NewSpeechSynthesizerFromConfig(speechConfig, audioConfig)
				if err != nil {

					req.ResponseChan <- SpeechResponse{
						Error: fmt.Sprintf("Got an error: %s", err),
					}
					return
				}

			}

			task := speechSynthesizer.SpeakTextAsync(req.Content)
			var outcome speech.SpeechSynthesisOutcome
			select {
			case outcome = <-task:
			case <-time.After(60 * time.Second):
				req.ResponseChan <- SpeechResponse{
					Error: "Timed out",
				}
				return
			}
			defer outcome.Close()
			if outcome.Error != nil {
				req.ResponseChan <- SpeechResponse{
					Error: fmt.Sprintf("Got an error: %s", outcome.Error),
				}
				return
			}

			if outcome.Result.Reason == common.SynthesizingAudioCompleted {
				// fmt.Printf("Speech synthesized to speaker for text [%s].\n", req.Content)
				fileName := fmt.Sprintf("./%d.wav", time.Now().UnixMicro()) // Change the filename and extension as needed
				if req.FileName != "" {
					fileName = fmt.Sprintf("./%s.wav", req.FileName)
				}

				if err := os.WriteFile(fileName, outcome.Result.AudioData, 0644); err != nil {
					req.ResponseChan <- SpeechResponse{
						Error: fmt.Sprintf("Error saving audio to file: %s", err),
					}
					return
				}

				req.ResponseChan <- SpeechResponse{
					FileName: fileName,
				}

			} else {
				err := ""
				cancellation, _ := speech.NewCancellationDetailsFromSpeechSynthesisResult(outcome.Result)
				err = fmt.Sprintf("CANCELED: Reason=%d.\n", cancellation.Reason)
				if cancellation.Reason == common.Error {
					err += fmt.Sprintf("CANCELED: ErrorCode=%d\nCANCELED: ErrorDetails=[%s]\nCANCELED: Did you set the speech resource key and region values?\n",
						cancellation.ErrorCode,
						cancellation.ErrorDetails)
				}
				req.ResponseChan <- SpeechResponse{
					Error: err,
				}
				return
			}
		}(req)
	}
}
