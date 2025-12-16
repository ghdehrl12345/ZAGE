package main

import (
	"fmt"
	"log"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

type AgeCircuit struct {
	CurrentYear frontend.Variable `gnark:",public"`
	LimitAge    frontend.Variable `gnark:",public"`
	BirthYear   frontend.Variable
}

func (circuit *AgeCircuit) Define(api frontend.API) error {
	myAge := api.Sub(circuit.CurrentYear, circuit.BirthYear)
	api.AssertIsLessOrEqual(circuit.LimitAge, myAge)
	return nil
}

func main() {
	fmt.Println("ZAGE 프로젝트 시동 중")

	var circuit AgeCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		log.Fatal("회로 컴파일 실패:", err)
	}

	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		log.Fatal("Setup 실패:", err)
	}

	witness, err := frontend.NewWitness(&AgeCircuit{
		CurrentYear: 2025,
		LimitAge:    19,
		BirthYear:   2005,
	}, ecc.BN254.ScalarField())
	if err != nil {
		log.Fatal("비공개 위트니스 생성 실패:", err)
	}

	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		log.Fatal("증명 생성 실패:", err)
	}
	fmt.Println("성인 인증 증명서 생성 완료 (생년월일은 숨겨짐)")

	publicWitness, err := frontend.NewWitness(&AgeCircuit{
		CurrentYear: 2025,
		LimitAge:    19,
	}, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		log.Fatal("공개 위트니스 생성 실패:", err)
	}

	err = groth16.Verify(proof, vk, publicWitness)
	if err == nil {
		fmt.Println("검증 성공! 이 사용자는 성인이 확실합니다.")
	} else {
		fmt.Println("검증 실패! 거짓말쟁이거나 미성년자입니다.")
	}
}
