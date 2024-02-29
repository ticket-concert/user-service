package response

import (
	"user-service/internal/pkg/constants"
)

type Province struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ProvinceResp struct {
	CollectionData []Province
	MetaData       constants.MetaData
}

type City struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	ProvinceId   string `json:"provinceId"`
	ProvinceName string `json:"provinceName"`
}

type CityResp struct {
	CollectionData []City
	MetaData       constants.MetaData
}

type District struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	CityId       string `json:"cityId"`
	CityName     string `json:"cityName"`
	ProvinceId   string `json:"provinceId"`
	ProvinceName string `json:"provinceName"`
}

type DistrictResp struct {
	CollectionData []District
	MetaData       constants.MetaData
}

type SubDistrict struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	DistrictId   string `json:"districtId"`
	DistrictName string `json:"districtName"`
	CityId       string `json:"cityId"`
	CityName     string `json:"cityName"`
	ProvinceId   string `json:"provinceId"`
	ProvinceName string `json:"provinceName"`
}

type SubDistrictResp struct {
	CollectionData []SubDistrict
	MetaData       constants.MetaData
}

type Country struct {
	Id            int    `json:"id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	ContinentCode string `json:"continentCode"`
	ContinentName string `json:"continentName"`
	FullName      string `json:"fullName"`
}

type CountryResp struct {
	CollectionData []Country
	MetaData       constants.MetaData
}

type Continent struct {
	Code string `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
}

type ContinentResp struct {
	CollectionData []Continent
	MetaData       constants.MetaData
}
