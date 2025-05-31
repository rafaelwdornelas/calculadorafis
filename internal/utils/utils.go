package utils

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"strconv"
	"strings"
)

// ProcessarValorMonetario processa um valor monetário em formato brasileiro para float64
func ProcessarValorMonetario(valorStr string) (float64, error) {
	valorStr = strings.TrimSpace(valorStr)
	valorStr = strings.Replace(valorStr, "R$", "", -1)
	valorStr = strings.Replace(valorStr, " ", "", -1)

	valorStr = strings.Replace(valorStr, ".", "", -1)  // Remove pontos de milhar
	valorStr = strings.Replace(valorStr, ",", ".", -1) // Substitui vírgula por ponto decimal

	valor, err := strconv.ParseFloat(valorStr, 64)
	if err != nil {
		return 0, err
	}

	return valor, nil
}

// ToFloat64 converte um valor para float64
func ToFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, false
		}
		return f, true
	default:
		return 0, false
	}
}

// GetTemplateFuncs retorna as funções para uso nos templates
func GetTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"sub": func(values ...interface{}) float64 {
			if len(values) != 2 {
				log.Printf("ERRO: Função 'sub' chamada com %d argumentos, esperava 2", len(values))
				return 0
			}
			a, aOk := ToFloat64(values[0])
			b, bOk := ToFloat64(values[1])
			if !aOk || !bOk {
				log.Printf("ERRO: Função 'sub' chamada com tipos incompatíveis")
				return 0
			}
			return a - b
		},
		"add": func(values ...interface{}) float64 {
			if len(values) != 2 {
				log.Printf("ERRO: Função 'add' chamada com %d argumentos, esperava 2", len(values))
				return 0
			}
			a, aOk := ToFloat64(values[0])
			b, bOk := ToFloat64(values[1])
			if !aOk || !bOk {
				log.Printf("ERRO: Função 'add' chamada com tipos incompatíveis")
				return 0
			}
			return a + b
		},
		"mul": func(values ...interface{}) float64 {
			if len(values) != 2 {
				log.Printf("ERRO: Função 'mul' chamada com %d argumentos, esperava 2", len(values))
				return 0
			}
			a, aOk := ToFloat64(values[0])
			b, bOk := ToFloat64(values[1])
			if !aOk || !bOk {
				log.Printf("ERRO: Função 'mul' chamada com tipos incompatíveis")
				return 0
			}
			return a * b
		},
		"div": func(values ...interface{}) float64 {
			if len(values) != 3 {
				log.Printf("ERRO: Função 'div' chamada com %d argumentos, esperava 3", len(values))
				return 0
			}
			a, aOk := ToFloat64(values[0])
			b, bOk := ToFloat64(values[1])
			multiplier, mOk := ToFloat64(values[2])
			if !aOk || !bOk || !mOk {
				log.Printf("ERRO: Função 'div' chamada com tipos incompatíveis")
				return 0
			}
			if b == 0 {
				return 0
			}
			return (a / b) * multiplier
		},
		"float64": func(i interface{}) float64 {
			v, _ := ToFloat64(i)
			return v
		},
		"formatMoney": func(value interface{}) string {
			v, ok := ToFloat64(value)
			if !ok {
				return "0,00"
			}

			// Se o valor for muito pequeno (menor que 0.01), considera como 0
			if math.Abs(v) < 0.01 {
				return "0,00"
			}
			formatted := fmt.Sprintf("%.2f", v)
			parts := strings.Split(formatted, ".")
			intPart := parts[0]
			decPart := parts[1]

			var result string
			for i := len(intPart) - 1; i >= 0; i-- {
				if (len(intPart)-i-1)%3 == 0 && i < len(intPart)-1 {
					result = "." + result
				}
				result = string(intPart[i]) + result
			}

			return result + "," + decPart
		},
		"calcularDiferenca": func(final, inicial float64) float64 {
			return final - inicial
		},
		"calcularPercentual": func(final, inicial float64) float64 {
			if inicial <= 0 {
				return 0
			}
			return ((final - inicial) / inicial) * 100
		},
	}
}
