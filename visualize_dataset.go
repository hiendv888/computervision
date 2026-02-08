package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// BoundingBox đại diện cho một bounding box
type BoundingBox struct {
	ClassID int
	X       int
	Y       int
	Width   int
	Height  int
}

// Class colors
var classColors = map[int]color.RGBA{
	0: {R: 255, G: 0, B: 0, A: 255},   // person - red
	1: {R: 0, G: 255, B: 0, A: 255},   // bag - green
	2: {R: 0, G: 0, B: 255, A: 255},   // car - blue
	3: {R: 255, G: 255, B: 0, A: 255}, // bicycle - yellow
	4: {R: 255, G: 0, B: 255, A: 255}, // dog - magenta
}

// Class names
var classNames = map[int]string{
	0: "person",
	1: "bag",
	2: "car",
	3: "bicycle",
	4: "dog",
}

func main() {
	fmt.Println("=== VISUALIZE DATASET ===\n")

	// Đọc danh sách ảnh
	imageFiles, err := filepath.Glob("dataset/images/*.jpg")
	if err != nil {
		log.Fatalf("Lỗi đọc thư mục images: %v", err)
	}

	if len(imageFiles) == 0 {
		log.Fatal("Không tìm thấy ảnh nào trong dataset/images/")
	}

	fmt.Printf("Tìm thấy %d ảnh\n\n", len(imageFiles))

	// Xử lý từng ảnh
	for _, imagePath := range imageFiles {
		// Lấy tên file không có extension
		baseName := filepath.Base(imagePath)
		nameWithoutExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))

		fmt.Printf("Đang xử lý: %s\n", baseName)

		// Đọc ảnh
		img, err := loadImage(imagePath)
		if err != nil {
			log.Printf("  ✗ Lỗi đọc ảnh: %v\n", err)
			continue
		}

		// Chuyển sang RGBA
		bounds := img.Bounds()
		rgba := image.NewRGBA(bounds)
		draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

		// Đọc labels
		labelPath := filepath.Join("dataset/labels", nameWithoutExt+".txt")
		boxes, err := readLabels(labelPath)
		if err != nil {
			log.Printf("  ✗ Lỗi đọc labels: %v\n", err)
			continue
		}

		fmt.Printf("  Tìm thấy %d bounding boxes\n", len(boxes))

		// Vẽ các bounding boxes
		for i, box := range boxes {
			color := classColors[box.ClassID]
			drawBoundingBox(rgba, box.X, box.Y, box.Width, box.Height, color, 3)

			// Vẽ label text
			className := classNames[box.ClassID]
			fmt.Printf("    Box %d: %s (%d,%d,%d,%d)\n", i+1, className, box.X, box.Y, box.Width, box.Height)
		}

		// Lưu ảnh kết quả
		outputPath := filepath.Join("dataset/visualized", baseName)
		if err := saveImage(rgba, outputPath); err != nil {
			log.Printf("  ✗ Lỗi lưu ảnh: %v\n", err)
			continue
		}

		fmt.Printf("  ✓ Đã lưu: %s\n\n", outputPath)
	}

	fmt.Println("=== HOÀN THÀNH ===")
	fmt.Println("Kiểm tra kết quả tại: dataset/visualized/")
}

// readLabels đọc file labels và trả về danh sách bounding boxes
func readLabels(filepath string) ([]BoundingBox, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var boxes []BoundingBox
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 5 {
			continue
		}

		classID, _ := strconv.Atoi(parts[0])
		x, _ := strconv.Atoi(parts[1])
		y, _ := strconv.Atoi(parts[2])
		w, _ := strconv.Atoi(parts[3])
		h, _ := strconv.Atoi(parts[4])

		boxes = append(boxes, BoundingBox{
			ClassID: classID,
			X:       x,
			Y:       y,
			Width:   w,
			Height:  h,
		})
	}

	return boxes, scanner.Err()
}

// drawBoundingBox vẽ hình chữ nhật lên ảnh
func drawBoundingBox(img *image.RGBA, x, y, w, h int, c color.Color, thickness int) {
	// Vẽ 4 cạnh của hình chữ nhật

	// Cạnh trên
	for t := 0; t < thickness; t++ {
		for i := x; i < x+w; i++ {
			if i >= 0 && i < img.Bounds().Dx() && y+t >= 0 && y+t < img.Bounds().Dy() {
				img.Set(i, y+t, c)
			}
		}
	}

	// Cạnh dưới
	for t := 0; t < thickness; t++ {
		for i := x; i < x+w; i++ {
			if i >= 0 && i < img.Bounds().Dx() && y+h-t-1 >= 0 && y+h-t-1 < img.Bounds().Dy() {
				img.Set(i, y+h-t-1, c)
			}
		}
	}

	// Cạnh trái
	for t := 0; t < thickness; t++ {
		for i := y; i < y+h; i++ {
			if x+t >= 0 && x+t < img.Bounds().Dx() && i >= 0 && i < img.Bounds().Dy() {
				img.Set(x+t, i, c)
			}
		}
	}

	// Cạnh phải
	for t := 0; t < thickness; t++ {
		for i := y; i < y+h; i++ {
			if x+w-t-1 >= 0 && x+w-t-1 < img.Bounds().Dx() && i >= 0 && i < img.Bounds().Dy() {
				img.Set(x+w-t-1, i, c)
			}
		}
	}
}

// loadImage đọc ảnh từ file
func loadImage(filepath string) (image.Image, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

// saveImage lưu ảnh vào file
func saveImage(img *image.RGBA, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
}
