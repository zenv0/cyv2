package http_files

import (
	"fmt"
	"net/http"
)

const compiled_path = "handle/methods/compiled/"
const uncompiled_path = "handle/methods/uncompiled/"
const inst_path = "handle/methods/inst/"

func ServFiles() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path[1:] {
		case "scripts/http-flood":
			http.ServeFile(w, r, compiled_path+"http")
			break
		case "scripts/udp-bypass.pl":
			http.ServeFile(w, r, uncompiled_path+"byp.pl")
			break
		case "scripts/udp-flood.pl":
			http.ServeFile(w, r, uncompiled_path+"udp.pl")
			break
		case "scripts/mix.py":
			http.ServeFile(w, r, uncompiled_path+"mix.py")
			break
		case "scripts/vpn-down.pl":
			http.ServeFile(w, r, uncompiled_path+"vpn.pl")
			break
		case "scripts/home-freeze":
			http.ServeFile(w, r, compiled_path+"tch")
			break
		case "scripts/home-stomp":
			http.ServeFile(w, r, compiled_path+"home-stomp")
			break
		case "scripts/connect":
			http.ServeFile(w, r, inst_path+"connect.py")
			break
		default:
			fmt.Fprintf(w, "wyd over here bro")
			break
		}
	})

	http.ListenAndServe(":9485", nil)

}
