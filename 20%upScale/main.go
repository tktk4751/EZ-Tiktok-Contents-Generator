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
	outputFileName := "output_20%upscale_" + time.Now().Format("20060102150405") + ".mp4"
	outputFile := filepath.Join("src/output", outputFileName)
	backgroundColor := "2f2f2f"
	text := "今日のハイライト"
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

	scaleFilter := fmt.Sprintf("scale=iw*%f:ih*%f", zoomFactor, zoomFactor)
	cropFilter := fmt.Sprintf("crop=1080:1920*0.6")
	padFilter := fmt.Sprintf("pad=1080:1920:(ow-iw)/2:(oh-ih)/2:%s", backgroundColor)

	textFilter := fmt.Sprintf("drawtext=text='%s':x=(w-text_w)/2:y=(h-text_h)/8:fontsize=%d:fontfile=%s:fontcolor=%s:bordercolor=%s:borderw=%d:shadowcolor=%s:shadowx=%d:shadowy=%d", text, fontsize, fontfile, textColor, borderColor, borderWidth, shadowColor, shadowXOffset, shadowYOffset)

	filter := fmt.Sprintf("%s,%s,%s,%s", scaleFilter, cropFilter, padFilter, textFilter)

	cmd := exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-vf", filter,
		"-c:a", "copy",
		outputFile,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
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
