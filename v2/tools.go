package toolkit

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const randomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+"

// Tools is the type used to instantiate this module
// Any variable of this type will have access to all methods with the receiver *Tools
// Also see https://developer.mozilla.org/en-US/docs/Web/HTTP/MIME_types/Common_types
type Tools struct {
	MaxFileSize        int64
	AllowedFileTypes   []string
	MaxJSONSize        int
	AllowUnknownFields bool
}

// RandomString returns a string of random characters of length n
func (t *Tools) RandomString(n int) string {
	// old expensive version for reference
	//s, r := make([]rune, n), []rune(randomStringSource)
	//for i := range s {
	//	p, _ := rand.Prime(rand.Reader, len(r))
	//	x, y := p.Uint64(), uint64(len(r))
	//	s[i] = r[x%y]
	//}
	//
	//return string(s)

	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	for i := range b {
		b[i] = randomStringSource[int(b[i])%len(randomStringSource)]
	}

	return string(b)

}

// UploadedFile is a struct used to save information about an uploaded file
type UploadedFile struct {
	Key              string
	NewFileName      string
	OriginalFileName string
	FileSize         int64
}

// UploadOneFile is a convenience methods that calls UploadFiles, but expects only one file
func (t *Tools) UploadOneFile(r *http.Request, uploadDir string, rename ...bool) (*UploadedFile, error) {
	renameFile := true
	if len(rename) > 0 {
		renameFile = rename[0]
	}

	files, err := t.UploadFiles(r, uploadDir, renameFile)
	if err != nil {
		return nil, err
	}

	return files[0], nil
}

// UploadFiles uploads one of more file to a specified directory, and give the files
// a random name if required. If the optional last parameter is set to false, it will
// not rename and use the original name.
func (t *Tools) UploadFiles(r *http.Request, uploadDir string, rename ...bool) ([]*UploadedFile, error) {
	renameFile := true
	if len(rename) > 0 {
		renameFile = rename[0]
	}

	var uploadedFiles []*UploadedFile

	if t.MaxFileSize == 0 {
		t.MaxFileSize = 1024 * 1024 * 1024
	}

	err := t.CreateDirIfNotExist(uploadDir)
	if err != nil {
		return nil, err
	}

	err = r.ParseMultipartForm(t.MaxFileSize)
	if err != nil {
		return nil, errors.New("uploaded file too big")
	}

	for key, fHeaders := range r.MultipartForm.File {
		for _, hdr := range fHeaders {
			// we will be deferring stuff below, so
			// since we are in a loop in-line a func
			uploadedFiles, err = func(uploadedFiles []*UploadedFile) ([]*UploadedFile, error) {
				var uploadedFile UploadedFile
				infile, err := hdr.Open()
				if err != nil {
					return nil, err
				}
				defer infile.Close()

				buff := make([]byte, 512)
				_, err = infile.Read(buff)
				if err != nil {
					return nil, err
				}

				// check to see if file type is permitted
				allowed := false
				fileType := http.DetectContentType(buff)

				if len(t.AllowedFileTypes) > 0 {
					for _, x := range t.AllowedFileTypes {
						// case insensitive compare
						if strings.EqualFold(fileType, x) {
							allowed = true
						}
					}
				} else {
					allowed = true
				}

				if !allowed {
					return nil, errors.New("the uploaded file type is not permitted")
				}

				_, err = infile.Seek(0, 0)
				if err != nil {
					return nil, err
				}

				if renameFile {
					uploadedFile.NewFileName = fmt.Sprintf("%s%s", t.RandomString(25), filepath.Ext(hdr.Filename))
				} else {
					uploadedFile.NewFileName = hdr.Filename
				}
				uploadedFile.OriginalFileName = hdr.Filename
				uploadedFile.Key = key

				var outfile *os.File
				defer outfile.Close()

				if outfile, err = os.Create(filepath.Join(uploadDir, uploadedFile.NewFileName)); err != nil {
					return nil, err
				} else {
					fileSize, err := io.Copy(outfile, infile)
					if err != nil {
						return nil, err
					}
					uploadedFile.FileSize = fileSize
				}

				uploadedFiles = append(uploadedFiles, &uploadedFile)
				return uploadedFiles, nil
			}(uploadedFiles)

			if err != nil {
				// return the files that were uploaded, if any
				return uploadedFiles, err
			}
		}
	}
	return uploadedFiles, nil
}

// CreateDirIfNotExist creates a directory, and all necessary parents
func (t *Tools) CreateDirIfNotExist(path string) error {
	const mode = 0755
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, mode)
		if err != nil {
			return err
		}
	}

	return nil
}

// Slugify a very simple means of creating a slug from a string
func (t *Tools) Slugify(s string) (string, error) {
	if s == "" {
		return "", errors.New("empty string not permitted")
	}

	//any letters or digits of any length
	var re = regexp.MustCompile(`[^a-z\d]+`)
	slug := strings.Trim(re.ReplaceAllString(strings.ToLower(s), "-"), "-")
	if s == "" {
		return "", errors.New("after removing characters, slug is zero length")
	}

	return slug, nil
}

// DownloadStaticFile downloads a file tries to force the browser to
// avoid displaying it in the browser window by setting content disposition.
// It also allows specification of the display name
func (t *Tools) DownloadStaticFile(w http.ResponseWriter, r *http.Request, pathName, displayName string) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", displayName))

	http.ServeFile(w, r, pathName)
}

// JSONResponse is the type used for sending JSON around
type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ReadJSON tries to read the body of a request and converts from json into a go data variable
func (t *Tools) ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1024 * 1024 // 1MB
	if t.MaxJSONSize != 0 {
		maxBytes = t.MaxJSONSize
	}

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	if !t.AllowUnknownFields {
		dec.DisallowUnknownFields()
	}

	err := dec.Decode(data)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.Trim(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger that %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			return fmt.Errorf("error unmarshalling JSON: %s", err.Error())

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contain only one JSON structure")
	}

	return nil
}

// WriteJSON takes a response status code and arbitrary data and writes json to the client
func (t *Tools) WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

// ErrorJSON takes an error and optionally a status code, and sends a JSON error message
func (t *Tools) ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	return t.WriteJSON(w, statusCode, payload)
}

// HandleError wraps ErrorJSON - Returns true if error was handled and caller should abort.
// Outputs the error JSON or writes to the logger on failure
func (t *Tools) HandleError(w http.ResponseWriter, err error, logger *log.Logger) bool {
	if err == nil {
		return false
	}

	if jsonErr := t.ErrorJSON(w, err); jsonErr != nil {
		logger.Printf("Error writing JSON response: %v", jsonErr)
	}
	return true
}

// PushJSONToRemote posts arbitrary data to some URL as JSON, and returns the response, status code, and error, if any
// The final parameter, client, is optional. If none is specified, we use the standard http.Client
func (t *Tools) PushJSONToRemote(uri string, data interface{}, client ...*http.Client) (*http.Response, int, error) {
	// create json
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, 0, err
	}

	// check for custom http client
	httpClient := &http.Client{}
	if len(client) > 0 {
		httpClient = client[0]
	}

	// build request and set header
	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, 0, err
	}
	request.Header.Set("Content-Type", "application/json")

	// call remote url
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer response.Body.Close()

	// send response back
	return response, response.StatusCode, nil
}

// WriteXML takes a response status code and arbitrary data and writes xml to the client
func (t *Tools) WriteXML(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	// in production, we would not use indent
	out, err := xml.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}
