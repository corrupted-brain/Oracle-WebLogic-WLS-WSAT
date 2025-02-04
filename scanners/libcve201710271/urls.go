package libcve201710271

var (
	// DefaultURLs is the endpoint URL that we should scan by default.
	DefaultURLs = []string{
		"/wls-wsat/CoordinatorPortType",
	}

	// AllURLs is the endpoint URLs that are known to be vulnerable and may be
	// desirable to scan in certain cases.
	AllURLs = []string{
		"/wls-wsat/CoordinatorPortType",
		"/wls-wsat/CoordinatorPortType11",
		"/wls-wsat/ParticipantPortType",
		"/wls-wsat/ParticipantPortType11",
		"/wls-wsat/RegistrationPortTypeRPC",
		"/wls-wsat/RegistrationPortTypeRPC11",
		"/wls-wsat/RegistrationRequesterPortType",
		"/wls-wsat/RegistrationRequesterPortType11",
	}
)
