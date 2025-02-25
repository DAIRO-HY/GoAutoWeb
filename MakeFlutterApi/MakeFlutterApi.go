package MakeFlutterApi

import (
	"GoAutoWeb/MakeFlutterApi/MakeApiConst"
	"GoAutoWeb/MakeFlutterApi/MakeApiRequest"
	"GoAutoWeb/MakeFlutterApi/MakeModel"
)

func Make() {
	MakeApiConst.Make()
	MakeModel.Make()
	MakeApiRequest.Make()
}
