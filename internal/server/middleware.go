package server

import (
	"compress/gzip"
	"net/http"
)

func DecompressGZIP(next http.Handler) http.Handler {
	// приводим возвращаемую функцию к типу функций HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(`Content-Encoding`) == `gzip` { //	если входящий пакет сжат GZIP
			gz, err := gzip.NewReader(r.Body) //	изготавливаем reader-декомпрессор GZIP
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			r.Body = gz //	подменяем стандартный reader из Request на декомпрессор GZIP
			defer gz.Close()
		}
		next.ServeHTTP(w, r) // передаём управление следующему обработчику
	})
}
