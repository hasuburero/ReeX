package api

import (
	"fmt"
	"io"
	"net/http"
)

// internal package
import (
	"github.com/hasuburero/ReeX/lib/common"
)

func Kill(w http.ResponseWriter, r *http.Request) {
	var ctx common.Post_Kill_Struct
	req_body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
}
