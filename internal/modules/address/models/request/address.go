package request

type Province struct {
	Page   int64  `query:"page" validate:"required"`
	Size   int64  `query:"size" validate:"required"`
	Search string `query:"search"`
}

type City struct {
	ProvinceId string `query:"provinceId" validate:"required"`
	Page       int64  `query:"page" validate:"required"`
	Size       int64  `query:"size" validate:"required"`
	Search     string `query:"search"`
}

type District struct {
	ProvinceId string `query:"provinceId" validate:"required"`
	CityId     string `query:"cityId" validate:"required"`
	Page       int64  `query:"page" validate:"required"`
	Size       int64  `query:"size" validate:"required"`
	Search     string `query:"search"`
}

type SubDistrict struct {
	ProvinceId string `query:"provinceId" validate:"required"`
	CityId     string `query:"cityId" validate:"required"`
	DistrictId string `query:"districtId" validate:"required"`
	Page       int64  `query:"page" validate:"required"`
	Size       int64  `query:"size" validate:"required"`
	Search     string `query:"search"`
}

type Country struct {
	Page   int64  `query:"page" validate:"required"`
	Size   int64  `query:"size" validate:"required"`
	Search string `query:"search"`
}
