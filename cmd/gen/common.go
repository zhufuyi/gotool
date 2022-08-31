package gen

import (
	"github.com/zhufuyi/goctl/pkg/replacer"
	"github.com/zhufuyi/goctl/pkg/utils"
	"strings"
)

var (
	// 清除标记的模板代码片段标记
	startMark = []byte("// delete the templates code start")
	endMark   = []byte("// delete the templates code end")
)

func adjustmentOfIDType(handlerCodes string) string {
	return idTypeToStr(idTypefixToUint64(handlerCodes))
}

func idTypefixToUint64(handlerCodes string) string {
	subStart := "ByIDRequest struct {"
	subEnd := "`" + `json:"id" binding:""` + "`"
	if subBytes := utils.FindSubBytesNotIn([]byte(handlerCodes), []byte(subStart), []byte(subEnd)); len(subBytes) > 0 {
		old := subStart + string(subBytes) + subEnd
		newStr := subStart + "\n\tID uint64 " + subEnd + " // uint64 id\n"
		handlerCodes = strings.ReplaceAll(handlerCodes, old, newStr)
	}

	return handlerCodes
}

func idTypeToStr(handlerCodes string) string {
	subStart := "ByIDRespond struct {"
	subEnd := "`" + `json:"id"` + "`"
	if subBytes := utils.FindSubBytesNotIn([]byte(handlerCodes), []byte(subStart), []byte(subEnd)); len(subBytes) > 0 {
		old := subStart + string(subBytes) + subEnd
		newStr := subStart + "\n\tID string " + subEnd + " // covert to string id\n"
		handlerCodes = strings.ReplaceAll(handlerCodes, old, newStr)
	}

	return handlerCodes
}

func addTheDeleteFields(r replacer.Replacer, filename string) []replacer.Field {
	var fields []replacer.Field

	data, err := r.ReadFile(filename)
	if err != nil {
		return fields
	}
	if subBytes := utils.FindSubBytes(data, startMark, endMark); len(subBytes) > 0 {
		fields = append(fields,
			replacer.Field{ // 清除标记的模板代码
				Old: string(subBytes),
				New: "",
			},
		)
	}

	return fields
}
