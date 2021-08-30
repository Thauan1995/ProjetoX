package utils

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestJSONUnmarshalValidate(t *testing.T) {

	type SubEstrutura struct {
		SubAtributoBool bool
		SubAtributoInt  int
	}

	type Estrutura struct {
		AtributoBool              bool
		AtributoFloat32           float32
		AtributoFloat64           float64
		AtributoInt               int
		AtributoInt32             int32
		AtributoInt64             int64
		AtributoMap               map[string]interface{}
		AtributoString            string
		AtributoSlice             []string
		AtributoTime              time.Time
		AtributoDate              JsonSpecialDate
		AtributoDateTime          JsonSpecialDateTime
		AtributoSubEstrutura      SubEstrutura
		AtributoSliceSubEstrutura []SubEstrutura
	}

	jsonValido := `{"AtributoBool": true, "AtributoFloat32": 15.23, "AtributoFloat64": 15.23, "AtributoInt": 15, "AtributoInt32": 16, "AtributoInt64": 17, "AtributoMap": {"Teste": "abc"}, "AtributoString": "teste", "AtributoSlice": ["teste1", "teste2", "teste3"], "AtributoTime": "2009-11-10T23:00:00Z", "AtributoDate": "0001-01-01", "AtributoDateTime": "2020-01-01 12:32:11", "AtributoSubEstrutura": {"SubAtributoBool": true, "SubAtributoInt": 15}, "AtributoSliceSubEstrutura": [{"SubAtributoBool": true, "SubAtributoInt": 15}]}`
	err := JSONUnmarshalValidate(jsonValido, reflect.TypeOf(Estrutura{}))
	if err != nil {
		t.Errorf("Erro ao validar JSON: %v", err)
		return
	}

	var estrutura Estrutura
	err = json.Unmarshal([]byte(jsonValido), &estrutura)
	if err != nil {
		t.Errorf("Erro ao realizar unmarshal do JSON: %v", err)
		return
	}
	return
}

func TestGetValorCampo(t *testing.T) {
	type Test struct {
		DataHora JsonSpecialDateTime
		Data     time.Time
	}

	var dataTest string = "2015-12-02 23:50:00"
	dataHora, err := ParseJsonSpecialDateTime("2006-01-02 15:04:05", dataTest)
	if err != nil {
		t.Errorf("dataTest inválida: %v", err)
	}

	data, err := time.Parse("2006-01-02 15:04:05", dataTest)
	if err != nil {
		t.Errorf("dataTest inválida: %v", err)
	}

	var test = Test{
		DataHora: dataHora,
		Data:     data,
	}

	// JsonSpecialDateTime
	tipo, valor, err := GetValorCampo("DataHora", reflect.ValueOf(test))
	if err != nil {
		t.Errorf("TestGetValorCampo retornou: %v", err)
	}
	if tipo != "time" {
		t.Errorf("TestGetValorCampo retornou tipo inválido")
	}
	if valor != dataTest {
		t.Errorf("TestGetValorCampo retornou valor inválido: %v", valor)
	}

	// time.Time
	tipo, valor, err = GetValorCampo("Data", reflect.ValueOf(test))
	if err != nil {
		t.Errorf("TestGetValorCampo retornou: %v", err)
	}
	if tipo != "time" {
		t.Errorf("TestGetValorCampo retornou tipo inválido")
	}
	if valor != dataTest {
		t.Errorf("TestGetValorCampo retornou valor inválido para campo data: %v", valor)
	}
}

func TestSetValorCampo(t *testing.T) {
	type OutroSliceStruct struct {
		AtributoFloat  float64
		AtributoInt    int64
		AtributoSlice  []string
		AtributoString string
		AtributoTime   time.Time
	}
	type SliceStruct struct {
		OutroSlice []OutroSliceStruct
	}
	type Estrutura struct {
		Slice []SliceStruct
	}

	dataHora, err := time.Parse("2006-01-02 15:04:05", "2020-10-30 19:20:18")
	if err != nil {
		t.Errorf("Falha ao criar data: %v", err)
	}

	var estrutura = Estrutura{
		[]SliceStruct{
			{
				[]OutroSliceStruct{
					{
						AtributoFloat: 1.5,
						AtributoInt:   1,
						AtributoSlice: []string{
							"Valor Inicial Slice",
						},
						AtributoString: "Valor Inicial String",
						AtributoTime:   dataHora,
					},
				},
			},
		},
	}

	//STRING
	var valorString string = "Valor Inserido"
	var campoString string = "Slice.OutroSlice.AtributoString"

	if err := SetValorCampo(campoString, &estrutura, valorString); err != nil {
		t.Errorf("Falha ao Setar string: %v", err)
		return
	}

	if estrutura.Slice[0].OutroSlice[0].AtributoString != valorString {
		t.Errorf("Atributo string inválido, esperado: %#v, retornou: %#v", valorString, estrutura.Slice[0].OutroSlice[0].AtributoString)
	}

	//INT64
	var valorInt int64 = 2
	var campoInt string = "Slice.OutroSlice.AtributoInt"

	if err := SetValorCampo(campoInt, &estrutura, valorInt); err != nil {
		t.Errorf("Falha ao Setar int64: %v", err)
		return
	}

	if estrutura.Slice[0].OutroSlice[0].AtributoInt != valorInt {
		t.Errorf("Atributo Int inválido, esperado: %#v, retornou: %#v", valorInt, estrutura.Slice[0].OutroSlice[0].AtributoInt)
	}

	//FLOAT64
	var valorFloat float64 = 2.5
	var campoFloat string = "Slice.OutroSlice.AtributoFloat"

	if err := SetValorCampo(campoFloat, &estrutura, valorFloat); err != nil {
		t.Errorf("Falha ao Setar float64: %v", err)
		return
	}

	if estrutura.Slice[0].OutroSlice[0].AtributoFloat != valorFloat {
		t.Errorf("Atributo Float inválido, esperado: %#v, retornou: %#v", valorFloat, estrutura.Slice[0].OutroSlice[0].AtributoFloat)
	}

	//TIME
	var valorTime time.Time = dataHora.AddDate(0, 0, 1)
	var campoTime string = "Slice.OutroSlice.AtributoTime"

	if err := SetValorCampo(campoTime, &estrutura, valorTime); err != nil {
		t.Errorf("Falha ao Setar Time64: %v", err)
		return
	}

	if !estrutura.Slice[0].OutroSlice[0].AtributoTime.Equal(valorTime) {
		t.Errorf("Atributo Time inválido, esperado: %#v, retornou: %#v", valorTime.Format("2006-01-02 15:04:05"), estrutura.Slice[0].OutroSlice[0].AtributoTime.Format("2006-01-02 15:04:05"))
	}

	//Slice String
	var valorSlice = []string{"Valor Inserido Slice"}
	var campoSlice string = "Slice.OutroSlice.AtributoSlice"

	if err := SetValorCampo(campoSlice, &estrutura, valorSlice); err != nil {
		t.Errorf("Falha ao Setar Slice: %v", err)
		return
	}

	if len(estrutura.Slice[0].OutroSlice[0].AtributoSlice) != len(valorSlice) {
		t.Errorf("Atributo Slice inválido, esperado: %#v, retornou: %#v", valorSlice, estrutura.Slice[0].OutroSlice[0].AtributoSlice)
		return
	}
	for _, indx := range estrutura.Slice[0].OutroSlice[0].AtributoSlice {
		for _, indxValor := range valorSlice {
			if indx != indxValor {
				t.Errorf("Atributo Slice inválido, esperado: %#v, retornou: %#v", valorSlice, estrutura.Slice[0].OutroSlice[0].AtributoSlice)
			}
		}
	}

	return
}
