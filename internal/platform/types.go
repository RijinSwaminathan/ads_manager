package platform

// Credentials stores API authentication details
type Credentials struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
}
