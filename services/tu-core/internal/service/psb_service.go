package service

import (
	"context"
	"fmt"
	"sort"

	"sekolah-platform/services/tu-core/internal/repository"
)

type PSBService struct {
	repo repository.PSBRepository
}

func NewPSBService(repo repository.PSBRepository) *PSBService {
	return &PSBService{repo: repo}
}

func (s *PSBService) GetCalonSiswa(ctx context.Context, schoolID, periodID string) ([]repository.CalonSiswa, error) {
	return s.repo.GetCalonSiswa(ctx, schoolID, periodID)
}

func (s *PSBService) GenerateKelas(ctx context.Context, schoolID, periodID string) error {
	siswaList, err := s.repo.GetCalonSiswa(ctx, schoolID, periodID)
	if err != nil {
		return err
	}

	// Grouping by Jurusan
	grouped := make(map[string][]repository.CalonSiswa)
	for _, s := range siswaList {
		jur := s.JurusanID
		if jur == "" { jur = "TANPA_JURUSAN" }
		grouped[jur] = append(grouped[jur], s)
	}

	MAX_KAPASITAS := 36
	hasilRombel := make(map[string][]repository.CalonSiswa)

	for jur, list := range grouped {
		// Sort by name
		sort.Slice(list, func(i, j int) bool { return list[i].Nama < list[j].Nama })
		
		total := len(list)
		if total == 0 { continue }
		jumlahKelas := (total + MAX_KAPASITAS - 1) / MAX_KAPASITAS
		baseSiswa := total / jumlahKelas
		sisa := total % jumlahKelas

		idx := 0
		for i := 0; i < jumlahKelas; i++ {
			label := string(rune(65 + i)) // A, B, C
			namaKelas := fmt.Sprintf("%s %s", jur, label)
			target := baseSiswa
			if sisa > 0 {
				target++
				sisa--
			}
			if idx+target > total {
				target = total - idx
			}
			hasilRombel[namaKelas] = list[idx : idx+target]
			idx += target
		}
	}

	return s.repo.GenerateKelas(ctx, schoolID, periodID, hasilRombel)
}
