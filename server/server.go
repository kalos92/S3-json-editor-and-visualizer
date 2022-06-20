package server

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"strings"

	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"s3-file-editor/appconfig"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/gorilla/mux"
)

type APIResponse struct {
	Err         bool                     `json:"error"`
	ErrorString string                   `json:"error_string"`
	Response    []map[string]interface{} `json:"response,omitempty"`
	GoError     string                   `json:"go_error"`
}

func ReturnError(err error, humanError string) []byte {
	log.Println(err)
	bytesResp, _ := json.Marshal(APIResponse{
		Err:         true,
		ErrorString: humanError,
		GoError:     err.Error(),
	})

	return bytesResp
}

//go:embed build/*
var vueSite embed.FS

type LastS3Request struct {
	s3key  string
	bucket string
	region string
}

type Server struct {
	s3Client   *s3.Client
	lastState  LastS3Request
	ctx        context.Context
	athenaType bool
	conf       appconfig.Config
}

type OverWriteRequest struct {
	Region    string `json:"region"`
	Bucket    string `json:"bucket"`
	JsonBody  string `json:"jsonBody"`
	S3Key     string `json:"S3Key"`
	Athena    bool   `json:"athena"`
	Formatted bool   `json:"formatted"`
}

func NewServer(c appconfig.Config) *Server {
	return &Server{
		ctx:  context.TODO(),
		conf: c,
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/build/", http.StatusMovedPermanently)
}

func createS3Client(ctx context.Context, region string) *s3.Client {

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	if err != nil {
		panic(err)
	}

	return s3.NewFromConfig(cfg)
}

//GET
func (srv *Server) requestJSONtoS3(w http.ResponseWriter, r *http.Request) {
	log.Println("Requested a file")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	var err error
	srv.lastState.bucket = r.URL.Query().Get("bucket")
	srv.lastState.s3key = r.URL.Query().Get("s3key")
	srv.lastState.region = r.URL.Query().Get("region")
	srv.athenaType, err = strconv.ParseBool(r.URL.Query().Get("athena"))

	if err != nil {
		w.Write(ReturnError(err, "Bad Request"))
		return
	}

	srv.s3Client = createS3Client(srv.ctx, srv.lastState.region)

	output, err := srv.s3Client.GetObject(srv.ctx, &s3.GetObjectInput{
		Bucket: aws.String(srv.lastState.bucket),
		Key:    aws.String(srv.lastState.s3key),
	})

	log.Println("File downloaded", srv.lastState.bucket, srv.lastState.s3key)

	if err != nil {
		w.Write(ReturnError(err, "Error while downloading file from S3"))
		return
	}

	bodyBytes, err := io.ReadAll(output.Body)

	if err != nil {
		w.Write(ReturnError(err, "Error while reading the body from the response"))
		return
	}

	var mm []map[string]interface{}

	if srv.athenaType {
		log.Println(string(bodyBytes))

		jsonStrings := strings.Split(string(bodyBytes), "\n")

		for _, s := range jsonStrings {
			m := map[string]interface{}{}

			tempBytes, err := convertJsonToBytes(s, false, false)
			if err != nil {
				w.Write(ReturnError(err, "Probably not a valid athena json file... Download it manually"))
				return
			}
			newline := []byte{'\n'}
			tempBytes = append(tempBytes, newline...)

			err = json.Unmarshal(tempBytes, &m)
			if err != nil {
				w.Write(ReturnError(err, "Probably not a valid athena json file... Download it manually"))
				return
			}
			mm = append(mm, m)
		}

	} else {
		m := map[string]interface{}{}
		err = json.Unmarshal(bodyBytes, &m)

		if err != nil {
			w.Write(ReturnError(err, "Probably not a valid RFC json file... Download it manually"))
			return
		}
		mm = append(mm, m)
	}

	bytesResp, err := json.Marshal(APIResponse{
		Err:         false,
		ErrorString: "",
		Response:    mm,
	})

	if err != nil {
		w.Write(ReturnError(err, "Error while sending the real response"))
		return
	}

	w.Write(bytesResp)
	log.Println("Request Finished")
}

//POST
func (srv *Server) overwriteJSONtoS3(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	log.Println("Requested to overwrite a file")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		w.Write(ReturnError(err, "Error while reading the body from the response"))
		return
	}

	postRequest := &OverWriteRequest{}
	err = json.Unmarshal(bodyBytes, postRequest)
	if err != nil {
		w.Write(ReturnError(err, "Cannot unmarsharl your body!"))
		return
	}

	if !srv.CheckIfIsTheSameFile(postRequest.Region, postRequest.Bucket, postRequest.S3Key, postRequest.Athena) {
		w.Write(ReturnError(errors.New("you are trying to save a file that you not requested"), "You are changing the type of the dowloaded file, this is too risky"))
		return
	}

	fileBytes := []byte{}

	if postRequest.Athena {

		unquotedJson, err := strconv.Unquote(postRequest.JsonBody)
		if err != nil {
			panic(err)
		}

		jsonFiles := strings.Split(unquotedJson, "\n")
		fmt.Println(jsonFiles)

		for i, s := range jsonFiles {
			if s == "" {
				continue
			}
			tempBytes, err := convertJsonToBytes(s, false, false)

			if err != nil {
				w.Write(ReturnError(err, "Invalid athena Json, cannot manage it"))
				return
			}

			if i != 0 {
				newline := []byte{'\n'}
				fileBytes = append(fileBytes, newline...)
			}

			fileBytes = append(fileBytes, tempBytes...)
		}
	} else {
		tempBytes, err := convertJsonToBytes(postRequest.JsonBody, true, postRequest.Formatted)

		if err != nil {
			w.Write(ReturnError(err, "Invalid Json, cannot manage it"))
			return
		}
		fileBytes = append(fileBytes, tempBytes...)
	}

	reader := bytes.NewReader(fileBytes)

	input := &s3.PutObjectInput{
		Bucket:      aws.String(srv.lastState.bucket),
		Key:         aws.String(srv.lastState.s3key),
		Body:        reader,
		ContentType: aws.String("application/json"),
	}

	_, err = srv.s3Client.PutObject(srv.ctx, input)

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			code := apiErr.ErrorCode()
			message := apiErr.ErrorMessage()
			log.Println(code, message)
			w.Write(ReturnError(err, "cannot upload on S3"))
			return
		}
		w.Write(ReturnError(err, "cannot upload on S3"))
		return
	}

	bytesResp, err := json.Marshal(APIResponse{
		Err:         false,
		ErrorString: "",
		Response: []map[string]interface{}{
			{"message": "Upload completed"},
		},
	})

	if err != nil {
		w.Write(ReturnError(err, "Error while sending the real response"))
		return
	}

	w.Write(bytesResp)
	log.Println("Request Finished")

}
func convertJsonToBytes(s string, unquote bool, format bool) ([]byte, error) {

	unquotedJson := ""
	var err error
	if unquote {
		unquotedJson, err = strconv.Unquote(s)
		if err != nil {
			return []byte{}, err
		}
	} else {
		unquotedJson = s
	}

	m := make(map[string]interface{})

	err = json.Unmarshal([]byte(unquotedJson), &m)
	if err != nil {
		return []byte{}, err
	}

	var tempBytes []byte
	if format {
		tempBytes, err = json.MarshalIndent(m, "", "\t")
	} else {
		tempBytes, err = json.Marshal(m)
	}

	if err != nil {
		return []byte{}, err
	}

	return tempBytes, nil
}

func (srv *Server) HandleRequests() {

	router := mux.NewRouter().StrictSlash(true)
	statikFS := http.FS(vueSite)

	staticServer := http.FileServer(statikFS)
	sh := http.StripPrefix("", staticServer)
	router.PathPrefix("/build").Handler(sh)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/api/request-json-to-s3", srv.requestJSONtoS3).Methods("GET")
	router.HandleFunc("/api/overwrite-json-to-s3", srv.overwriteJSONtoS3).Methods("POST", "OPTIONS")

	log.Println("Starting Server")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(srv.conf.WebPort), router))
}

func (srv *Server) CheckIfIsTheSameFile(region, bucket, S3Key string, athena bool) bool {

	if region == srv.lastState.region && bucket == srv.lastState.bucket && S3Key == srv.lastState.s3key && athena == srv.athenaType {
		return true
	}

	return false

}
