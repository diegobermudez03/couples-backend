package files

import "errors"

var (
	ErrPathNotLongEnough = errors.New("PATH_NOT_LONG_ENOUGH")
	ErrDetectingImageType = errors.New("ERROR_DETECTING_TYPE")
	ErrInvalidImageType = errors.New("INVALID_IMAGE_TYPE")
	ErrUploadingImage 	= errors.New("UNABLE_TO_UPLOAD_IMAGE")
)