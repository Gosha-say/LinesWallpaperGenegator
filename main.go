package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math/rand"
	"net/http"
	url2 "net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("Version 0.1")
	http.HandleFunc("/", httpRoot)
	http.HandleFunc("/image", httpImage)

	err := http.ListenAndServe(":8888", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server clsed")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func httpRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got %s request\n", r.URL.Path)
	w.Header().Set("Content-Type", "text/html")
	html := "<form method='get' action='/image'><input placeholder='Width' name='w'><input placeholder='Height' name='h'><button type='submit'>Generate</button></form>\n"
	_, err := w.Write([]byte(html))
	if err != nil {
		fmt.Printf("WriteString error: %s\n", err)
		os.Exit(1)
	}
}

func httpImage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got %s request\n", r.URL.Path)
	buf := new(bytes.Buffer)
	url, err := url2.Parse(r.URL.String())
	if err != nil {
		fmt.Printf("URL error: %s", err)
	}

	query := url.Query()
	ws, err := strconv.Atoi(strings.Join(query["w"], ""))
	hs, err := strconv.Atoi(strings.Join(query["h"], ""))
	if err != nil {
		fmt.Printf("URL convertion error: %s", err)
	}
	genImage := generateImage(ws, hs)
	err = jpeg.Encode(buf, genImage, nil)
	send := buf.Bytes()
	if err != nil {
		fmt.Printf("Image error: %s", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	write, err := w.Write(send)
	if err != nil {
		fmt.Printf("Image error %d: %s", write, err)
	}
}

func generateImage(width int, height int) image.Image {

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}

	point := getPoint()
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	rand.Seed(time.Now().UnixNano())
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x%10 == 0 {
				point = getPoint()
			}
			img.Set(x, y, point)
		}
	}
	return img
}

func getPoint() color.RGBA {
	return color.RGBA{R: uint8(rand.Intn(255)), G: uint8(rand.Intn(255)), B: uint8(rand.Intn(255)), A: 0xff}
}
