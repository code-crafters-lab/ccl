package license

import "github.com/code-crafters-lab/ccl/pkg/grpc/license"

type fileLicense struct {
	license.LicenseFile
}

func (fl *fileLicense) ID() string {
	return *fl.Payload.LicenseId
}

func NewFileLicense() License {
	return &fileLicense{}
}
