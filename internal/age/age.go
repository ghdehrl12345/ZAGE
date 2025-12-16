package age

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// ScalarField은 프로젝트 전반에서 사용하는 곡선 필드를 고정해 둡니다.
var ScalarField = ecc.BN254.ScalarField()

// Circuit는 나이 비교 회로 정의입니다.
type Circuit struct {
	CurrentYear frontend.Variable `gnark:",public"`
	LimitAge    frontend.Variable `gnark:",public"`
	BirthYear   frontend.Variable
}

// Define는 "나이가 기준 이상인지" 제약식을 구성합니다.
func (circuit *Circuit) Define(api frontend.API) error {
	myAge := api.Sub(circuit.CurrentYear, circuit.BirthYear)
	api.AssertIsLessOrEqual(circuit.LimitAge, myAge)
	return nil
}

// Compile은 AgeCircuit을 R1CS 형태로 변환합니다.
func Compile() (constraint.ConstraintSystem, error) {
	var circuit Circuit
	return frontend.Compile(ScalarField, r1cs.NewBuilder, &circuit)
}

// NewPrivateWitness는 비공개 증명용 witness를 생성합니다.
func NewPrivateWitness(currentYear, limitAge, birthYear int) (witness.Witness, error) {
	assignment := &Circuit{
		CurrentYear: currentYear,
		LimitAge:    limitAge,
		BirthYear:   birthYear,
	}
	witness, err := frontend.NewWitness(assignment, ScalarField)
	if err != nil {
		return nil, fmt.Errorf("private witness 생성 실패: %w", err)
	}
	return witness, nil
}

// NewPublicWitness는 서버 검증용 공개 witness를 생성합니다.
func NewPublicWitness(currentYear, limitAge int) (witness.Witness, error) {
	assignment := &Circuit{
		CurrentYear: currentYear,
		LimitAge:    limitAge,
	}
	witness, err := frontend.NewWitness(assignment, ScalarField, frontend.PublicOnly())
	if err != nil {
		return nil, fmt.Errorf("public witness 생성 실패: %w", err)
	}
	return witness, nil
}
