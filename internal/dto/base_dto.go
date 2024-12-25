package dto

type BaseDto struct {
	Skip int32 `json:"skip"`
	Take int32 `json:"take"`
}

func NewBaseDto(skip, take int32) BaseDto {
	return BaseDto{
		Skip: skip,
		Take: take,
	}
}

func (b *BaseDto) GetSkip() int32 {
	if b.Skip == 0 {
		return 0
	}
	return b.Skip
}

func (b *BaseDto) GetTake() int32 {
	if b.Take == 0 {
		return 10
	}
	return b.Take
}
