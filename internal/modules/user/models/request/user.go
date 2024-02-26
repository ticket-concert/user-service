package request

const (
	RoleUser        = `user`
	RoleAdmin       = `admin`
	RoleStackHolder = `stackholder`
)

var MapOfRole = map[string]string{
	RoleUser:        RoleUser,
	RoleAdmin:       RoleAdmin,
	RoleStackHolder: RoleStackHolder,
}

type RegisterUser struct {
	FullName      string `json:"fullName" validate:"required"`
	Email         string `json:"email" validate:"required,min=1,max=50"`
	Password      string `json:"password" validate:"required,min=8,max=20"`
	NIK           string `json:"nik" validate:"required"`
	MobileNumber  string `json:"mobileNumber" validate:"required"`
	Address       string `json:"address"`
	ProvinceId    string `json:"provinceId" validate:"required_if=CountryId 100"`
	CityId        string `json:"cityId" validate:"required_if=CountryId 100"`
	DistrictId    string `json:"districtId" validate:"required_if=CountryId 100"`
	SubdictrictId string `json:"subdictrictId" validate:"required_if=CountryId 100"`
	CountryId     string `json:"countryId" validate:"required"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	RtRw          string `json:"rtRw"`
	Role          string `json:"role" validate:"required"`
	KKNumber      string `json:"kkNumber"`
}

type UpdateUser struct {
	FullName      string `json:"fullName" validate:"required"`
	MobileNumber  string `json:"mobileNumber" validate:"required"`
	Address       string `json:"address"`
	ProvinceId    string `json:"provinceId" validate:"required"`
	CityId        string `json:"cityId" validate:"required"`
	DistrictId    string `json:"districtId" validate:"required"`
	SubdictrictId string `json:"subdictrictId" validate:"required"`
	CountryId     string `json:"countryId" validate:"required"`
	RtRw          string `json:"rtRw"`
	Role          string `json:"role"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
}

type VerifyRegisterUser struct {
	Email string `json:"email" validate:"required,min=1,max=50"`
	Otp   string `json:"otpNumber" validate:"required"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,min=1,max=50"`
	Password string `json:"password" validate:"required"`
}

type GetProfile struct {
	UserId string
}
