package MakeKotlinApi

import (
	"GoAutoWeb/MakeKotlinApi/MakeApiConst"
	"GoAutoWeb/MakeKotlinApi/MakeApiRequest"
	"GoAutoWeb/MakeKotlinApi/MakeModel"
)

func Make() {
	MakeApiConst.Make()
	MakeModel.Make()
	MakeApiRequest.Make()
}
