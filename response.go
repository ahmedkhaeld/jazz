package jazz

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gocarina/gocsv"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// WriteJSON Writes json to response writer
func (j *Jazz) WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	//assign headers to response
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

// ReadJSON  allows reading generic json,
// provide clean way to read any kind of json from a request, assuming
// that request body has only a single json value
func (j *Jazz) ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}, AllowUnknownFields bool, MaxJSONSize ...int) error {
	// maxBytes define the max size of request allowed prevents clients from accidentally or maliciously sending a large
	//request and wasting server resources.
	maxBytes := 1048576 // 1MB default
	// if a custom max size is specified, use that instead of default
	if len(MaxJSONSize) > 0 {
		maxBytes = MaxJSONSize[0]
	}
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	//should we allow unknown fields in json
	if !AllowUnknownFields {
		dec.DisallowUnknownFields()
	}

	//attempt to decode the data, and figure out what the error is, to send back a human-readable response
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
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			return fmt.Errorf("error unmarshalling json: %s", err.Error())

		default:
			return err
		}
	}

	// assume to decode a json file that has one entry
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single JSON value")
	}

	return nil
}

func (j *Jazz) ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	payload.Error = true
	payload.Message = err.Error()

	return j.WriteJSON(w, statusCode, payload)
}

func (j *Jazz) ReadCSV(b []byte, data interface{}) error {
	err := gocsv.UnmarshalBytes(b, data)
	if err != nil {
		return err
	}

	return nil
}

func (j *Jazz) WriteXML(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := xml.MarshalIndent(data, "", "   ")
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

func (j *Jazz) ReadXML(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048756
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := xml.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single XML value")
	}

	return nil
}

func (j *Jazz) DownloadFile(w http.ResponseWriter, r *http.Request, pathToFile, fileName string) error {
	fp := path.Join(pathToFile, fileName)
	fileToServe := filepath.Clean(fp)
	w.Header().Set("Content-Type", fmt.Sprintf("attachment; file=\"%s\"", fileName))
	http.ServeFile(w, r, fileToServe)
	return nil
}

// UploadedFile is a struct used for the uploaded file
type UploadedFile struct {
	NewFileName      string
	OriginalFileName string
	FileSize         int64
}

// UploadOneFile is just a convenience method that calls UploadFiles, but expects only one file to
// be in the upload.
func (j *Jazz) UploadOneFile(r *http.Request, uploadDir string, rename bool, AllowedFileTypes []string, MaxFileSize ...int) (*UploadedFile, error) {

	files, err := j.UploadFiles(r, uploadDir, rename, AllowedFileTypes, MaxFileSize...)
	if err != nil {
		return nil, err
	}

	return files[0], nil
}

// UploadFiles uploads one or more file to a specified directory, and gives the files a random name.
// It returns a slice containing the newly named files, the original file names, the size of the files,
// and potentially an error. If the optional last parameter is set to true, then we will not rename
// the files, but will use the original file names.
func (j *Jazz) UploadFiles(r *http.Request, uploadDir string, rename bool, AllowedFileTypes []string, MaxFileSize ...int) ([]*UploadedFile, error) {
	// check to see if we are renaming the uploadedFiles
	renameFile := true
	if rename == false {
		renameFile = rename
	}

	var uploadedFiles []*UploadedFile

	// create the upload directory if it does not exist
	err := j.CreateDirIfNotExist(uploadDir)
	if err != nil {
		return nil, err
	}

	// sanity check on t.MaxFileSize
	if MaxFileSize[0] == 0 {
		MaxFileSize[0] = 1024 * 1024 * 5 // 5 megabytes
	}

	// parse the form so we have access to the file
	err = r.ParseMultipartForm(int64(MaxFileSize[0]))
	if err != nil {
		return nil, fmt.Errorf("the uploaded file is too big, and must be less than %d", MaxFileSize[0])
	}

	for _, fHeaders := range r.MultipartForm.File {
		for _, hdr := range fHeaders {
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

				allowed := false
				filetype := http.DetectContentType(buff)
				if len(AllowedFileTypes) > 0 {
					for _, x := range AllowedFileTypes {
						if strings.EqualFold(filetype, x) {
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
					fmt.Println(err)
					return nil, err
				}

				if renameFile {
					uploadedFile.NewFileName = fmt.Sprintf("%s%s", j.RandomString(25), filepath.Ext(hdr.Filename))
				} else {
					uploadedFile.NewFileName = hdr.Filename
				}
				uploadedFile.OriginalFileName = hdr.Filename

				var outfile *os.File
				defer outfile.Close()

				if outfile, err = os.Create(filepath.Join(uploadDir, uploadedFile.NewFileName)); nil != err {
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
				return uploadedFiles, err
			}
		}
	}
	return uploadedFiles, nil
}

// Slugify is a (very) simple means of creating a slug from a provided string.
//
//example: "This is a test" -> "this-is-a-test"
func (j *Jazz) Slugify(s string) (string, error) {
	if s == "" {
		return "", errors.New("empty string not permitted")
	}
	var re = regexp.MustCompile(`[^a-z\d]+`)
	slug := strings.Trim(re.ReplaceAllString(strings.ToLower(s), "-"), "-")
	if len(slug) == 0 {
		return "", errors.New("after removing characters, slug is zero length")
	}

	return slug, nil
}

///////////////////////// ERROR Response  //////////////////////////

// ErrorStatus returns a response with the supplied http status
func (j *Jazz) ErrorStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Error404 returns page not found response
func (j *Jazz) Error404(w http.ResponseWriter, r *http.Request) {
	j.ErrorStatus(w, http.StatusNotFound)
}

// Error500 returns internal server error response
func (j *Jazz) Error500(w http.ResponseWriter, r *http.Request) {
	j.ErrorStatus(w, http.StatusInternalServerError)
}

// ErrorUnauthorized sends an unauthorized status (client is not known)
func (j *Jazz) ErrorUnauthorized(w http.ResponseWriter, r *http.Request) {
	j.ErrorStatus(w, http.StatusUnauthorized)
}

// ErrorForbidden returns a forbidden status message (client is known)
func (j *Jazz) ErrorForbidden(w http.ResponseWriter, r *http.Request) {
	j.ErrorStatus(w, http.StatusForbidden)
}
