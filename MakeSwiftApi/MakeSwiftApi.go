package MakeSwiftApi

import (
	"GoAutoWeb/MakeSwiftApi/MakeApiConst"
	"GoAutoWeb/MakeSwiftApi/MakeApiRequest"
	"GoAutoWeb/MakeSwiftApi/MakeModel"
)

func Make() {
	MakeApiConst.Make()
	MakeModel.Make()
	MakeApiRequest.Make()
}
