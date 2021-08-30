package utils

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

var (
	allNumRe  = regexp.MustCompile("[^0-9]")
	allLetRe  = regexp.MustCompile("[^a-zA-Z]")
	allSame14 = regexp.MustCompile("0{14}|1{14}|2{14}|3{14}|4{14}|5{14}|6{14}|7{14}|8{14}|9{14}")
	allSame11 = regexp.MustCompile("0{11}|1{11}|2{11}|3{11}|4{11}|5{11}|6{11}|7{11}|8{11}|9{11}")
)

const (
	DiasCorridos = 0
	DiasUteis    = 1

	EDIPastaCompletos      = "completos"
	EDIPastaNaoProcessados = "naoprocessados"
	EDIPastaPartes         = "partes"
)

type Intervalo struct {
	Minuto    string
	Hora      string
	DiaSemana string
	DiaMes    string
}

func RespondWithError(w http.ResponseWriter, code, errorCode int, message string) {
	RespondWithJSON(w, code, map[string]interface{}{"error": message, "code": errorCode})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func Encrypt(pass string) string {
	return encrypt(pass, randomSalt())
}

func Matches(clear string, encrypted string) bool {
	parts := strings.Split(encrypted, "$")
	if len(parts) != 3 {
		return false
	}
	return encrypted == encrypt(clear, parts[2])
}

func randomSalt() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var salt string
	for i := 0; i <= 5; i++ {
		n := r.Intn(100)
		salt = salt + strconv.Itoa(n)
	}
	return salt
}

func encrypt(pass string, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt+pass)
	password := fmt.Sprintf("sha1$%x$%s", h.Sum(nil), salt)
	return password
}

func IsValidEmail(email string) bool {
	reMail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(email) > 254 || !reMail.MatchString(email) {
		return false
	}
	return true
}

func IsValidCNPJ(cnpj string) bool {
	n := allNumRe.ReplaceAllString(cnpj, "")
	if n == "" {
		return false
	}
	if len(n) != 14 {
		return false
	}
	if allSame14.Match([]byte(n)) {
		return false
	}
	if n[:3] == "000" {
		return false
	}
	size := len(n) - 2
	numbers := n[0:size]
	digits := n[size:]
	var sum int
	pos := size - 7
	for i := size; i >= 1; i-- {
		num, _ := strconv.Atoi(string(numbers[size-i]))
		sum += num * pos
		pos = pos - 1
		if pos < 2 {
			pos = 9
		}
	}
	var result int
	if sum%11 < 2 {
		result = 0
	} else {
		result = 11 - sum%11
	}
	x, _ := strconv.Atoi(string(digits[0]))
	if result != x {
		return false
	}
	size = size + 1
	numbers = n[0:size]
	sum = 0
	pos = size - 7
	for i := size; i >= 1; i-- {
		num, _ := strconv.Atoi(string(numbers[size-i]))
		sum += num * pos
		pos = pos - 1
		if pos < 2 {
			pos = 9
		}
	}
	if sum%11 < 2 {
		result = 0
	} else {
		result = 11 - sum%11
	}
	num, _ := strconv.Atoi(string(digits[1]))
	if result != num {
		return false
	}
	return true
}

func IsValidCPF(cpf string) bool {
	cpf = allNumRe.ReplaceAllString(cpf, "")
	if cpf == "" {
		return false
	}
	if len(cpf) != 11 {
		return false
	}
	if allSame11.Match([]byte(cpf)) {
		return false
	}
	var sum int
	var res int
	for i := 1; i <= 9; i++ {
		num, _ := strconv.Atoi(cpf[i-1 : i])
		sum = sum + num*(11-i)
	}
	res = (sum * 10) % 11
	if (res == 10) || (res == 11) {
		res = 0
	}
	num, _ := strconv.Atoi(cpf[9:10])
	if res != num {
		return false
	}
	sum = 0
	for i := 1; i <= 10; i++ {
		num, _ := strconv.Atoi(cpf[i-1 : i])
		sum = sum + num*(12-i)
	}
	res = (sum * 10) % 11
	if (res == 10) || (res == 11) {
		res = 0
	}
	num, _ = strconv.Atoi(cpf[10:11])
	if res != num {
		return false
	}
	return true
}

func NormalizeCPFCNPJ(cpfcnpj string) (string, bool) {

	cpfcnpj = OnlyNumbers(cpfcnpj)
	if IsValidCNPJ(cpfcnpj) {
		return cpfcnpj, true
	}

	if len(cpfcnpj) > 11 {
		cpfcnpj = cpfcnpj[3:]
	}
	if IsValidCPF(cpfcnpj) {
		return cpfcnpj, true
	}

	return cpfcnpj, false

}

func OnlyNumbers(s string) string {
	return allNumRe.ReplaceAllString(s, "")
}

func OnlyLetters(s string) string {
	return allLetRe.ReplaceAllString(s, "")
}

func IsNumber(n string) bool {
	_, err := strconv.ParseInt(n, 10, 64)
	return err == nil
}

func InArray(s string, list []string) bool {
	for _, el := range list {
		if el == s {
			return true
		}
	}
	return false
}

func CompareSlices(slice1 []int64, slice2 []int64) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := 0; i < len(slice1); i++ {
		if !InIntArray(slice1[i], slice2) {
			return false
		}
	}
	return true
}

func GetDiffSlices(slice1 []int64, slice2 []int64) []int64 {
	var diff []int64

	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			if found {
				diff = append(diff, s1)
			}
		}
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}
	return diff
}

func InIntArray(n int64, list []int64) bool {
	for _, el := range list {
		if el == n {
			return true
		}
	}
	return false
}

func Abrev(texto string, tamanho int) string {

	var newString string = ""

	if len(texto) > tamanho {
		palavras := strings.Split(texto, " ")
		for i, v := range palavras {
			if i > 0 {
				if len(v) > 3 {
					v = fmt.Sprint(v[0:4], ".")
				} else {
					v = fmt.Sprint(v, " ")
				}
				newString += v
			} else {
				newString += fmt.Sprint(v, " ")
			}
		}
		return newString
	} else {
		return texto
	}

}

func Modulo11(numero string) (int64, error) {

	var (
		soma          int64
		multiplicador int64 = 2
		digito        int64
	)

	for i := int64(len(numero)) - 1; i >= 0; i-- {

		if multiplicador == 10 {
			multiplicador = 2
		}

		codigo, err := strconv.ParseInt(string(numero[i]), 10, 64)
		if err != nil {
			return 0, err
		}

		soma += (codigo * multiplicador)
		multiplicador++
	}

	resto := soma % 11

	if resto == 0 || resto == 1 {
		digito = 0
	} else {
		digito = 11 - resto
	}

	return digito, nil

}

func GetAliquota(codMunIni int64, codMunFim int64, ie string) (int64, error) {

	var aliquota int64
	sulSudesteExcetoES := []string{"31", "33", "35", "41", "42", "43"}
	norteNordesteCentroOesteES := []string{"11", "12", "13", "14", "15", "16", "17", "21", "22", "23", "24", "25", "26", "27", "28", "29", "50", "51", "52", "53", "32"}

	if codMunIni < 99 || codMunFim < 99 {
		return aliquota, fmt.Errorf("Codigo de municipio não informado: ini %d, fim: %d", codMunIni, codMunFim)
	}

	strCodMunIni := fmt.Sprint(codMunIni)[:2]
	strCodMunFim := fmt.Sprint(codMunFim)[:2]

	switch strCodMunIni {
	case "35":
		if strCodMunFim == "35" {
			aliquota = 12
		} else if InArray(strCodMunIni, sulSudesteExcetoES) {
			aliquota = 12
		} else if InArray(strCodMunFim, norteNordesteCentroOesteES) {
			aliquota = 7
		}
	case "52":
		if ie == "ISENTO" && strCodMunFim == "52" {
			aliquota = 17
		} else if ie != "ISENTO" && strCodMunFim == "52" {
			aliquota = 0
		} else if InArray(strCodMunFim, sulSudesteExcetoES) {
			aliquota = 12
		} else if InArray(strCodMunFim, norteNordesteCentroOesteES) {
			aliquota = 12
		}
	case "53":
		if strCodMunFim == "53" {
			aliquota = 0
		} else if InArray(strCodMunFim, sulSudesteExcetoES) {
			aliquota = 12
		} else if InArray(strCodMunFim, norteNordesteCentroOesteES) {
			aliquota = 12
		}
	case "32":
		if strCodMunFim == "32" {
			aliquota = 12
		} else if InArray(strCodMunFim, sulSudesteExcetoES) {
			aliquota = 12
		} else if InArray(strCodMunFim, norteNordesteCentroOesteES) {
			aliquota = 12
		}
	case "31":
		if strCodMunFim == "31" {
			aliquota = 0
		} else if InArray(strCodMunFim, sulSudesteExcetoES) {
			aliquota = 12
		} else if InArray(strCodMunFim, norteNordesteCentroOesteES) {
			aliquota = 7
		}
	case "33":
		if strCodMunFim == "33" {
			aliquota = 12
		} else if InArray(strCodMunFim, sulSudesteExcetoES) {
			aliquota = 12
		} else if InArray(strCodMunFim, norteNordesteCentroOesteES) {
			aliquota = 7
		}
	case "41":
		if strCodMunFim == "41" {
			aliquota = 0
		} else if InArray(strCodMunFim, sulSudesteExcetoES) {
			aliquota = 12
		} else if InArray(strCodMunFim, norteNordesteCentroOesteES) {
			aliquota = 7
		}
	case "42":
		if InArray(strCodMunFim, sulSudesteExcetoES) {
			aliquota = 12
		} else if InArray(strCodMunFim, norteNordesteCentroOesteES) {
			aliquota = 7
		}
	case "43":
		if strCodMunFim == "43" {
			aliquota = 0
		} else {
			aliquota = 12
		}
	default:
		aliquota = 0
	}

	return aliquota, nil

}

func RandomNumber(number int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(number)
}

func FloatToCurrency(f float64) int64 {
	valorFormatado := OnlyNumbers(fmt.Sprintf("%.2f", f))
	currency, _ := strconv.ParseInt(valorFormatado, 10, 64)
	return currency
}

func CurrencyToFloat(d int64) float64 {
	return float64(d) / 100
}

func MilimetroToMetro(d int64) float64 {
	return float64(d) / 1000
}

func ToFixed(val float64, dec int) float64 {
	pow := math.Pow(10, float64(dec))
	digit := pow * val
	round := math.Floor(digit)

	return round / pow
}

func GetStringInBetween(str string, start string, end string) string {
	s := strings.Index(str, start)
	if s == -1 {
		return ""
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return ""
	}
	e += s
	if s > e {
		return ""
	}
	return str[s:e]
}

type JsonSpecialDate struct {
	time.Time
}

func (sd *JsonSpecialDate) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)
	if strInput == "0000-00-00" {
		sd.Time = time.Time{}
		return nil
	}
	newTime, err := time.Parse("2006-01-02", strInput)
	if err != nil {
		return err
	}
	sd.Time = newTime
	return nil
}

func (sd *JsonSpecialDate) MarshalJSON() ([]byte, error) {
	if sd.Time.IsZero() {
		return []byte("\"0001-01-01\""), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", sd.Time.Format("2006-01-02"))), nil
}

type JsonSpecialDateTime struct {
	time.Time
}

func (sd *JsonSpecialDateTime) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)
	if strInput == "0000-00-00 00:00:00" {
		sd.Time = time.Time{}
		return nil
	}
	newTime, err := time.Parse("2006-01-02 15:04:05", strInput)
	if err != nil {
		return err
	}
	sd.Time = newTime
	return nil
}

func (sd *JsonSpecialDateTime) MarshalJSON() ([]byte, error) {
	if sd.Time.IsZero() {
		return []byte("\"0001-01-01 00:00:00\""), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", sd.Time.Format("2006-01-02 15:04:05"))), nil
}

func GetValorCampo(campo string, structValor reflect.Value) (tipo string, valor string, err error) {
	if campo == "" {
		err = fmt.Errorf("Campo não informado")
		return
	}
	if strings.Contains(structValor.Type().String(), "*") {
		err = fmt.Errorf("Não permitido acessar estrutura tipo ponteiro")
		return
	}

	if structValor.Type().Kind().String() == "slice" {
		var valores []string
		for i := 0; i < structValor.Len(); i++ {
			tipoRecursivo, valorRecursivo, errRecursivo := GetValorCampo(campo, structValor.Index(i))
			tipo = tipoRecursivo

			if errRecursivo != nil {
				err = errRecursivo
				return
			}

			valores = append(valores, valorRecursivo)
		}
		valor = strings.Join(valores, "|*")
		return
	}

	atributos := strings.Split(campo, ".")
	for i := 0; i < len(atributos); i++ {
		atributo := atributos[i]

		if structValor.FieldByName(atributo).Kind() == reflect.Invalid {
			err = fmt.Errorf("Entre os campos %s o atributo %s não encontrado na Struct\nstruct:%#v", campo, atributo, structValor)
			return
		}

		structValor = structValor.FieldByName(atributo)

		tipo = structValor.Type().Kind().String()
		switch tipo {
		case "string":
			valor = structValor.String()

		case "int64", "int32", "int":
			valor = fmt.Sprintf("%d", structValor.Int())

		case "float64", "float32":
			valor = fmt.Sprintf("%.2f", structValor.Float())

		case "struct":
			// Tratativa para estruturas jsonSpecialDate e jsonSpecialDateTime
			if structValor.NumField() == 1 && structValor.Field(0).Type().Name() == "Time" {
				if structValor.Field(0).Type().ConvertibleTo(reflect.TypeOf(time.Time{})) {
					tipo = "time"
					valor = structValor.Field(0).Interface().(time.Time).Format("2006-01-02 15:04:05")
				}
			}
			if structValor.Type().ConvertibleTo(reflect.TypeOf(time.Time{})) {
				tipo = "time"
				valor = structValor.Interface().(time.Time).Format("2006-01-02 15:04:05")
			}

		case "slice":
			if structValor.Len() == 0 {
				indx := reflect.MakeSlice(structValor.Type(), 1, 1)
				structValor = reflect.Append(structValor, indx.Index(0))
			}

			var valores []string
			for j := 0; j < structValor.Len(); j++ {
				if len(atributos) == i+1 {
					err = fmt.Errorf("Campo Slice não possui atributo para receber valor")
					return
				}

				proximosCampos := strings.Join(atributos[i+1:], ".")
				tipoRecursivo, valorRecursivo, errRecursivo := GetValorCampo(proximosCampos, structValor.Index(j))
				if errRecursivo != nil {
					err = errRecursivo
					return
				}

				tipo = tipoRecursivo

				valores = append(valores, valorRecursivo)
			}
			valor = strings.Join(valores, "|*")
			return
		}
	}
	return
}

func SetValorCampo(campo string, estrutura interface{}, valor interface{}) (err error) {
	rStruct := reflect.ValueOf(estrutura)
	rValor := reflect.ValueOf(valor)

	if rStruct.Kind() != reflect.Ptr {
		return fmt.Errorf("Tipo inválido")
	}

	return setValorCampo(campo, rStruct.Elem(), rValor)
}

func setValorCampo(campo string, estrutura reflect.Value, valor reflect.Value) error {
	atributos := strings.Split(campo, ".")
	var estruturaLocal reflect.Value

	for i, atributo := range atributos {

		if estruturaLocal = estrutura.FieldByName(atributo); !estruturaLocal.CanSet() {
			return fmt.Errorf("can't set %#v", estruturaLocal)
		}

		switch estruturaLocal.Type().Kind() {
		case reflect.Struct:
			if !estruturaLocal.Type().ConvertibleTo(reflect.TypeOf(time.Time{})) {
				return setValorCampo(strings.Join(atributos[i+1:], "."), estruturaLocal, valor)
			}
		case reflect.Slice:
			if estruturaLocal.Len() == 0 {
				indx := reflect.MakeSlice(estruturaLocal.Type(), 1, 1)
				estruturaLocal.Set(indx)
			}

			if i == len(atributos)-1 {
				break
			}

			for j := 0; j < estruturaLocal.Len(); j++ {
				err := setValorCampo(strings.Join(atributos[i+1:], "."), estruturaLocal.Index(j), valor)
				if err != nil {
					return err
				}
			}
			return nil
		}

		estruturaLocal.Set(valor)
	}
	return nil
}

func JSONUnmarshalValidate(jsonString string, estrutura reflect.Type) error {

	var estruturaJson interface{}
	if err := json.Unmarshal([]byte(jsonString), &estruturaJson); err != nil {
		return err
	}

	switch reflect.TypeOf(estruturaJson).Kind() {
	case reflect.Slice:
		for _, estruturaSlice := range estruturaJson.([]interface{}) {
			estruturaMap := estruturaSlice.(map[string]interface{})
			err := validaMapaEstrutura(estruturaMap, estrutura)
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		estruturaMap := estruturaJson.(map[string]interface{})
		return validaMapaEstrutura(estruturaMap, estrutura)
	default:
		return fmt.Errorf("JSON mal-formado")
	}
	return nil
}

func validaMapaEstrutura(mapa map[string]interface{}, estrutura reflect.Type) error {

	var erros []string
	for campo, valor := range mapa {

		structField, ok := estrutura.FieldByName(campo)
		if !ok {
			erros = append(erros, fmt.Sprintf("Campo %s não localizado na estrutura %s", campo, estrutura.Name()))
			continue
		}

		tipoEstrutura := structField.Type.Kind().String()
		if valor == nil && InArray(tipoEstrutura, []string{"map", "slice"}) {
			continue
		} else if valor == nil {
			erros = append(erros, fmt.Sprintf("Campo %s enviado nulo e é esperado o tipo %s", campo, tipoEstrutura))
			continue
		}

		tipoValor := reflect.TypeOf(valor).Kind().String()

		// Estrutura é int, porém enviou float, então precisa validar se o valor é "convertível" para int
		if (tipoEstrutura == "int" || tipoEstrutura == "int32" || tipoEstrutura == "int64") && (tipoValor == "float32" || tipoValor == "float64") {
			floatValor, err := strconv.ParseFloat(fmt.Sprint(valor), 64)
			if err != nil || !isIntegral(floatValor) {
				erros = append(erros, fmt.Sprintf("Campo %s enviado no tipo %s e é esperado %s", campo, tipoValor, tipoEstrutura))
				continue
			}
		} else if tipoEstrutura == "float32" && tipoValor == "float64" {
			_, err := strconv.ParseFloat(fmt.Sprint(valor), 32)
			if err != nil {
				erros = append(erros, fmt.Sprintf("Campo %s enviado no tipo %s e é esperado %s", campo, tipoValor, tipoEstrutura))
				continue
			}
		} else if tipoEstrutura == "float64" && tipoValor == "float32" {
			_, err := strconv.ParseFloat(fmt.Sprint(valor), 64)
			if err != nil {
				erros = append(erros, fmt.Sprintf("Campo %s enviado no tipo %s e é esperado %s", campo, tipoValor, tipoEstrutura))
				continue
			}
		} else if tipoEstrutura == "bool" && tipoEstrutura != tipoValor {
			_, err := strconv.ParseBool(fmt.Sprint(valor))
			if err != nil {
				erros = append(erros, fmt.Sprintf("Campo %s enviado no tipo %s e é esperado %s", campo, tipoValor, tipoEstrutura))
				continue
			}
		} else if structField.Type.Kind() == reflect.Struct && structField.Type == reflect.TypeOf(time.Time{}) {
			if _, err := time.Parse("2006-01-02T15:04:05Z", fmt.Sprint(valor)); err != nil {
				erros = append(erros, fmt.Sprintf("Campo %s enviado com formato de time inválido [YYYY-mm-ddTHH:mm:ssZ", campo))
				continue
			}
		} else if structField.Type.Kind() == reflect.Struct && structField.Type == reflect.TypeOf(JsonSpecialDate{}) {
			if _, err := time.Parse("2006-01-02", fmt.Sprint(valor)); err != nil {
				erros = append(erros, fmt.Sprintf("Campo %s enviado com formato de date inválido [YYYY-mm-dd].", campo))
				continue
			}
		} else if structField.Type.Kind() == reflect.Struct && structField.Type == reflect.TypeOf(JsonSpecialDateTime{}) {
			if _, err := time.Parse("2006-01-02 15:04:05", fmt.Sprint(valor)); err != nil {
				erros = append(erros, fmt.Sprintf("Campo %s enviado com formato de datetime inválido [YYYY-mm-dd HH:mm:ss].", campo))
				continue
			}
		} else if tipoEstrutura == "struct" && tipoValor == "map" {
			err := validaMapaEstrutura(valor.(map[string]interface{}), structField.Type)
			if err != nil {
				erros = append(erros, err.Error())
				continue
			}
		} else if tipoEstrutura == "slice" && tipoValor == "slice" {
			for _, subAtributo := range valor.([]interface{}) {
				tipoValorSlice := reflect.TypeOf(subAtributo).Kind().String()
				if tipoValorSlice == "map" {
					err := validaMapaEstrutura(subAtributo.(map[string]interface{}), structField.Type.Elem())
					if err != nil {
						erros = append(erros, err.Error())
						continue
					}
				}
			}
		} else if tipoEstrutura != tipoValor {
			erros = append(erros, fmt.Sprintf("Campo %s enviado no tipo %s e é esperado %s", campo, tipoValor, tipoEstrutura))
			continue
		} else if !InArray(tipoValor, []string{"string", "int", "int32", "int64", "bool", "float32", "float64", "map", "slice"}) {
			erros = append(erros, fmt.Sprintf("Campo %s enviado em um tipo não tratado [%s]", campo, tipoValor))
			continue
		}

	}

	if len(erros) > 0 {
		return fmt.Errorf("%s", strings.Join(erros, "\n"))
	}
	return nil
}

func isIntegral(val float64) bool {
	return val == float64(int(val))
}

func LimparString(s string) string {
	var (
		b           = bytes.NewBufferString("")
		acentuacoes = map[rune]string{'À': "A", 'Á': "A", 'Â': "A", 'Ã': "A", 'Ä': "A", 'Å': "AA", 'Æ': "AE", 'Ç': "C", 'È': "E", 'É': "E", 'Ê': "E", 'Ë': "E", 'Ì': "I", 'Í': "I", 'Î': "I", 'Ï': "I", 'Ð': "D", 'Ł': "L", 'Ñ': "N", 'Ò': "O", 'Ó': "O", 'Ô': "O", 'Õ': "O", 'Ö': "OE", 'Ø': "OE", 'Œ': "OE", 'Ù': "U", 'Ú': "U", 'Ü': "UE", 'Û': "U", 'Ý': "Y", 'Þ': "TH", 'ẞ': "SS", 'à': "a", 'á': "a", 'â': "a", 'ã': "a", 'ä': "ae", 'å': "aa", 'æ': "ae", 'ç': "c", 'è': "e", 'é': "e", 'ê': "e", 'ë': "e", 'ì': "i", 'í': "i", 'î': "i", 'ï': "i", 'ð': "d", 'ł': "l", 'ñ': "n", 'ń': "n", 'ò': "o", 'ó': "o", 'ô': "o", 'õ': "o", 'ō': "o", 'ö': "oe", 'ø': "oe", 'œ': "oe", 'ś': "s", 'ù': "u", 'ú': "u", 'û': "u", 'ū': "u", 'ü': "ue", 'ý': "y", 'ÿ': "y", 'ż': "z", 'þ': "th", 'ß': "ss"}
	)
	for _, c := range s {
		if val, ok := acentuacoes[c]; ok {
			b.WriteString(val)
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func IncrementaUltimaLetra(texto string) string {
	if len(texto)-1 < 0 {
		return texto
	}
	ultimaLetra := []rune(texto[len(texto)-1 : len(texto)])
	if len(ultimaLetra) < 0 {
		return texto
	}
	proximaLetra := string(ultimaLetra[0] + 1)
	if OnlyLetters(proximaLetra) != "" {
		return fmt.Sprintf("%s%s", texto[:len(texto)-1], proximaLetra)
	}
	return texto
}

func GetTimeNow() time.Time {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return time.Now()
	}
	return time.Now().In(loc)
}

func GetSpecialTimeNow() JsonSpecialDateTime {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return JsonSpecialDateTime{Time: time.Now()}
	}
	return JsonSpecialDateTime{Time: time.Now().In(loc)}
}

func GetSpecialDateNow() JsonSpecialDate {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return JsonSpecialDate{Time: time.Now()}
	}
	return JsonSpecialDate{Time: time.Now().In(loc)}
}

func ValidaExecucao(intervalo Intervalo) bool {
	now := GetTimeNow()
	var validaMinuto bool
	var validaHora bool
	var validaSemana bool
	var validaDiaMes bool
	intervalo.Minuto = strings.TrimSpace(intervalo.Minuto)
	intervalo.Hora = strings.TrimSpace(intervalo.Hora)
	intervalo.DiaSemana = strings.TrimSpace(intervalo.DiaSemana)
	intervalo.DiaMes = strings.TrimSpace(intervalo.DiaMes)
	if intervalo.Minuto == "*" {
		validaMinuto = true
	} else {
		minutoAtual, err := strconv.ParseInt(now.Format("04"), 10, 64)
		if err != nil {
			return false
		}
		if strings.Contains(intervalo.Minuto, "-") { //Até
			minutoDeAte := strings.Split(intervalo.Minuto, "-")
			if len(minutoDeAte) != 2 {
				return false
			}
			minutoDe, err := strconv.ParseInt(minutoDeAte[0], 10, 64)
			if err != nil {
				return false
			}
			minutoAte, err := strconv.ParseInt(minutoDeAte[1], 10, 64)
			if err != nil {
				return false
			}
			if minutoAtual >= minutoDe && minutoAtual <= minutoAte {
				validaMinuto = true
			}
		} else if strings.Contains(intervalo.Minuto, ",") { //separados
			minutosSeparados := strings.Split(intervalo.Minuto, ",")
			for _, strMinutoIntervalo := range minutosSeparados {
				minutoIntervalo, err := strconv.ParseInt(strMinutoIntervalo, 10, 64)
				if err != nil {
					return false
				}
				if minutoIntervalo == minutoAtual {
					validaMinuto = true
				}
			}
		} else {
			minutoIntervalo, err := strconv.ParseInt(intervalo.Minuto, 10, 64)
			if err != nil {
				return false
			}
			if minutoIntervalo == minutoAtual {
				validaMinuto = true
			}
		}
	}
	if intervalo.Hora == "*" {
		validaHora = true
	} else {
		horaAtual, err := strconv.ParseInt(now.Format("15"), 10, 64)
		if err != nil {
			return false
		}
		if strings.Contains(intervalo.Hora, "-") { //Até
			horaDeAte := strings.Split(intervalo.Hora, "-")
			if len(horaDeAte) != 2 {
				return false
			}
			horaDe, err := strconv.ParseInt(horaDeAte[0], 10, 64)
			if err != nil {
				return false
			}
			horaAte, err := strconv.ParseInt(horaDeAte[1], 10, 64)
			if err != nil {
				return false
			}
			if horaAtual >= horaDe && horaAtual <= horaAte {
				validaHora = true
			}
		} else if strings.Contains(intervalo.Hora, ",") { //separados
			horaSeparada := strings.Split(intervalo.Hora, ",")
			for _, strHoraIntervalo := range horaSeparada {
				horaIntervalo, err := strconv.ParseInt(strHoraIntervalo, 10, 64)
				if err != nil {
					return false
				}
				if horaIntervalo == horaAtual {
					validaHora = true
				}
			}
		} else {
			horaIntervalo, err := strconv.ParseInt(intervalo.Hora, 10, 64)
			if err != nil {
				return false
			}
			if horaIntervalo == horaAtual {
				validaHora = true
			}
		}
	}
	if intervalo.DiaSemana == "*" {
		validaSemana = true
	} else {
		diaSemanaAtual := int64(now.Weekday())
		if strings.Contains(intervalo.DiaSemana, "-") { //Até
			diaSemanaDeAte := strings.Split(intervalo.DiaSemana, "-")
			if len(diaSemanaDeAte) != 2 {
				return false
			}
			diaSemanaDe, err := strconv.ParseInt(diaSemanaDeAte[0], 10, 64)
			if err != nil {
				return false
			}
			diaSemanaAte, err := strconv.ParseInt(diaSemanaDeAte[1], 10, 64)
			if err != nil {
				return false
			}
			if diaSemanaAtual >= diaSemanaDe && diaSemanaAtual <= diaSemanaAte {
				validaSemana = true
			}
		} else if strings.Contains(intervalo.DiaSemana, ",") {
			diaSemanaSeparado := strings.Split(intervalo.DiaSemana, ",")
			for _, strDiaSemanaIntervalo := range diaSemanaSeparado {
				diaSemanaIntervalo, err := strconv.ParseInt(strDiaSemanaIntervalo, 10, 64)
				if err != nil {
					return false
				}
				if diaSemanaIntervalo == diaSemanaAtual {
					validaSemana = true
				}
			}
		} else {
			diaSemanaIntervalo, err := strconv.ParseInt(intervalo.DiaSemana, 10, 64)
			if err != nil {
				return false
			}
			if diaSemanaIntervalo == diaSemanaAtual {
				validaSemana = true
			}
		}
	}
	if intervalo.DiaMes == "*" {
		validaDiaMes = true
	} else {
		diaMesAtual, err := strconv.ParseInt(now.Format("02"), 10, 64)
		if err != nil {
			return false
		}
		if strings.Contains(intervalo.DiaMes, "-") { //Até
			diaMesDeAte := strings.Split(intervalo.DiaMes, "-")
			if len(diaMesDeAte) != 2 {
				return false
			}
			diaMesDe, err := strconv.ParseInt(diaMesDeAte[0], 10, 64)
			if err != nil {
				return false
			}
			diaMesAte, err := strconv.ParseInt(diaMesDeAte[1], 10, 64)
			if err != nil {
				return false
			}
			if diaMesAtual >= diaMesDe && diaMesAtual <= diaMesAte {
				validaDiaMes = true
			}
		} else if strings.Contains(intervalo.DiaMes, ",") { //separados
			diaMesSeparado := strings.Split(intervalo.DiaMes, ",")
			for _, strDiaMesIntervalo := range diaMesSeparado {
				diaMesIntervalo, err := strconv.ParseInt(strDiaMesIntervalo, 10, 64)
				if err != nil {
					return false
				}
				if diaMesIntervalo == diaMesAtual {
					validaDiaMes = true
				}
			}
		} else {
			diaMesIntervalo, err := strconv.ParseInt(intervalo.DiaMes, 10, 64)
			if err != nil {
				return false
			}
			if diaMesIntervalo == diaMesAtual {
				validaDiaMes = true
			}
		}
	}
	return validaMinuto && validaHora && validaDiaMes && validaSemana
}

func GetDadosBase64(fullBase64 string) (string, string, error) {
	regBase64 := regexp.MustCompile("^data:(.*?);base64,")

	if fullBase64 == "" {
		return "", "", fmt.Errorf("Base64 não informado")
	}

	var mimeType string
	aMimeType := regBase64.FindStringSubmatch(fullBase64)
	if len(aMimeType) == 2 {
		mimeType = aMimeType[1]
	}

	base64Content := regBase64.ReplaceAllString(fullBase64, "")
	if _, err := base64.StdEncoding.DecodeString(base64Content); err != nil {
		return "", "", fmt.Errorf("Base64 inválido: %v", err)
	}

	return base64Content, mimeType, nil
}

func GeraPDF(filename string, baseImagem string) error {

	base64Content, mimeType, err := GetDadosBase64(baseImagem)
	if err != nil {
		return fmt.Errorf("Erro ao pegar os dados do base64 %v", err)
	}

	dec, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return fmt.Errorf("Erro ao realizar decode do base64 %v", err)
	}

	extension, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return fmt.Errorf("Falha ao obter a extensão do arquivo %v", err)
	}

	if len(extension) < 1 {
		return fmt.Errorf("Tipo de imagem inválido")
	}

	f, err := os.Create(fmt.Sprintf("/tmp/filename%v", extension[0]))
	if err != nil {
		return fmt.Errorf("Erro ao criar arquivo %v", err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		return fmt.Errorf("Erro ao escrever arquivo")
	}
	if err := f.Sync(); err != nil {
		return fmt.Errorf("Erro ao salvar arquivo")
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.ImageOptions(
		f.Name(),
		40, 10,
		100, 210,
		false,
		gofpdf.ImageOptions{ReadDpi: true},
		0,
		"",
	)
	err = pdf.OutputFileAndClose(fmt.Sprintf("/tmp/%v", filename))

	return err
}

func ParseJsonSpecialDate(layout string, unparsedData string) (JsonSpecialDate, error) {

	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return JsonSpecialDate{}, err
	}

	data, err := time.ParseInLocation(layout, unparsedData, loc)
	if err != nil {
		return JsonSpecialDate{}, err
	}

	return JsonSpecialDate{Time: data}, nil
}

func ParseJsonSpecialDateTime(layout string, unparsedData string) (JsonSpecialDateTime, error) {

	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return JsonSpecialDateTime{}, err
	}

	data, err := time.ParseInLocation(layout, unparsedData, loc)
	if err != nil {
		return JsonSpecialDateTime{}, err
	}

	return JsonSpecialDateTime{Time: data}, nil
}

func RemoveDuplicidadeLista(list []string) []string {
	check := make(map[string]int64, 0)

	if len(list) == 0 {
		return list
	}
	for i := range list {
		check[list[i]] = 1
	}
	list = []string{}
	for apelido := range check {
		list = append(list, apelido)
	}
	return list
}
