//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"syscall/js"

	_ "embed"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"

	"github.com/ghdehrl12345/ZAGE/internal/age"
)

//go:embed zage.pk
var pkData []byte

var (
	ccs        constraint.ConstraintSystem
	compileErr error
)

func init() {
	ccs, compileErr = age.Compile()
}

func GenerateProof(this js.Value, args []js.Value) interface{} {
	if compileErr != nil {
		return fmt.Sprintf("Error: 회로 초기화 실패 - %v", compileErr)
	}
	if len(args) < 3 {
		return "Error: 인자가 부족합니다 (년도, 기준나이, 생년 필요)"
	}

	currentYear := args[0].Int()
	limitAge := args[1].Int()
	birthYear := args[2].Int()
	fmt.Printf("go-wasm: 증명 생성 시작 (현재:%d, 기준:%d, 생년:%d)\n", currentYear, limitAge, birthYear)

	pk := groth16.NewProvingKey(ecc.BN254)
	if _, err := pk.ReadFrom(bytes.NewReader(pkData)); err != nil {
		return fmt.Sprintf("Error: proving key 로딩 실패 - %v", err)
	}

	witness, err := age.NewPrivateWitness(currentYear, limitAge, birthYear)
	if err != nil {
		return fmt.Sprintf("Error: witness 생성 실패 - %v", err)
	}

	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		return fmt.Sprintf("Error: 증명 생성 실패 - %v", err)
	}

	var proofBuf bytes.Buffer
	if _, err := proof.WriteTo(&proofBuf); err != nil {
		return fmt.Sprintf("Error: proof 직렬화 실패 - %v", err)
	}
	return hex.EncodeToString(proofBuf.Bytes())
}

func main() {
	wait := make(chan struct{})
	js.Global().Set("verifyAgeZKP", js.FuncOf(GenerateProof))
	fmt.Println("✅ [Wasm] Go 성인인증 모듈 로드 완료")
	<-wait
}
