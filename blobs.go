package clccam

// BlobResponse is returned when uploading a file via POST.
type BlobResponse struct {
	// Url from which the uploaded file can be retrieved, e.g. "/services/blobs/download/5c1abf95939a600ea38a8661/test.sh"
	Url string `json:"url"`

	// Length of the file in bytes
	Length int `json:"length"`

	// MIME type of the file
	ContentType string `json:"content_type"`

	// Date/time of the upload
	When Timestamp `json:"upload_date"`
}

// UploadFile uploads the contents of file @name contained in @b.
func (c *Client) UploadFile(name string, b []byte) (res BlobResponse, err error) {
	return res, c.getResponse("/services/blobs/upload/"+name, "POST", b, &res)
}
