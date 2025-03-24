package apiversion

type Info struct {
	SpecRootUrl     string `json:"specRootUrl"`
	Major           *int   `json:"major"`
	Minor           *int   `json:"minor"`
	SupportedMajors []int  `json:"supportedMajors"`
}

type ApiVersionResponse struct {
	Specs []Info `json:"specs"`
}
