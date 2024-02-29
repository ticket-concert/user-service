package entity

type Province struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type City struct {
	Id           string `json:"id" bson:"id"`
	Name         string `json:"name" bson:"name"`
	ProvinceId   string `json:"provinceId" bson:"provinceId"`
	ProvinceName string `json:"provinceName" bson:"provinceName"`
}

type District struct {
	Id           string `json:"id" bson:"id"`
	Name         string `json:"name" bson:"name"`
	CityId       string `json:"cityId" bson:"cityId"`
	CityName     string `json:"cityName" bson:"cityName"`
	ProvinceId   string `json:"provinceId" bson:"provinceId"`
	ProvinceName string `json:"provinceName" bson:"provinceName"`
}

type SubDistrict struct {
	Id           string `json:"id" bson:"id"`
	Name         string `json:"name" bson:"name"`
	DistrictId   string `json:"districtId" bson:"districtId"`
	DistrictName string `json:"districtName" bson:"districtName"`
	CityId       string `json:"cityId" bson:"cityId"`
	CityName     string `json:"cityName" bson:"cityName"`
	ProvinceId   string `json:"provinceId" bson:"provinceId"`
	ProvinceName string `json:"provinceName" bson:"provinceName"`
}

type Country struct {
	Id            int    `json:"id" bson:"id"`
	Code          string `json:"code" bson:"code"`
	Name          string `json:"name" bson:"name"`
	Iso3          string `json:"iso3" bson:"iso3"`
	Number        int    `json:"number" bson:"number"`
	ContinentCode string `json:"continentCode" bson:"continentCode"`
	ContinentName string `json:"continentName" bson:"continentName"`
	DisplayOrder  int    `json:"displayOrder" bson:"displayOrder"`
	FullName      string `json:"fullName" bson:"fullName"`
}

type Continent struct {
	Code string `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
}
