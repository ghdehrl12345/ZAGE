package main

import (
	"fmt"
	"log"
	"os"

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
	fmt.Println("키 생성 및 저장 시작")

	var circuit AgeCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		log.Fatal("회로 컴파일 실패:", err)
	}

	// 증명 키(PK)와 검증 키(VK) 생성
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		log.Fatal("Setup 실패:", err)
	}

	// 증명 키 저장
	pkFile, _ := os.Create("zage.pk")
	pk.WriteTo(pkFile)
	pkFile.Close()
	fmt.Println("증명 키(zage.pk) 저장 완료!")

	// 검증 키(Verifying Key) 저장
	vkFile, _ := os.Create("zage.vk")
	vk.WriteTo(vkFile)
	vkFile.Close()
	fmt.Println("검증 키(zage.vk) 저장 완료")

	// 5. 테스트
	fmt.Println("\n--- [테스트: 저장된 키로 증명 해보기] ---")

	// 예시: 2005년생(20세)
	witness, _ := frontend.NewWitness(&AgeCircuit{
		CurrentYear: 2025,
		LimitAge:    19,
		BirthYear:   2005,
	}, ecc.BN254.ScalarField())

	// 증명 생성
	proof, _ := groth16.Prove(ccs, pk, witness)

	// 검증
	publicWitness, _ := frontend.NewWitness(&AgeCircuit{
		CurrentYear: 2025,
		LimitAge:    19,
	}, ecc.BN254.ScalarField(), frontend.PublicOnly())

	err = groth16.Verify(proof, vk, publicWitness)
	if err == nil {
		fmt.Println("🎉 테스트 성공: 로직에 문제 없습니다.")
	}
}
