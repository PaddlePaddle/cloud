package pfsserver

import (
	//"encoding/json"
	//"github.com/cloud/go/file_manager/pfscommon"
	"fmt"
	"github.com/cloud/go/file_manager/pfsmodules"
	"log"
	"net/http"
	"strconv"
)

func lsCmdHandler(w http.ResponseWriter, req *pfsmodules.CmdAttr) {
	resp := pfsmodules.LsCmdResponse{}

	/*
		if req.Method != "ls" {
			resp.SetErr("not surported method:" + req.Method)
			pfsmodules.WriteCmdJsonResponse(w, &resp, http.StatusMethodNotAllowed)
			return
		}
	*/

	log.Print(req)

	cmd := pfsmodules.NewLsCmd(req, &resp)
	cmd.RunAndResponse(w)

	return
}

func MD5SumCmdHandler(w http.ResponseWriter, req *pfsmodules.CmdAttr) {
	resp := pfsmodules.MD5SumResponse{}
	log.Print(req)

	cmd := pfsmodules.NewMD5SumCmd(req, &resp)
	cmd.RunAndResponse(w)
}

func GetFilesHandler(w http.ResponseWriter, r *http.Request) {
	resp := pfsmodules.LsCmdResponse{}
	req, err := pfsmodules.GetJsonRequestCmdAttr(r)
	if err != nil {
		resp.SetErr(err.Error())
		pfsmodules.WriteCmdJsonResponse(w, &resp, 422)
		return
	}

	if len(req.Args) == 0 {
		resp.SetErr("no args")
		pfsmodules.WriteCmdJsonResponse(w, &resp, http.StatusExpectationFailed)
		return

	}

	switch req.Method {
	case "ls":
		lsCmdHandler(w, req)
	case "md5sum":
		MD5SumCmdHandler(w, req)
	default:
		resp.SetErr(http.StatusText(http.StatusMethodNotAllowed))
		pfsmodules.WriteCmdJsonResponse(w, &resp, http.StatusMethodNotAllowed)
	}

	log.Print(req)
}

/*
func SendHttpTxtResponse(w http.ResponseWriter, status int32) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "%s %s",
		strconv.Itoa(http.StatusMethodNotAllowed),
		http.StatusText(http.StatusMethodNotAllowed))
}

func SendJsonResponse(w http.ResponseWriter, status int32) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "%s %s",
		strconv.Itoa(http.StatusMethodNotAllowed),
		http.StatusText(http.StatusMethodNotAllowed))
}
*/
func rmCmdHandler(w http.ResponseWriter, req *pfsmodules.CmdAttr) {
	resp := pfsmodules.RmCmdResponse{}

	log.Print(req)

	cmd := pfsmodules.NewRmCmd(req, &resp)
	cmd.RunAndResponse(w)

	return
}

func touchHandler(w http.ResponseWriter, req *pfsmodules.CmdAttr) {
	resp := pfsmodules.TouchCmdResponse{}

	log.Print(req)

	cmd := pfsmodules.NewTouchCmd(req, &resp)
	cmd.RunAndResponse(w)

	return
}

func PostFilesHandler(w http.ResponseWriter, r *http.Request) {
	resp := pfsmodules.JsonResponse{}
	req, err := pfsmodules.GetJsonRequestCmdAttr(r)
	if err != nil {
		resp.SetErr(err.Error())
		pfsmodules.WriteCmdJsonResponse(w, &resp, http.StatusExpectationFailed)
		return
	}

	if len(req.Args) == 0 {
		resp.SetErr("no args")
		pfsmodules.WriteCmdJsonResponse(w, &resp, http.StatusExpectationFailed)
		return

	}

	log.Print(req)

	switch req.Method {
	case "rm":
		rmCmdHandler(w, req)
	case "touch":
		if len(req.Args) != 1 {
			resp.SetErr("please create only one file")
			pfsmodules.WriteCmdJsonResponse(w, &resp, http.StatusExpectationFailed)
			return
		}
		touchHandler(w, req)
	default:
		resp.SetErr(http.StatusText(http.StatusMethodNotAllowed))
		pfsmodules.WriteCmdJsonResponse(w, &resp, http.StatusMethodNotAllowed)
	}
}

func GetChunkMetaHandler(w http.ResponseWriter, r *http.Request) {
	method := r.URL.Query().Get("method")

	log.Println(r.URL.String())

	switch method {
	case "getchunkmeta":
		cmd := pfsmodules.GetChunkMetaCmd(w, r)
		if cmd == nil {
			return
		}
		cmd.RunAndResponse(w)
	default:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		//w.Write(http.StatusText(http.StatusMethodNotAllowed))
		fmt.Fprintf(w, "%s %s",
			strconv.Itoa(http.StatusMethodNotAllowed),
			http.StatusText(http.StatusMethodNotAllowed))
	}
}

func GetChunksHandler(w http.ResponseWriter, r *http.Request) {
}

func PatchChunksHandler(w http.ResponseWriter, r *http.Request) {
}

/*
func GetChunksHandler(w http.ResponseWriter, r *http.Request) {
	method := r.URL.Query().Get("method")

	log.Println(r.URL.String())

	switch method {
	case "rm":
		cmd := pfsmodules.GetRMCmd(w, r)
		if cmd == nil {
			return
		}
		cmd.RunAndResponse(w)
	default:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		//w.Write(http.StatusText(http.StatusMethodNotAllowed))
		fmt.Fprintf(w, "%s %s",
			strconv.Itoa(http.StatusMethodNotAllowed),
			http.StatusText(http.StatusMethodNotAllowed))
	}
}
*/

func PostChunksHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("w")
	//for k
	path := r.URL.Query().Get("path")

	log.Println(path)
}
