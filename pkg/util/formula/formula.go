package formula

import (
	"strconv"
	"strings"

	"github.com/dengsgo/math-engine/engine"
	"github.com/pkg/errors"
)

func Calc(formula string, params map[string]float64) (float64, error) {
	if formula == "" || len(params) == 0 {
		return 0, errors.Errorf("params is nil")
	}

	result, err := engine.ParseAndExec(replace(formula, params))
	if err != nil {
		return 0, errors.Wrapf(err, "formula calc failed")
	}
	return result, nil
}

func replace(f string, params map[string]float64) (ret string) {
	for k, v := range params {
		ret = strings.ReplaceAll(f, k, strconv.FormatFloat(v, 'f', 4, 64))
	}

	return ret
}
