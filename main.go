package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

var oneDigitNames = map[int]string{
	0: "sıfır",
	1: "bir",
	2: "iki",
	3: "üc",
	4: "dörd",
	5: "beş",
	6: "altı",
	7: "yeddi",
	8: "səkkiz",
	9: "doqquz",
}

var twoDigitNames = map[int]string{
	1: "on",
	2: "iyirmi",
	3: "otuz",
	4: "qırx",
	5: "əlli",
	6: "altımış",
	7: "yetmiş",
	8: "səksən",
	9: "doxsan",
}

var threeDigitName = map[int]string{
	1: "yüz",
	2: "iki yüz",
	3: "üc yüz",
	4: "dörd yüz",
	5: "beş yüz",
	6: "altı yüz",
	7: "yeddi yüz",
	8: "səkkiz yüz",
	9: "doqquz yüz",
}

const fourDigitName = " min"
const sevenDigitName = " milyon"
const tenDigitsName = " milyard"

func addAdditionalWordInFrontIfNeeded(last, next int, current, additional string) string {
	if last == 1 && next == 0 {
		return additional + current
	}
	return current
}

func convert(num int) string {
	if num == 0 {
		return oneDigitNames[num]
	}
	return helper(num, true, 1, 1)
}

func helper(num int, isCategoryWordAdded bool, digitPosition int, digitCount int) string {
	cur := num % 10
	leftNum := num / 10

	if digitPosition > 3 {
		digitPosition = 1
		isCategoryWordAdded = false
	}

	if cur == 0 { //we just pass it to the next lvl
		return helper(leftNum, isCategoryWordAdded, digitPosition+1, digitCount+1)
	}

	var res string
	switch digitPosition {
	case 1:
		res = oneDigitNames[cur] //1
	case 2:
		res = twoDigitNames[cur] //11
	case 3:
		res = threeDigitName[cur] //111
		res = addAdditionalWordInFrontIfNeeded(cur, leftNum, res, "bir ")
	}

	if !isCategoryWordAdded && digitCount > 3 {
		switch digitCount {
		case 4, 5, 6: //1_000 ... 100_000
			res += fourDigitName
		case 7, 8, 9: //1_000_000 ... 100_000_000
			res += sevenDigitName
		case 10, 11, 12: //1_000_000_000 ... 100_000_000_000
			res += tenDigitsName
		}
		isCategoryWordAdded = true
	}

	if leftNum == 0 { //if cur digit is the last one
		return res
	}
	return helper(leftNum, isCategoryWordAdded, digitPosition+1, digitCount+1) + " " + res
}

func handler(w http.ResponseWriter, r *http.Request) {
	manatStr := r.URL.Query().Get("manat")
	qepikStr := r.URL.Query().Get("qepik")

	manat, err := strconv.Atoi(manatStr)
	if err != nil {
		http.Error(w, "Invalid manat value", http.StatusBadRequest)
		return
	}

	qepik, err := strconv.Atoi(qepikStr)
	if err != nil {
		http.Error(w, "Invalid qepik value", http.StatusBadRequest)
		return
	}

	res := convert(manat) + " manat"
	if qepik != 0 {
		res += " " + convert(qepik) + " qepik"
	}

	resp := &Response{
		Lang:   "az",
		Result: res,
	}

	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

type Response struct {
	Lang   string `json:"lang"`
	Result string `json:"result"`
}

func main() {
	//fmt.Printf("%s%c\n", convert(1), '*')               //ok
	//fmt.Printf("%s%c\n", convert(10), '*')              //1 2 3
	//fmt.Printf("%s%c\n", convert(100), '*')             //1 2 3
	//fmt.Printf("%s%c\n", convert(1_000), '*')           //1 2 3
	//fmt.Printf("%s%c\n", convert(10_000), '*')          //1 2 3
	//fmt.Printf("%s%c\n", convert(100_000), '*')         //1 2 3
	//fmt.Printf("%s%c\n", convert(1_000_000), '*')       //1 2 3
	//fmt.Printf("%s%c\n", convert(10_000_000), '*')      //1 2 3
	//fmt.Printf("%s%c\n", convert(100_000_000), '*')     //1 2 3
	//fmt.Printf("%s%c\n", convert(1_000_000_000), '*')   //1 2 3
	//fmt.Printf("%s%c\n", convert(10_000_000_000), '*')  //1 2 3
	//fmt.Printf("%s%c\n", convert(113_168_135_431), '*') //1 2 3
	//fmt.Printf("%s%c\n", convert(1123), '*')            //1 2 3
	//fmt.Printf("%s%c\n", convert(353), '*')             //1 2 3
	//fmt.Printf("%s%c\n", convert(349_111_123), '*')     //1 2 3
	//fmt.Printf("%s%c\n", convert(88_999_123), '*')      //1 2 3
	//fmt.Printf("%s%c\n", convert(0), '*')               //1 2 3
	//
	//num, err := strconv.Atoi("40")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(convert(num))

	//a, b := math.Modf(11.22)
	//fmt.Println(a)
	//fmt.Println(b)

	http.HandleFunc("/convert-num-to-az", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
