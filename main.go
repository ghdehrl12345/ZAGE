package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/consensys/gnark/backend/groth16"

	"github.com/ghdehrl12345/ZAGE/internal/age"
)

const (
	currentYear      = 2025
	limitAge         = 19
	exampleBirth     = 2005
	provingKeyFile   = "zage.pk"
	verifyingKeyFile = "zage.vk"
)

// keyWriter는 Groth16 키가 공통으로 구현하는 WriteTo 인터페이스를 추상화합니다.
type keyWriter interface {
	WriteTo(io.Writer) (int64, error)
}

func main() {
	fmt.Println("🚀 ZAGE: 키 생성 및 증명 테스트 시작")

	ccs, err := age.Compile()
	if err != nil {
		log.Fatalf("회로 컴파일 실패: %v", err)
	}

	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		log.Fatalf("Groth16 Setup 실패: %v", err)
	}

	if err := writeKey(provingKeyFile, pk); err != nil {
		log.Fatalf("증명 키 저장 실패: %v", err)
	}
	if err := writeKey(verifyingKeyFile, vk); err != nil {
		log.Fatalf("검증 키 저장 실패: %v", err)
	}
	fmt.Printf("✅ 키 저장 완료 (PK: %s, VK: %s)\n", provingKeyFile, verifyingKeyFile)

	witness, err := age.NewPrivateWitness(currentYear, limitAge, exampleBirth)
	if err != nil {
		log.Fatalf("비공개 witness 생성 실패: %v", err)
	}

	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		log.Fatalf("증명 생성 실패: %v", err)
	}
	fmt.Println("🧾 예시 증명 생성 완료")

	publicWitness, err := age.NewPublicWitness(currentYear, limitAge)
	if err != nil {
		log.Fatalf("공개 witness 생성 실패: %v", err)
	}

	if err := groth16.Verify(proof, vk, publicWitness); err != nil {
		log.Fatalf("검증 실패: %v", err)
	}
	fmt.Println("🎉 검증 성공: 로직이 정상적으로 동작합니다.")
}

// writeKey는 키 파일을 안전하게 생성하고 내용을 기록합니다.
func writeKey(path string, key keyWriter) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := key.WriteTo(file); err != nil {
		return err
	}
	return nil
}
