package main

import (
	"fmt"
	"github.com/CAFxX/httpcompression"
	"github.com/klauspost/compress/gzip"
	"io"
	"net/http"
	"os"
)

func main() {

	data, err := os.ReadFile("./data.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	//got := snappy.Encode(nil, data)
	//fmt.Println("原始长度为:", len(data))
	//fmt.Println("压缩后长度为:", len(got))
	//
	//gotS2 := s2.Encode(nil, data)
	//fmt.Println("原始长度为:", len(data))
	//fmt.Println("压缩后长度为:", len(gotS2))
	//
	//var encoder, _ = zstd.NewWriter(nil)
	//gotZstd := encoder.EncodeAll(data, make([]byte, 0, len(data)))
	//fmt.Println("原始长度为:", len(data))
	//fmt.Println("压缩后长度为:", len(gotZstd))
	//
	////打包：tar cvf /root/Desktop/centos.tar.gz /root/Desktop/centos
	////解包：tar xvf /root/Desktop/centos.tar.gz

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		str := string(data)
		io.WriteString(w, str)
	})
	//zstdComp, err := zstd.New()
	//compression := httpcompression.ZstandardCompressor(zstdComp)
	//compress, _ := httpcompression.DefaultAdapter(compression) // Use the default configuration

	compress, _ := httpcompression.DefaultAdapter(httpcompression.GzipCompressionLevel(gzip.DefaultCompression)) // Use the default configuration

	http.Handle("/", compress(handler))
	http.ListenAndServe("0.0.0.0:8080", nil)
}
