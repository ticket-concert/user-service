package response

type RegisterUser struct {
	Email string `json:"email"`
}

type VerifyRegister struct {
	AuthToken    string `json:"authToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiredAt    string `json:"expiredAt"`
}

type LoginUserResp struct {
	AuthToken    string `json:"authToken" bson:"authToken"`
	RefreshToken string `json:"refreshToken" bson:"refreshToken"`
	ExpiredAt    string `json:"expiredAt" bson:"password"`
}

type GetProfile struct {
	UserId        string `json:"userId"`
	FullName      string `json:"fullName"`
	Email         string `json:"email"`
	NIK           string `json:"nik"`
	MobileNumber  string `json:"mobileNumber"`
	Address       string `json:"address"`
	CountryCode   string `json:"countryCode"`
	CountryName   string `json:"countryName"`
	ContinentName string `json:"continentName" bson:"continentName"`
	Latitude      string `json:"latitude" bson:"latitude"`
	Longitude     string `json:"longitude" bson:"longitude"`
	RtRw          string `json:"rtRw"`
	Role          string `json:"role"`
	KKNumber      string `json:"kkNumber"`
}
