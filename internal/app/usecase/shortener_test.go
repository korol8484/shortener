package usecase

import "testing"

func TestGenAlias(t *testing.T) {
	alias := GenAlias(6, "testString")
	if alias != "Jlf8iW" {
		t.Fatal("invalid alias generated")
	}
}

func BenchmarkGenAlias(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenAlias(5, "https://ya.ru")
	}
}
