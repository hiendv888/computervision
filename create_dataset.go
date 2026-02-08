package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("=== TẠO DATASET GIẢ LẬP ===\n")

	// Tạo cấu trúc thư mục
	dirs := []string{
		"dataset/images",
		"dataset/labels",
		"dataset/visualized",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Lỗi tạo thư mục %s: %v", dir, err)
		}
		fmt.Printf("✓ Đã tạo thư mục: %s\n", dir)
	}

	// Tải 5 ảnh mẫu
	fmt.Println("\n=== TẢI ẢNH MẪU ===")
	imageURLs := []string{
		"https://picsum.photos/800/600?random=1",
		"https://picsum.photos/800/600?random=2",
		"https://picsum.photos/800/600?random=3",
		"https://picsum.photos/800/600?random=4",
		"https://picsum.photos/800/600?random=5",
	}

	for i, url := range imageURLs {
		filename := fmt.Sprintf("img%03d.jpg", i+1)
		filepath := filepath.Join("dataset/images", filename)

		fmt.Printf("Đang tải %s...", filename)
		if err := downloadImage(url, filepath); err != nil {
			log.Fatalf("\nLỗi tải ảnh %s: %v", filename, err)
		}
		fmt.Println(" ✓")
	}

	// Tạo labels mẫu
	fmt.Println("\n=== TẠO LABELS MẪU ===")
	labels := map[string][]string{
		"img001.txt": {
			"0 150 100 200 350", // person
			"1 400 300 100 120", // bag
		},
		"img002.txt": {
			"0 200 150 180 300", // person
			"2 500 200 150 100", // car
		},
		"img003.txt": {
			"0 100 80 150 280",  // person
			"0 350 120 160 290", // person
			"1 550 350 80 90",   // bag
		},
		"img004.txt": {
			"3 250 200 200 180", // bicycle
			"0 100 150 140 250", // person
		},
		"img005.txt": {
			"4 300 250 180 150", // dog
			"0 50 100 130 280",  // person
		},
	}

	for filename, lines := range labels {
		filepath := filepath.Join("dataset/labels", filename)

		file, err := os.Create(filepath)
		if err != nil {
			log.Fatalf("Lỗi tạo file %s: %v", filename, err)
		}

		for _, line := range lines {
			fmt.Fprintln(file, line)
		}
		file.Close()

		fmt.Printf("✓ Đã tạo label: %s (%d boxes)\n", filename, len(lines))
	}

	fmt.Println("\n=== HOÀN THÀNH ===")
	fmt.Println("Dataset đã được tạo tại: dataset/")
	fmt.Println("- 5 ảnh trong dataset/images/")
	fmt.Println("- 5 file labels trong dataset/labels/")
	fmt.Println("\nChạy 'go run visualize_dataset.go' để xem kết quả")
}

// downloadImage tải ảnh từ URL và lưu vào file
func downloadImage(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.ReadFrom(resp.Body)
	return err
}
