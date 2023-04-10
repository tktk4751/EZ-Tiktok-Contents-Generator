package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	inputFile := "src/input/input.mp4"
	outputFileName := "output_fullscale_" + time.Now().Format("20060102150405") + ".mp4"
	outputFile := filepath.Join("src/output", outputFileName)
	backgroundColor := "2f2f2f" //Feel free to change the color code
	text := "今日のハイライト"          //Feel free to change text
	zoomFactor := 1.2
	err := resizeAndAddText(inputFile, outputFile, text, backgroundColor, zoomFactor)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func resizeAndAddText(inputFile, outputFile, text, backgroundColor string, zoomFactor float64) error {
	fontfile := "src/font/JNRfont_n.ttf"
	fontsize := 120
	textColor := "white"
	borderColor := "Red"
	borderWidth := 2
	shadowColor := "black@0.5"
	shadowXOffset := 3
	shadowYOffset := 3

	ffprobeCmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_entries", "stream=width,height",
		"-of", "default=nw=1:nk=1",
		inputFile,
	)
	widthHeightBytes, err := ffprobeCmd.Output()
	if err != nil {
		return fmt.Errorf("ffprobe command execution failed: %w", err)
	}
	widthHeight := string(widthHeightBytes)
	var width, height int
	fmt.Sscanf(widthHeight, "%d\n%d", &width, &height)

	cropFilter := ""
	if float64(height)*9/16 <= float64(width) {
		cropFilter = fmt.Sprintf("crop=ih*(9/16):ih")
	} else {
		cropFilter = fmt.Sprintf("crop=iw:iw*(16/9)")
	}

	scaleFilter := fmt.Sprintf("scale=%d:%d", 1080, 1920)

	textFilter := fmt.Sprintf("drawtext=text='%s':x=(w-text_w)/2:y=(h-text_h)/8:fontsize=%d:fontfile=%s:fontcolor=%s:bordercolor=%s:borderw=%d:shadowcolor=%s:shadowx=%d:shadowy=%d", text, fontsize, fontfile, textColor, borderColor, borderWidth, shadowColor, shadowXOffset, shadowYOffset)

	padFilter := fmt.Sprintf("scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:%s,%s", backgroundColor, textFilter)

	// padFilter := fmt.Sprintf("scale=%d:1080:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2", scaledWidth)

	filter := fmt.Sprintf("%s,%s,%s,%s", cropFilter, scaleFilter, padFilter, textFilter)

	cmd := exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-vf", filter,
		"-c:a", "copy",
		outputFile,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("ffmpeg command execution failed: %w", err)
	}

	json, err := json.Marshal(map[string]interface{}{
		"input_file":    inputFile,
		"output_file":   outputFile,
		"background":    backgroundColor,
		"text":          text,
		"zoom_factor":   zoomFactor,
		"font_file":     fontfile,
		"font_size":     fontsize,
		"text_color":    textColor,
		"output_width":  1080,
		"output_height": 1920,
		"output_aspect": "9:16",
	})
	if err != nil {
		return fmt.Errorf("failed to encode output JSON: %w", err)
	}

	fmt.Println(string(json))

	return nil
}
